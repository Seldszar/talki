//go:build linux
// +build linux

package autorun

import (
	"fmt"
	"os"
	"text/template"
)

const desktopTemplate = `[Desktop Entry]
Type=Application
Name={{.DisplayName}}
Exec={{.Executable}}
X-GNOME-Autostart-enabled=true
`

func (a *AutoRun) path() string {
	return fmt.Sprintf("~/.config/autostart/%s.desktop", a.Name)
}

func (a *AutoRun) Enable() error {
	t, err := template.New("desktop").Parse(desktopTemplate)

	if err != nil {
		return err
	}

	f, err := os.Create(a.path())

	if err != nil {
		return err
	}

	defer f.Close()

	if err = t.Execute(f, a); err != nil {
		return err
	}

	return nil
}

func (a *AutoRun) Disable() error {
	if err := os.Remove(a.path()); err != nil {
		return err
	}

	return nil
}

func (a *AutoRun) IsEnabled() bool {
	if _, err := os.Stat(a.path()); err != nil {
		return false
	}

	return true
}
