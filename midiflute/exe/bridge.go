package main

// ToSアドオン midiflute とのファイル連携
//   events.txt  exe → アドオン  "# <セッション>;" の後に "<連番> <P|S|T> <音番号>;"
//   ack.txt     アドオン → exe  "<セッション> <処理済みの連番>;"
//   helper.txt  exe → アドオン  exe が動いている証拠(UNIX時刻)
//   status.txt  アドオン → exe  "<UNIX時刻> <演奏モード 0/1> <3オクターブ目 0/1>;"

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type event struct {
	seq  int
	kind byte
	val  int
	at   time.Time
}

// 古くなった未処理イベントは送り直さない(遅れて鳴るより、鳴らない方がよい)
const staleAfter = 1500 * time.Millisecond

type Bridge struct {
	dir       string
	session   string
	mu        sync.Mutex
	seq       int
	pending   []event
	sentAt    map[int]time.Time
	latencies []float64
}

func NewBridge(dir string) *Bridge {
	b := &Bridge{dir: dir, session: strconv.FormatInt(time.Now().Unix(), 10), sentAt: map[int]time.Time{}}
	b.mu.Lock()
	b.flushLocked()
	b.mu.Unlock()
	return b
}

func (b *Bridge) path(name string) string { return filepath.Join(b.dir, name) }

func (b *Bridge) Send(kind byte, val int) int {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.seq++
	b.pending = append(b.pending, event{b.seq, kind, val, time.Now()})
	b.sentAt[b.seq] = time.Now()
	b.flushLocked()
	return b.seq
}

func (b *Bridge) readAck() int {
	data, err := os.ReadFile(b.path("ack.txt"))
	if err != nil {
		return 0
	}
	s := string(data)
	if !strings.HasSuffix(s, ";") {
		return 0
	}
	parts := strings.Fields(strings.TrimSuffix(s, ";"))
	if len(parts) != 2 || parts[0] != b.session {
		return 0
	}
	n, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0
	}
	return n
}

func (b *Bridge) flushLocked() {
	acked := b.readAck()
	b.consumeAckLocked(acked)
	if len(b.pending) > 200 {
		b.pending = b.pending[len(b.pending)-200:]
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "# %s;\n", b.session)
	for _, e := range b.pending {
		fmt.Fprintf(&sb, "%d %c %d;\n", e.seq, e.kind, e.val)
	}
	data := []byte(sb.String())
	for i := 0; i < 20; i++ {
		if err := os.WriteFile(b.path("events.txt"), data, 0644); err == nil {
			return
		}
		time.Sleep(time.Millisecond)
	}
}

func (b *Bridge) consumeAckLocked(acked int) {
	now := time.Now()
	keep := b.pending[:0]
	for _, e := range b.pending {
		if e.seq > acked && now.Sub(e.at) < staleAfter {
			keep = append(keep, e)
		}
	}
	b.pending = keep
	for seq, t := range b.sentAt {
		if seq <= acked {
			b.latencies = append(b.latencies, float64(now.Sub(t).Microseconds())/1000)
			delete(b.sentAt, seq)
		} else if now.Sub(t) > 5*time.Second {
			delete(b.sentAt, seq) // 演奏モード外などで処理されなかったもの
		}
	}
	if len(b.latencies) > 20 {
		b.latencies = b.latencies[len(b.latencies)-20:]
	}
}

// Poll は処理済みの連番を確認して遅延を記録する。未処理があるかを返す。
func (b *Bridge) Poll() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	if len(b.sentAt) == 0 {
		return false
	}
	b.consumeAckLocked(b.readAck())
	return len(b.sentAt) > 0
}

// Latency は直近の遅延と中央値(ミリ秒)。記録が無ければ ok=false。
func (b *Bridge) Latency() (last, median float64, ok bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if len(b.latencies) == 0 {
		return 0, 0, false
	}
	s := append([]float64(nil), b.latencies...)
	sort.Float64s(s)
	return b.latencies[len(b.latencies)-1], s[len(s)/2], true
}

// ForeignWriter は、events.txt を別のプログラムが書き換えていれば true
// (古い midiflute.py などが同時に動いていると起きる)
func (b *Bridge) ForeignWriter() bool {
	data, err := os.ReadFile(b.path("events.txt"))
	if err != nil {
		return false
	}
	line := strings.SplitN(string(data), "\n", 2)[0]
	if !strings.HasPrefix(line, "# ") || !strings.HasSuffix(line, ";") {
		return false
	}
	return strings.TrimSuffix(strings.TrimPrefix(line, "# "), ";") != b.session
}

func (b *Bridge) Heartbeat() {
	_ = os.WriteFile(b.path("helper.txt"), []byte(strconv.FormatInt(time.Now().Unix(), 10)), 0644)
}

type GameStatus struct {
	Linked  bool // アドオンが動いている(ToSでキャラクターがマップにいる)
	Playing bool // 演奏モード中
	Oct3    bool // 3オクターブ目が解放済み
}

func ReadStatus(dir string) GameStatus {
	data, err := os.ReadFile(filepath.Join(dir, "status.txt"))
	if err != nil {
		return GameStatus{}
	}
	s := string(data)
	if !strings.HasSuffix(strings.TrimSpace(s), ";") {
		return GameStatus{}
	}
	f := strings.Fields(strings.TrimSuffix(strings.TrimSpace(s), ";"))
	if len(f) != 3 {
		return GameStatus{}
	}
	t, err := strconv.ParseInt(f[0], 10, 64)
	if err != nil {
		return GameStatus{}
	}
	age := time.Now().Unix() - t
	if age < -5 || age > 5 {
		return GameStatus{}
	}
	return GameStatus{Linked: true, Playing: f[1] == "1", Oct3: f[2] == "1"}
}
