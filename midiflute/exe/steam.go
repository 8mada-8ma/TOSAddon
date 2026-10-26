package main

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// 日本版 Tree of Savior の Steam アプリID(steam_appid.txt で確認)
const tosAppID = "1175730"

var vdfPathRe = regexp.MustCompile(`"path"\s+"((?:[^"\\]|\\.)*)"`)
var vdfInstallDirRe = regexp.MustCompile(`"installdir"\s+"((?:[^"\\]|\\.)*)"`)

func vdfUnescape(s string) string { return strings.ReplaceAll(s, `\\`, `\`) }

// steamLibraries は libraryfolders.vdf からライブラリフォルダの一覧を返す(Steam本体のフォルダを含む)
func steamLibraries(steamDir string) []string {
	libs := []string{steamDir}
	data, err := os.ReadFile(filepath.Join(steamDir, "steamapps", "libraryfolders.vdf"))
	if err != nil {
		return libs
	}
	for _, m := range vdfPathRe.FindAllStringSubmatch(string(data), -1) {
		p := vdfUnescape(m[1])
		dup := false
		for _, l := range libs {
			if strings.EqualFold(filepath.Clean(l), filepath.Clean(p)) {
				dup = true
			}
		}
		if !dup {
			libs = append(libs, p)
		}
	}
	return libs
}

// findTosInSteam は Steam のライブラリから ToS のインストールフォルダを探す
func findTosInSteam(steamDirs []string) string {
	for _, sd := range steamDirs {
		if sd == "" {
			continue
		}
		for _, lib := range steamLibraries(sd) {
			apps := filepath.Join(lib, "steamapps")
			if data, err := os.ReadFile(filepath.Join(apps, "appmanifest_"+tosAppID+".acf")); err == nil {
				if m := vdfInstallDirRe.FindStringSubmatch(string(data)); m != nil {
					if d := filepath.Join(apps, "common", vdfUnescape(m[1])); IsTosDir(d) {
						return d
					}
				}
			}
			if d := filepath.Join(apps, "common", "Tree of Savior (Japanese Ver.)"); IsTosDir(d) {
				return d
			}
		}
	}
	return ""
}
