package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/genuinetools/bane/apparmor"
	"github.com/genuinetools/bane/version"
	"github.com/genuinetools/pkg/cli"
	"github.com/sirupsen/logrus"
)

var (
	apparmorProfileDir string

	debug bool
)

func main() {
	// Create a new cli program.
	p := cli.NewProgram()
	p.Name = "bane"
	p.Description = "Custom AppArmor profile generator for docker containers"

	// Set the GitCommit and Version.
	p.GitCommit = version.GITCOMMIT
	p.Version = version.VERSION

	// Setup the global flags.
	p.FlagSet = flag.NewFlagSet("global", flag.ExitOnError)
	p.FlagSet.StringVar(&apparmorProfileDir, "profile-dir", "/etc/apparmor.d/containers", "directory for saving the profiles")
	p.FlagSet.BoolVar(&debug, "d", false, "enable debug logging")

	// Set the before function.
	p.Before = func(ctx context.Context) error {
		// Set the log level.
		if debug {
			logrus.SetLevel(logrus.DebugLevel)
		}

		if p.FlagSet.NArg() < 1 {
			return fmt.Errorf("pass the path to the config file")
		}

		return nil
	}

	// Set the main program action.
	p.Action = func(ctx context.Context, args []string) error {
		// parse the profile config file
		profileConfig := args[0]

		// make sure the file exists
		if _, err := os.Stat(profileConfig); os.IsNotExist(err) {
			logrus.Fatalf("No such file or directory: %s", profileConfig)
		}

		file, err := ioutil.ReadFile(profileConfig)
		if err != nil {
			logrus.Fatalf("Reading file %s failed: %v", profileConfig, err)
		}

		var profile apparmor.ProfileConfig
		if _, err := toml.Decode(string(file), &profile); err != nil {
			logrus.Fatalf("Parsing config file %s failed: %q", profileConfig, err)
		}

		// clean the profile name so we are sure it starts with `docker-`
		if !strings.HasPrefix(profile.Name, "docker-") {
			profile.Name = fmt.Sprintf("docker-%s", profile.Name)
		}

		// install the profile
		if err := profile.Install(apparmorProfileDir); err != nil {
			logrus.Fatalf("Installing profile %s failed: %v", profile.Name, err)
		}

		fmt.Printf("Profile installed successfully you can now run the profile with\n`docker run --security-opt=\"apparmor:%s\"`\n", profile.Name)

		return nil
	}

	// Run our program.
	p.Run()
}
