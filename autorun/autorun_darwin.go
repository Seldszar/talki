package autorun

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

const jobTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
  <dict>
    <key>Label</key>
    <string>{{.Name}}</string>
    <key>ProgramArguments</key>
		<array>
			<string>{{.Executable}}</string>
		</array>
    <key>RunAtLoad</key>
    <true/>
    <key>AbandonProcessGroup</key>
    <true/>
  </dict>
</plist>`

var launchDir = filepath.Join(os.Getenv("HOME"), "Library", "LaunchAgents")

func (a *AutoRun) path() string {
	return filepath.Join(launchDir, fmt.Sprintf("%s.plist", a.Name))
}

func (a *AutoRun) Enable() error {
	t := template.Must(template.New("job").Parse(jobTemplate))

	if err := os.MkdirAll(launchDir, 0777); err != nil {
		return err
	}

	f, err := os.Create(a.path())

	if err != nil {
		return err
	}

	defer f.Close()

	if err := t.Execute(f, a); err != nil {
		return err
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
