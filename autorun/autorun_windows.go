//go:build windows
// +build windows

package autorun

import (
	"golang.org/x/sys/windows/registry"
)

func (a *AutoRun) Enable() error {
	key, err := registry.OpenKey(registry.CURRENT_USER, `SOFTWARE\Microsoft\Windows\CurrentVersion\Run`, registry.ALL_ACCESS)

	if err != nil {
		return err
	}

	defer key.Close()

	if err = key.SetStringValue(a.Name, a.Executable); err != nil {
		return err
	}

	return nil
}

func (a *AutoRun) Disable() error {
	key, err := registry.OpenKey(registry.CURRENT_USER, `SOFTWARE\Microsoft\Windows\CurrentVersion\Run`, registry.ALL_ACCESS)

	if err != nil {
		return err
	}

	defer key.Close()

	if err = key.DeleteValue(a.Name); err != nil {
		return err
	}

	return nil
}

func (a *AutoRun) IsEnabled() bool {
	key, err := registry.OpenKey(registry.CURRENT_USER, `SOFTWARE\Microsoft\Windows\CurrentVersion\Run`, registry.READ)

	if err != nil {
		return false
	}

	defer key.Close()

	if _, _, err := key.GetStringValue(a.Name); err != nil {
		return false
	}

	return true
}
