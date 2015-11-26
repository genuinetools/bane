package apparmor

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"text/template"
)

// ProfileConfig defines the config for an
// apparmor profile to be generated from.
type ProfileConfig struct {
	Name         string
	Filesystem   FsConfig
	Network      NetConfig
	Capabilities CapConfig

	Imports      []string
	InnerImports []string
}

// FsConfig defines the filesystem options for a profile.
type FsConfig struct {
	ReadOnlyPaths   []string
	LogOnWritePaths []string
	WritablePaths   []string
	AllowExec       []string
	DenyExec        []string
}

// NetConfig defines the network options for a profile.
// For example you probably don't need NetworkRaw if your
// application doesn't `ping`.
// Currently limited to AppArmor 2.3-2.6 rules.
type NetConfig struct {
	Raw       bool
	Packet    bool
	Protocols []string
}

// CapConfig defines the allowed or denied kernel capabilities
// for a profile.
type CapConfig struct {
	Allow []string
	Deny  []string
}

// Generate uses the baseTemplate to generate an apparmor profile
// for the ProfileConfig passed.
func (profile *ProfileConfig) Generate(out io.Writer) error {
	compiled, err := template.New("apparmor_profile").Parse(baseTemplate)
	if err != nil {
		return err
	}
	if tunablesExists("global") {
		profile.Imports = append(profile.Imports, "#include <tunables/global>")
	} else {
		profile.Imports = append(profile.Imports, "@{PROC}=/proc/")
	}
	if abstractionsExists("base") {
		profile.InnerImports = append(profile.InnerImports, "#include <abstractions/base>")
	}
	if err := compiled.Execute(out, profile); err != nil {
		return err
	}
	return nil
}

// check if the tunables/global exist
func tunablesExists(name string) bool {
	_, err := os.Stat(path.Join("/etc/apparmor.d/tunables", name))
	return err == nil
}

// check if abstractions/base exist
func abstractionsExists(name string) bool {
	_, err := os.Stat(path.Join("/etc/apparmor.d/abstractions", name))
	return err == nil
}

// Install takes a profile config, generates the profile
// and installs it in the given directory with `apparmor_parser`.
func (profile *ProfileConfig) Install(dir string) error {
	// Make sure the path where they want to save the profile exists
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	f, err := os.OpenFile(filepath.Join(dir, profile.Name), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	if err := profile.Generate(f); err != nil {
		f.Close()
		return err
	}
	f.Close()

	cmd := exec.Command("/sbin/apparmor_parser", "-r", "-W", profile.Name)
	// to use the parser directly we have to make sure we are in the correct
	// dir with the profile
	cmd.Dir = dir

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Error loading apparmor profile %s: %v (%s)", profile.Name, err, output)
	}
	return nil
}
