package main

// 表示言語(日本語・English・한국어)

var langCodes = []string{"ja", "en", "ko"}
var langNames = []string{"日本語", "English", "한국어"}

var curLang = "ja"

func setLang(code string) {
	for _, c := range langCodes {
		if c == code {
			curLang = code
			return
		}
	}
	curLang = "ja"
}

func langIndex(code string) int {
	for i, c := range langCodes {
		if c == code {
			return i
		}
	}
	return 0
}

// T は現在の言語の文字列を返す(無ければ日本語)
func T(key string) string {
	if m, ok := texts[curLang]; ok {
		if s, ok := m[key]; ok {
			return s
		}
	}
	return texts["ja"][key]
}

var texts = map[string]map[string]string{
	"ja": {
		"title":           "midiflute - ToS 笛演奏用 MIDI 受信",
		"grp_conn":        "接続",
		"midi_device":     "MIDI機器",
		"connect":         "接続",
		"disconnect":      "切断",
		"refresh":         "一覧を更新",
		"tos_folder":      "ToSフォルダ",
		"change":          "変更…",
		"language":        "言語 / Language",
		"grp_status":      "状態",
		"latency":         "遅延:",
		"grp_keys":        "鍵盤の割り当て",
		"oct_down":        "－1オクターブ",
		"semi_down":       "－半音",
		"semi_up":         "＋半音",
		"oct_up":          "＋1オクターブ",
		"learn":           "鍵盤で指定",
		"legend":          "青=1オクターブ目　緑=2オクターブ目　橙=3オクターブ目　灰=鳴らない鍵盤(範囲外・未解放)　赤=押している鍵盤",
		"fg_only":         "ToSが前面にあるときだけ演奏を送る",
		"tray_tip":        "midiflute(クリックで表示)",
		"tray_show":       "表示",
		"tray_quit":       "終了",
		"already_running": "midiflute.exe はすでに起動しています。\n\nタスクトレイ(画面右下)のアイコンをクリックすると、起動中の画面を開けます。",
		"ui_error":        "画面を作れませんでした: ",
		"choose_tos":      "Tree of Savior のインストールフォルダを選んでください",
		"not_tos":         "このフォルダには release と data がありません。\nTree of Savior (Japanese Ver.) のフォルダを選んでください。",
		"midi_connected":  "接続中: ",
		"midi_waiting":    "待機中(機器をつなぐと自動で接続します)",
		"midi_off":        "切断中(他のソフトでキーボードを使えます)",
		"midi_busy":       "この機器は開けませんでした(他のソフトが使用中の可能性があります)",
		"midi_notfound":   "機器が見つかりません",
		"midi_removed":    "機器が外れました",
		"oct3_open":       "解放済み",
		"oct3_locked":     "未解放",
		"tos_foreign":     "⚠ ほかの midiflute.exe(別の版や別の場所のものを含む)が同時に動いているようです。1つだけにしてください。",
		"tos_nofolder":    "ToSフォルダが未設定です。「変更…」から選んでください。",
		"tos_unlinked":    "未連携(ToSを起動し、キャラクターでログインすると連携します)",
		"tos_linked":      "連携中(演奏モードではありません) / 3オクターブ目: ",
		"tos_playing":     "演奏モード中 ♪ / 3オクターブ目: ",
		"lat_fmt":         "直近 %.0f ms / 中央値 %.0f ms",
		"lat_none":        "―(演奏すると表示されます)",
		"base_fmt":        "1オクターブ目のド: %s(鍵盤番号 %d)　",
		"hint_learning":   "▶ 1オクターブ目の「ド」にしたい鍵盤を押してください",
		"hint_normal":     "キーボードのオクターブ切替を動かしたときは、「鍵盤で指定」でもう一度合わせてください。",
		"do":              "ド",
	},
	"en": {
		"title":           "midiflute - MIDI receiver for ToS flute playing",
		"grp_conn":        "Connection",
		"midi_device":     "MIDI device",
		"connect":         "Connect",
		"disconnect":      "Disconnect",
		"refresh":         "Refresh",
		"tos_folder":      "ToS folder",
		"change":          "Change…",
		"language":        "言語 / Language",
		"grp_status":      "Status",
		"latency":         "Latency:",
		"grp_keys":        "Key mapping",
		"oct_down":        "-1 octave",
		"semi_down":       "-1 semitone",
		"semi_up":         "+1 semitone",
		"oct_up":          "+1 octave",
		"learn":           "Set by key",
		"legend":          "Blue=1st octave  Green=2nd octave  Orange=3rd octave  Gray=silent keys (out of range / locked)  Red=pressed",
		"fg_only":         "Send notes only while ToS is in the foreground",
		"tray_tip":        "midiflute (click to show)",
		"tray_show":       "Show",
		"tray_quit":       "Quit",
		"already_running": "midiflute.exe is already running.\n\nClick its icon in the system tray (bottom right of the screen) to open it.",
		"ui_error":        "Could not create the window: ",
		"choose_tos":      "Select the Tree of Savior installation folder",
		"not_tos":         "This folder does not contain release and data.\nPlease select the Tree of Savior (Japanese Ver.) folder.",
		"midi_connected":  "Connected: ",
		"midi_waiting":    "Waiting (connects automatically when a device is plugged in)",
		"midi_off":        "Disconnected (other software can use the keyboard)",
		"midi_busy":       "Could not open this device (another program may be using it)",
		"midi_notfound":   "Device not found",
		"midi_removed":    "The device was unplugged",
		"oct3_open":       "unlocked",
		"oct3_locked":     "locked",
		"tos_foreign":     "⚠ Another midiflute.exe (possibly another version or copy) seems to be running. Keep only one.",
		"tos_nofolder":    "The ToS folder is not set. Choose it with \"Change…\".",
		"tos_unlinked":    "Not linked (start ToS and log in with a character)",
		"tos_linked":      "Linked (not in playing mode) / 3rd octave: ",
		"tos_playing":     "Playing mode ♪ / 3rd octave: ",
		"lat_fmt":         "Last %.0f ms / Median %.0f ms",
		"lat_none":        "- (shown after you play)",
		"base_fmt":        "1st octave C: %s (key number %d)  ",
		"hint_learning":   "▶ Press the key you want to use as C of the 1st octave",
		"hint_normal":     "If you change the octave on your keyboard, use \"Set by key\" again.",
		"do":              "Do",
	},
	"ko": {
		"title":           "midiflute - ToS 피리 연주용 MIDI 수신",
		"grp_conn":        "연결",
		"midi_device":     "MIDI 기기",
		"connect":         "연결",
		"disconnect":      "연결 해제",
		"refresh":         "목록 새로고침",
		"tos_folder":      "ToS 폴더",
		"change":          "변경…",
		"language":        "言語 / Language",
		"grp_status":      "상태",
		"latency":         "지연:",
		"grp_keys":        "건반 배치",
		"oct_down":        "-1 옥타브",
		"semi_down":       "-반음",
		"semi_up":         "+반음",
		"oct_up":          "+1 옥타브",
		"learn":           "건반으로 지정",
		"legend":          "파랑=1옥타브  초록=2옥타브  주황=3옥타브  회색=소리 나지 않는 건반(범위 밖・미해금)  빨강=누르고 있는 건반",
		"fg_only":         "ToS가 앞에 있을 때만 연주 보내기",
		"tray_tip":        "midiflute (클릭하면 표시)",
		"tray_show":       "표시",
		"tray_quit":       "종료",
		"already_running": "midiflute.exe가 이미 실행 중입니다.\n\n작업 표시줄 오른쪽 아래의 트레이 아이콘을 클릭하면 실행 중인 화면을 열 수 있습니다.",
		"ui_error":        "화면을 만들 수 없습니다: ",
		"choose_tos":      "Tree of Savior 설치 폴더를 선택하세요",
		"not_tos":         "이 폴더에는 release와 data가 없습니다.\nTree of Savior (Japanese Ver.) 폴더를 선택하세요.",
		"midi_connected":  "연결됨: ",
		"midi_waiting":    "대기 중 (기기를 연결하면 자동으로 연결됩니다)",
		"midi_off":        "연결 해제됨 (다른 프로그램에서 키보드를 사용할 수 있습니다)",
		"midi_busy":       "이 기기를 열 수 없습니다 (다른 프로그램이 사용 중일 수 있습니다)",
		"midi_notfound":   "기기를 찾을 수 없습니다",
		"midi_removed":    "기기가 분리되었습니다",
		"oct3_open":       "해금됨",
		"oct3_locked":     "미해금",
		"tos_foreign":     "⚠ 다른 midiflute.exe(다른 버전이나 다른 위치의 것 포함)가 동시에 실행 중인 것 같습니다. 하나만 실행하세요.",
		"tos_nofolder":    "ToS 폴더가 설정되지 않았습니다. \"변경…\"에서 선택하세요.",
		"tos_unlinked":    "연동 안 됨 (ToS를 실행하고 캐릭터로 로그인하면 연동됩니다)",
		"tos_linked":      "연동 중 (연주 모드 아님) / 3옥타브: ",
		"tos_playing":     "연주 모드 중 ♪ / 3옥타브: ",
		"lat_fmt":         "최근 %.0f ms / 중앙값 %.0f ms",
		"lat_none":        "- (연주하면 표시됩니다)",
		"base_fmt":        "1옥타브의 도: %s (건반 번호 %d)  ",
		"hint_learning":   "▶ 1옥타브의 \"도\"로 할 건반을 누르세요",
		"hint_normal":     "키보드의 옥타브 전환을 바꿨다면 \"건반으로 지정\"으로 다시 맞춰 주세요.",
		"do":              "도",
	},
}
