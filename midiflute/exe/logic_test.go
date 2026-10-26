package main

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestMapping(t *testing.T) {
	ids := map[int]bool{}
	names := []string{"C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B"}
	for off := 0; off < 36; off++ {
		id, kind := noteToID(60+off, 60, true)
		if id == 0 || ids[id] {
			t.Fatalf("off %d id %d", off, id)
		}
		ids[id] = true
		oct := off / 12
		sharp := strings.Contains(names[off%12], "#")
		if sharp != (id >= 15 && id <= 24 || id >= 32) || kind != KeyOct1+oct {
			t.Fatalf("off %d id %d kind %d", off, id, kind)
		}
	}
	if len(ids) != 36 {
		t.Fatal("36音でない")
	}
	if id, k := noteToID(84, 60, false); id != 0 || k != KeyLocked {
		t.Fatal("未解放の3オクターブ目が鳴る")
	}
	if id, k := noteToID(59, 60, true); id != 0 || k != KeyOut {
		t.Fatal("範囲外")
	}
	if id, _ := noteToID(96, 60, true); id != 0 {
		t.Fatal("範囲外(上)")
	}
	if noteName(60) != "C4" || noteName(48) != "C3" || noteName(61) != "C#4" || noteName(0) != "C-1" {
		t.Fatal(noteName(60), noteName(0))
	}
	if !isBlack(61) || isBlack(60) {
		t.Fatal("黒鍵判定")
	}
}

func TestBridge(t *testing.T) {
	dir := t.TempDir()
	b := NewBridge(dir)
	b.Send('P', 1)
	b.Send('S', 1)
	data, _ := os.ReadFile(filepath.Join(dir, "events.txt"))
	want := "# " + b.session + ";\n1 P 1;\n2 S 1;\n"
	if string(data) != want {
		t.Fatalf("events.txt=%q", data)
	}
	os.WriteFile(filepath.Join(dir, "ack.txt"), []byte(b.session+" 2;"), 0644)
	if b.Poll() {
		t.Fatal("未処理が残っている")
	}
	if _, _, ok := b.Latency(); !ok {
		t.Fatal("遅延が記録されない")
	}
	b.Send('P', 5)
	data, _ = os.ReadFile(filepath.Join(dir, "events.txt"))
	if string(data) != "# "+b.session+";\n3 P 5;\n" {
		t.Fatalf("処理済みが消えない %q", data)
	}
	// 別のプログラムが書いたら検知、自分が書けば解消
	os.WriteFile(filepath.Join(dir, "events.txt"), []byte("# 123;\n5 N 60;\n"), 0644)
	if !b.ForeignWriter() {
		t.Fatal("他の書き込みを検知しない")
	}
	b.Send('S', 5)
	if b.ForeignWriter() {
		t.Fatal("自分の書き込みを他と誤認")
	}
	// 古い未処理は捨てる
	b.mu.Lock()
	for i := range b.pending {
		b.pending[i].at = time.Now().Add(-2 * time.Second)
	}
	b.mu.Unlock()
	b.Send('P', 7)
	data, _ = os.ReadFile(filepath.Join(dir, "events.txt"))
	if strings.Contains(string(data), " P 5;") || !strings.Contains(string(data), " P 7;") {
		t.Fatalf("古いイベントが残る %q", data)
	}
	// 書きかけの ack は無視
	os.WriteFile(filepath.Join(dir, "ack.txt"), []byte(b.session+" 3"), 0644)
	if b.readAck() != 0 {
		t.Fatal("不完全なackを読んだ")
	}
	// 状態ファイル
	now := strconv.FormatInt(time.Now().Unix(), 10)
	os.WriteFile(filepath.Join(dir, "status.txt"), []byte(now+" 1 0;"), 0644)
	if s := ReadStatus(dir); !s.Linked || !s.Playing || s.Oct3 {
		t.Fatalf("%+v", s)
	}
	old := strconv.FormatInt(time.Now().Unix()-30, 10)
	os.WriteFile(filepath.Join(dir, "status.txt"), []byte(old+" 1 1;"), 0644)
	if s := ReadStatus(dir); s.Linked {
		t.Fatal("古い状態を連携中と判断")
	}
}

func TestSteamSearch(t *testing.T) {
	root := t.TempDir()
	steam := filepath.Join(root, "Steam")
	lib := filepath.Join(root, "lib2")
	os.MkdirAll(filepath.Join(steam, "steamapps"), 0755)
	vdf := "\"libraryfolders\"\n{\n\t\"0\"\n\t{\n\t\t\"path\"\t\t\"" + strings.ReplaceAll(steam, `\`, `\\`) + "\"\n\t}\n\t\"1\"\n\t{\n\t\t\"path\"\t\t\"" + strings.ReplaceAll(lib, `\`, `\\`) + "\"\n\t\t\"apps\"\n\t\t{\n\t\t\t\"1175730\"\t\t\"123\"\n\t\t}\n\t}\n}\n"
	os.WriteFile(filepath.Join(steam, "steamapps", "libraryfolders.vdf"), []byte(vdf), 0644)
	if findTosInSteam([]string{steam}) != "" {
		t.Fatal("無いのに見つけた")
	}
	tos := filepath.Join(lib, "steamapps", "common", "Tree of Savior (Japanese Ver.)")
	os.MkdirAll(filepath.Join(tos, "release"), 0755)
	os.MkdirAll(filepath.Join(tos, "data"), 0755)
	os.WriteFile(filepath.Join(lib, "steamapps", "appmanifest_1175730.acf"), []byte("\"AppState\"\n{\n\t\"installdir\"\t\t\"Tree of Savior (Japanese Ver.)\"\n}\n"), 0644)
	if got := findTosInSteam([]string{"", steam}); got != tos {
		t.Fatalf("got %q", got)
	}
}

func TestTextsComplete(t *testing.T) {
	for _, code := range langCodes {
		for k := range texts["ja"] {
			if texts[code][k] == "" {
				t.Errorf("%s に %s がない", code, k)
			}
		}
		if len(texts[code]) != len(texts["ja"]) {
			t.Errorf("%s の項目数が違う", code)
		}
	}
	setLang("ko")
	if T("connect") != "연결" {
		t.Fatal("切替できない")
	}
	setLang("xx")
	if curLang != "ja" {
		t.Fatal("不明な言語")
	}
}
