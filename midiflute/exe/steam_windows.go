//go:build windows

package main

import "golang.org/x/sys/windows/registry"

func steamDirsFromRegistry() []string {
	var dirs []string
	read := func(root registry.Key, path, name string) {
		k, err := registry.OpenKey(root, path, registry.QUERY_VALUE)
		if err != nil {
			return
		}
		defer k.Close()
		if v, _, err := k.GetStringValue(name); err == nil && v != "" {
			dirs = append(dirs, v)
		}
	}
	read(registry.CURRENT_USER, `Software\Valve\Steam`, "SteamPath")
	read(registry.LOCAL_MACHINE, `SOFTWARE\WOW6432Node\Valve\Steam`, "InstallPath")
	read(registry.LOCAL_MACHINE, `SOFTWARE\Valve\Steam`, "InstallPath")
	return dirs
}
