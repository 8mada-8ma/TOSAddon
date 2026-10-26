//go:build windows

package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
	"golang.org/x/sys/windows"
)

// 二重起動の防止。版が違っても同じ名前を使い続けること。
const instanceMutexName = "Local\\midiflute_tos_midi_helper"

var instanceMutex windows.Handle

const appVersion = "0.3.0"

type App struct {
	mu            sync.Mutex
	settings      Settings
	tosDir        string
	bridge        *Bridge
	midi          *MidiIn
	wantConnected bool
	held          map[int]int  // 押している鍵盤 -> 鳴らしている音番号
	pressed       map[int]bool // 画面表示用
	learning      bool
	status        GameStatus
	lastOct3      bool
	devices       []string
	midiErr       string
	fgVal         bool
	fgAt          time.Time
	foreign       bool
	dirty         int32

	mw      *walk.MainWindow
	devCB   *walk.ComboBox
	connPB  *walk.PushButton
	tosLE   *walk.LineEdit
	midiLbl *walk.Label
	tosLbl  *walk.Label
	latLbl  *walk.Label
	baseLbl *walk.Label
	hintLbl *walk.Label
	kb      *walk.CustomWidget
	fgCB    *walk.CheckBox
	ni      *walk.NotifyIcon
	font    *walk.Font
}

func main() {
	name, _ := windows.UTF16PtrFromString(instanceMutexName)
	h, err := windows.CreateMutex(nil, false, name)
	if err == windows.ERROR_ALREADY_EXISTS {
		walk.MsgBox(nil, "midiflute",
			"midiflute.exe はすでに起動しています。\n\nタスクトレイ(画面右下)のアイコンをクリックすると、起動中の画面を開けます。",
			walk.MsgBoxIconInformation)
		return
	}
	instanceMutex = h
	disablePowerThrottling()
	a := &App{held: map[int]int{}, pressed: map[int]bool{}}
	a.settings = LoadSettings()
	a.wantConnected = a.settings.AutoConnect
	a.devices = ListMidiDevices()
	if a.settings.Device == "" && len(a.devices) > 0 {
		a.settings.Device = a.devices[0]
	}
	if d := FindTosDir(a.settings); d != "" {
		a.setTosDir(d)
	}
	a.font, _ = walk.NewFont("Meiryo UI", 8, 0)

	if err := a.buildWindow(); err != nil {
		walk.MsgBox(nil, "midiflute", "画面を作れませんでした: "+err.Error(), walk.MsgBoxIconError)
		return
	}
	a.setupTray()
	go a.midiLoop()
	go a.maintenanceLoop()
	go a.ackLoop()
	go a.repaintLoop()
	a.refreshUI()
	a.mw.Run()

	// 終了処理: 鳴っている音を止めてから機器を閉じる
	a.mu.Lock()
	a.stopAllLocked()
	a.midi.Close()
	a.mu.Unlock()
	if a.ni != nil {
		a.ni.Dispose()
	}
}

func (a *App) setTosDir(d string) {
	a.tosDir = d
	a.bridge = NewBridge(AddonDir(d))
	a.bridge.Heartbeat()
	a.settings.TosDir = d
	a.settings.Save()
}

// ---- MIDI 受信 ----

func (a *App) midiLoop() {
	for msg := range midiCh {
		status := msg & 0xF0
		note := int((msg >> 8) & 0x7F)
		vel := (msg >> 16) & 0x7F
		on := status == 0x90 && vel > 0
		off := status == 0x80 || (status == 0x90 && vel == 0)
		if !on && !off {
			continue
		}
		a.mu.Lock()
		if on {
			a.noteOnLocked(note)
		} else {
			a.noteOffLocked(note)
		}
		a.mu.Unlock()
		atomic.StoreInt32(&a.dirty, 1)
	}
}

func (a *App) oct3Locked() bool {
	if a.status.Linked {
		return a.status.Oct3
	}
	return a.lastOct3
}

func (a *App) foreground() bool {
	if time.Since(a.fgAt) > 100*time.Millisecond {
		a.fgVal = foregroundIsToS()
		a.fgAt = time.Now()
	}
	return a.fgVal
}

func (a *App) noteOnLocked(note int) {
	a.pressed[note] = true
	if a.learning {
		a.learning = false
		a.setBaseLocked(note)
		return
	}
	if a.bridge == nil {
		return
	}
	if a.settings.ForegroundOnly && !a.foreground() {
		return
	}
	id, _ := noteToID(note, a.settings.Base, a.oct3Locked())
	if id == 0 {
		return
	}
	if prev, ok := a.held[note]; ok {
		a.bridge.Send('S', prev)
	}
	a.bridge.Send('P', id)
	a.held[note] = id
}

func (a *App) noteOffLocked(note int) {
	delete(a.pressed, note)
	id, ok := a.held[note]
	if !ok {
		return
	}
	delete(a.held, note)
	if a.bridge != nil {
		a.bridge.Send('S', id)
	}
}

func (a *App) stopAllLocked() {
	for note, id := range a.held {
		if a.bridge != nil {
			a.bridge.Send('S', id)
		}
		delete(a.held, note)
	}
}

func (a *App) setBaseLocked(base int) {
	if base < 0 {
		base = 0
	}
	if base > 127 {
		base = 127
	}
	a.stopAllLocked()
	a.settings.Base = base
	a.settings.Save()
	atomic.StoreInt32(&a.dirty, 1)
	go a.mw.Synchronize(a.refreshUI)
}

// ---- 定期処理(受信とは別の流れで動かし、受信を止めない) ----

func (a *App) maintenanceLoop() {
	tick := time.NewTicker(250 * time.Millisecond)
	n := 0
	for range tick.C {
		n++
		a.mu.Lock()
		if a.bridge != nil {
			if n%8 == 1 {
				a.bridge.Heartbeat()
			}
			a.status = ReadStatus(a.bridge.dir)
			a.foreign = a.bridge.ForeignWriter()
			if a.status.Linked {
				a.lastOct3 = a.status.Oct3
			}
		}
		a.mu.Unlock()

		if n%8 == 1 {
			devs := ListMidiDevices()
			a.mu.Lock()
			a.devices = devs
			if a.midi != nil && !contains(devs, a.midi.Name) {
				a.stopAllLocked()
				a.pressed = map[int]bool{}
				a.midi.Close()
				a.midi = nil
				a.midiErr = "機器が外れました"
			}
			if a.midi == nil && a.wantConnected && contains(devs, a.settings.Device) {
				a.openLocked(a.settings.Device)
			}
			a.mu.Unlock()
		}
		a.mw.Synchronize(a.refreshUI)
	}
}

func (a *App) ackLoop() {
	for {
		a.mu.Lock()
		b := a.bridge
		a.mu.Unlock()
		if b != nil && b.Poll() {
			time.Sleep(2 * time.Millisecond)
		} else {
			time.Sleep(20 * time.Millisecond)
		}
	}
}

func (a *App) repaintLoop() {
	for {
		time.Sleep(30 * time.Millisecond)
		if atomic.SwapInt32(&a.dirty, 0) == 1 && a.kb != nil {
			a.mw.Synchronize(func() { a.kb.Invalidate() })
		}
	}
}

func (a *App) openLocked(name string) {
	m, err := OpenMidi(name)
	if err != nil {
		a.midiErr = err.Error()
		return
	}
	a.midi = m
	a.midiErr = ""
}

func contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}

// ---- 画面 ----

func (a *App) buildWindow() error {
	return MainWindow{
		AssignTo: &a.mw,
		Title:    "midiflute - ToS 笛演奏用 MIDI 受信 v" + appVersion,
		Icon:     walk.IconApplication(),
		MinSize:  Size{Width: 820, Height: 470},
		Size:     Size{Width: 900, Height: 500},
		Font:     Font{Family: "Meiryo UI", PointSize: 9},
		Layout:   VBox{},
		Children: []Widget{
			GroupBox{
				Title:  "接続",
				Layout: Grid{Columns: 4},
				Children: []Widget{
					Label{Text: "MIDI機器"},
					ComboBox{AssignTo: &a.devCB, Model: a.devices, Value: a.settings.Device, Editable: false,
						OnCurrentIndexChanged: a.onDeviceChanged},
					PushButton{AssignTo: &a.connPB, Text: "接続", OnClicked: a.onConnect},
					PushButton{Text: "一覧を更新", OnClicked: a.onRefreshDevices},
					Label{Text: "ToSフォルダ"},
					LineEdit{AssignTo: &a.tosLE, ReadOnly: true, Text: a.tosDir},
					PushButton{Text: "変更…", OnClicked: a.onChooseTos},
					HSpacer{},
				},
			},
			GroupBox{
				Title:  "状態",
				Layout: Grid{Columns: 3},
				Children: []Widget{
					Label{Text: "MIDI:"}, Label{AssignTo: &a.midiLbl}, HSpacer{},
					Label{Text: "ToS:"}, Label{AssignTo: &a.tosLbl}, HSpacer{},
					Label{Text: "遅延:"}, Label{AssignTo: &a.latLbl}, HSpacer{},
				},
			},
			GroupBox{
				Title:  "鍵盤の割り当て",
				Layout: VBox{},
				Children: []Widget{
					Composite{
						Layout: HBox{MarginsZero: true},
						Children: []Widget{
							Label{AssignTo: &a.baseLbl},
							PushButton{Text: "－1オクターブ", OnClicked: func() { a.shiftBase(-12) }},
							PushButton{Text: "－半音", OnClicked: func() { a.shiftBase(-1) }},
							PushButton{Text: "＋半音", OnClicked: func() { a.shiftBase(1) }},
							PushButton{Text: "＋1オクターブ", OnClicked: func() { a.shiftBase(12) }},
							PushButton{Text: "鍵盤で指定", OnClicked: a.onLearn},
							HSpacer{},
						},
					},
					Label{AssignTo: &a.hintLbl},
					CustomWidget{
						AssignTo:            &a.kb,
						MinSize:             Size{Width: 760, Height: 150},
						ClearsBackground:    true,
						InvalidatesOnResize: true,
						PaintPixels:         a.paintKeyboard,
					},
					Label{Text: "青=1オクターブ目　緑=2オクターブ目　橙=3オクターブ目　灰=鳴らない鍵盤(範囲外・未解放)　赤=押している鍵盤"},
				},
			},
			Composite{
				Layout: HBox{MarginsZero: true},
				Children: []Widget{
					CheckBox{AssignTo: &a.fgCB, Text: "ToSが前面にあるときだけ演奏を送る", Checked: a.settings.ForegroundOnly,
						OnCheckedChanged: func() {
							a.mu.Lock()
							a.settings.ForegroundOnly = a.fgCB.Checked()
							a.settings.Save()
							a.mu.Unlock()
						}},
					HSpacer{},
				},
			},
		},
	}.Create()
}

func (a *App) setupTray() {
	ni, err := walk.NewNotifyIcon(a.mw)
	if err != nil {
		return
	}
	a.ni = ni
	ni.SetIcon(walk.IconApplication())
	ni.SetToolTip("midiflute(クリックで表示)")
	ni.MouseDown().Attach(func(x, y int, button walk.MouseButton) {
		if button == walk.LeftButton {
			a.restoreFromTray()
		}
	})
	show := walk.NewAction()
	show.SetText("表示")
	show.Triggered().Attach(a.restoreFromTray)
	quit := walk.NewAction()
	quit.SetText("終了")
	quit.Triggered().Attach(func() { a.mw.Close() })
	ni.ContextMenu().Actions().Add(show)
	ni.ContextMenu().Actions().Add(quit)

	// 最小化したらタスクトレイに格納する
	a.mw.SizeChanged().Attach(func() {
		if win.IsIconic(a.mw.Handle()) {
			a.mw.SetVisible(false)
			ni.SetVisible(true)
		}
	})
}

func (a *App) restoreFromTray() {
	a.mw.SetVisible(true)
	win.ShowWindow(a.mw.Handle(), win.SW_RESTORE)
	win.SetForegroundWindow(a.mw.Handle())
	if a.ni != nil {
		a.ni.SetVisible(false)
	}
}

func (a *App) onDeviceChanged() {
	name := a.devCB.Text()
	a.mu.Lock()
	defer a.mu.Unlock()
	if name == "" || name == a.settings.Device {
		return
	}
	a.settings.Device = name
	a.settings.Save()
	if a.midi != nil {
		a.stopAllLocked()
		a.midi.Close()
		a.midi = nil
		a.openLocked(name)
	}
}

func (a *App) onConnect() {
	a.mu.Lock()
	if a.midi != nil {
		a.stopAllLocked()
		a.pressed = map[int]bool{}
		a.midi.Close()
		a.midi = nil
		a.wantConnected = false
		a.midiErr = ""
	} else {
		a.wantConnected = true
		if a.devCB.Text() != "" {
			a.settings.Device = a.devCB.Text()
		}
		a.openLocked(a.settings.Device)
	}
	a.settings.AutoConnect = a.wantConnected
	a.settings.Save()
	a.mu.Unlock()
	atomic.StoreInt32(&a.dirty, 1)
	a.refreshUI()
}

func (a *App) onRefreshDevices() {
	devs := ListMidiDevices()
	a.mu.Lock()
	a.devices = devs
	cur := a.settings.Device
	a.mu.Unlock()
	a.devCB.SetModel(devs)
	for i, d := range devs {
		if d == cur {
			a.devCB.SetCurrentIndex(i)
		}
	}
	a.refreshUI()
}

func (a *App) onChooseTos() {
	dlg := walk.FileDialog{Title: "Tree of Savior のインストールフォルダを選んでください"}
	ok, err := dlg.ShowBrowseFolder(a.mw)
	if err != nil || !ok {
		return
	}
	if !IsTosDir(dlg.FilePath) {
		walk.MsgBox(a.mw, "midiflute", "このフォルダには release と data がありません。\nTree of Savior (Japanese Ver.) のフォルダを選んでください。", walk.MsgBoxIconWarning)
		return
	}
	a.mu.Lock()
	a.stopAllLocked()
	a.setTosDir(dlg.FilePath)
	a.mu.Unlock()
	a.tosLE.SetText(dlg.FilePath)
	a.refreshUI()
}

func (a *App) shiftBase(d int) {
	a.mu.Lock()
	a.setBaseLocked(a.settings.Base + d)
	a.mu.Unlock()
}

func (a *App) onLearn() {
	a.mu.Lock()
	a.learning = true
	a.mu.Unlock()
	a.refreshUI()
}

func (a *App) refreshUI() {
	if a.mw == nil {
		return
	}
	a.mu.Lock()
	connected := a.midi != nil
	midiName := ""
	if connected {
		midiName = a.midi.Name
	}
	midiErr := a.midiErr
	want := a.wantConnected
	st := a.status
	oct3 := a.oct3Locked()
	tosDir := a.tosDir
	base := a.settings.Base
	learning := a.learning
	b := a.bridge
	foreign := a.foreign
	a.mu.Unlock()

	switch {
	case connected:
		a.midiLbl.SetText("接続中: " + midiName)
		a.connPB.SetText("切断")
	case want:
		msg := "待機中(機器をつなぐと自動で接続します)"
		if midiErr != "" {
			msg = midiErr + " / " + msg
		}
		a.midiLbl.SetText(msg)
		a.connPB.SetText("切断")
	default:
		a.midiLbl.SetText("切断中(他のソフトでキーボードを使えます)")
		a.connPB.SetText("接続")
	}

	o3 := "未解放"
	if oct3 {
		o3 = "解放済み"
	}
	switch {
	case foreign:
		a.tosLbl.SetText("⚠ ほかの midiflute.exe(別の版や別の場所のものを含む)が同時に動いているようです。1つだけにしてください。")
	case tosDir == "":
		a.tosLbl.SetText("ToSフォルダが未設定です。「変更…」から選んでください。")
	case !st.Linked:
		a.tosLbl.SetText("未連携(ToSを起動し、キャラクターでログインすると連携します)")
	case !st.Playing:
		a.tosLbl.SetText("連携中(演奏モードではありません) / 3オクターブ目: " + o3)
	default:
		a.tosLbl.SetText("演奏モード中 ♪ / 3オクターブ目: " + o3)
	}

	if b != nil {
		if last, med, ok := b.Latency(); ok {
			a.latLbl.SetText(fmt.Sprintf("直近 %.0f ms / 中央値 %.0f ms", last, med))
		} else {
			a.latLbl.SetText("―(演奏すると表示されます)")
		}
	}

	a.baseLbl.SetText(fmt.Sprintf("1オクターブ目のド: %s(鍵盤番号 %d)　", noteName(base), base))
	if learning {
		a.hintLbl.SetText("▶ 1オクターブ目の「ド」にしたい鍵盤を押してください")
	} else {
		a.hintLbl.SetText("microKEYなどのオクターブ切替を動かしたときは、「鍵盤で指定」でもう一度合わせてください。")
	}
}

type keyColors struct{ white, black walk.Color }

var colorByKind = map[int]keyColors{
	KeyOut:    {walk.RGB(232, 232, 232), walk.RGB(55, 55, 55)},
	KeyOct1:   {walk.RGB(175, 210, 255), walk.RGB(35, 85, 170)},
	KeyOct2:   {walk.RGB(170, 230, 180), walk.RGB(30, 120, 60)},
	KeyOct3:   {walk.RGB(255, 210, 160), walk.RGB(190, 110, 30)},
	KeyLocked: {walk.RGB(205, 205, 205), walk.RGB(110, 110, 110)},
}

func (a *App) paintKeyboard(canvas *walk.Canvas, _ walk.Rectangle) error {
	a.mu.Lock()
	base := a.settings.Base
	oct3 := a.oct3Locked()
	pressed := make(map[int]bool, len(a.pressed))
	for k, v := range a.pressed {
		pressed[k] = v
	}
	a.mu.Unlock()

	bounds := a.kb.ClientBoundsPixels()
	lo, hi := base-12, base+47
	if lo < 0 {
		lo = 0
	}
	if hi > 127 {
		hi = 127
	}
	for lo > 0 && isBlack(lo) {
		lo--
	}
	var whites []int
	for n := lo; n <= hi; n++ {
		if !isBlack(n) {
			whites = append(whites, n)
		}
	}
	if len(whites) == 0 {
		return nil
	}
	ww := float64(bounds.Width-1) / float64(len(whites))
	H := bounds.Height - 1
	pen, _ := walk.NewCosmeticPen(walk.PenSolid, walk.RGB(90, 90, 90))
	defer pen.Dispose()

	fill := func(c walk.Color, r walk.Rectangle) {
		br, err := walk.NewSolidColorBrush(c)
		if err != nil {
			return
		}
		canvas.FillRectanglePixels(br, r)
		br.Dispose()
	}
	label := func(n int) string {
		switch n - base {
		case 0:
			return "ド1"
		case 12:
			return "ド2"
		case 24:
			return "ド3"
		}
		if n%12 == 0 {
			return noteName(n)
		}
		return ""
	}
	colorFor := func(n int, black bool) walk.Color {
		_, kind := noteToID(n, base, oct3)
		if pressed[n] {
			if kind == KeyOut || kind == KeyLocked {
				return walk.RGB(150, 150, 150)
			}
			return walk.RGB(235, 60, 60)
		}
		if black {
			return colorByKind[kind].black
		}
		return colorByKind[kind].white
	}

	whiteIdx := map[int]int{}
	for i, n := range whites {
		whiteIdx[n] = i
		r := walk.Rectangle{X: bounds.X + int(float64(i)*ww), Y: bounds.Y, Width: int(ww) + 1, Height: H}
		fill(colorFor(n, false), r)
		canvas.DrawRectanglePixels(pen, r)
		if s := label(n); s != "" {
			tr := walk.Rectangle{X: r.X, Y: r.Y + H - 22, Width: r.Width, Height: 20}
			canvas.DrawTextPixels(s, a.font, walk.RGB(20, 20, 20), tr, walk.TextCenter|walk.TextSingleLine|walk.TextVCenter)
		}
	}
	bw := int(ww * 0.62)
	bh := int(float64(H) * 0.6)
	for n := lo; n <= hi; n++ {
		if !isBlack(n) {
			continue
		}
		i, ok := whiteIdx[n-1]
		if !ok {
			continue
		}
		x := bounds.X + int(float64(i+1)*ww) - bw/2
		r := walk.Rectangle{X: x, Y: bounds.Y, Width: bw, Height: bh}
		fill(colorFor(n, true), r)
		canvas.DrawRectanglePixels(pen, r)
		if s := label(n); s != "" {
			tr := walk.Rectangle{X: r.X - 6, Y: r.Y + bh - 20, Width: bw + 12, Height: 18}
			canvas.DrawTextPixels(s, a.font, walk.RGB(255, 255, 255), tr, walk.TextCenter|walk.TextSingleLine|walk.TextVCenter)
		}
	}
	return nil
}
