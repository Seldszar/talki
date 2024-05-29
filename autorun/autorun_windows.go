//go:build windows
// +build windows

package autorun

// #cgo LDFLAGS: -lole32 -luuid
/*
#define WIN32_LEAN_AND_MEAN

#include <stdint.h>
#include <windows.h>

char CreateShortcut(char *shortcutA, char *path, char *args);
*/
import "C"

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var startupDir = filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Roaming", "Microsoft", "Windows", "Start Menu", "Programs", "Startup")

func (a *AutoRun) path() string {
	return filepath.Join(startupDir, fmt.Sprintf("%s.lnk", a.Name))
}

func (a *AutoRun) Enable() error {
	if err := os.MkdirAll(startupDir, 0777); err != nil {
		return err
	}

	if res := C.CreateShortcut(C.CString(a.path()), C.CString(a.Executable), C.CString("")); res != 0 {
		return errors.New("unable to create shortcut")
	}

	return nil
}

func (a *AutoRun) Disable() error {
	return os.Remove(a.path())
}

func (a *AutoRun) IsEnabled() bool {
	if _, err := os.Stat(a.path()); err != nil {
		return false
	}

	return true
}
