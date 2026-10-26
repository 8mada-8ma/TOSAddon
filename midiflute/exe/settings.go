package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type Settings struct {
	TosDir         string `json:"tos_dir"`
	Device         string `json:"device"`
	Base           int    `json:"base"`
	ForegroundOnly bool   `json:"foreground_only"`
	AutoConnect    bool   `json:"auto_connect"`
	Lang           string `json:"lang"`
}

func exeDir() string {
	p, err := os.Executable()
	if err != nil {
		return "."
	}
	return filepath.Dir(p)
}

// 設定は %APPDATA%\midiflute\settings.json に保存する(exe を置く場所の書き込み権限に左右されない)
func settingsPath() string {
	if d, err := os.UserConfigDir(); err == nil {
		return filepath.Join(d, "midiflute", "settings.json")
	}
	return filepath.Join(exeDir(), "midiflute_exe.json")
}

func LoadSettings() Settings {
	s := Settings{Base: 60, AutoConnect: true}
	data, err := os.ReadFile(settingsPath())
	if err != nil {
		// 旧版(exe と同じフォルダに保存していた版)からの引き継ぎ
		data, err = os.ReadFile(filepath.Join(exeDir(), "midiflute_exe.json"))
	}
	if err == nil {
		_ = json.Unmarshal(data, &s)
	}
	if s.Base < 0 || s.Base > 127 {
		s.Base = 60
	}
	return s
}

func (s Settings) Save() {
	data, _ := json.MarshalIndent(s, "", "  ")
	_ = os.MkdirAll(filepath.Dir(settingsPath()), 0755)
	_ = os.WriteFile(settingsPath(), data, 0644)
}

// IsTosDir はToSのインストールフォルダらしいか
func IsTosDir(dir string) bool {
	st, err := os.Stat(filepath.Join(dir, "release"))
	if err != nil || !st.IsDir() {
		return false
	}
	st, err = os.Stat(filepath.Join(dir, "data"))
	return err == nil && st.IsDir()
}

// FindTosDir は、exeの置き場所・設定・よくある場所の順にToSフォルダを探す
func FindTosDir(s Settings) string {
	d := exeDir()
	// ...\addons\midiflute\midiflute.exe に置かれている場合
	if strings.EqualFold(filepath.Base(d), "midiflute") && strings.EqualFold(filepath.Base(filepath.Dir(d)), "addons") {
		if root := filepath.Dir(filepath.Dir(d)); IsTosDir(root) {
			return root
		}
	}
	if s.TosDir != "" && IsTosDir(s.TosDir) {
		return s.TosDir
	}
	// Steam の登録情報とライブラリフォルダ(D:やF:などの追加ライブラリを含む)から探す
	return findTosInSteam(append(steamDirsFromRegistry(),
		`C:\Program Files (x86)\Steam`, `C:\Program Files\Steam`))
}

func AddonDir(tosDir string) string {
	d := filepath.Join(tosDir, "addons", "midiflute")
	_ = os.MkdirAll(d, 0755)
	return d
}
