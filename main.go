package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/Sirupsen/logrus"
	"github.com/jessfraz/bane/apparmor"
	"github.com/jessfraz/bane/version"
)

const (
	// BANNER is what is printed for help/info output
	BANNER = ` _
| |__   __ _ _ __   ___
| '_ \ / _` + "`" + ` | '_ \ / _ \
| |_) | (_| | | | |  __/
|_.__/ \__,_|_| |_|\___|
 Custom AppArmor profile generator for docker containers
 Version: %s

`
)

var (
	apparmorProfileDir string

	debug bool
	vrsn  bool
)

func init() {
	// parse flags
	flag.StringVar(&apparmorProfileDir, "profile-dir", "/etc/apparmor.d/containers", "directory for saving the profiles")

	flag.BoolVar(&vrsn, "version", false, "print version and exit")
	flag.BoolVar(&vrsn, "v", false, "print version and exit (shorthand)")
	flag.BoolVar(&debug, "d", false, "run in debug mode")

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(BANNER, version.VERSION))
		flag.PrintDefaults()
	}

	flag.Parse()

	if vrsn {
		fmt.Printf("bane version %s, build %s", version.VERSION, version.GITCOMMIT)
		os.Exit(0)
	}

	if flag.NArg() < 1 {
		usageAndExit("Pass the path to the config file.", 1)
	}

	// parse the arg
	arg := flag.Args()[0]

	if arg == "help" {
		usageAndExit("", 0)
	}

	if arg == "version" {
		fmt.Printf("bane version %s, build %s", version.VERSION, version.GITCOMMIT)
		os.Exit(0)
	}

	// set log level
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
}

func main() {
	// parse the profile config file
	profileConfig := flag.Args()[0]

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
}

func usageAndExit(message string, exitCode int) {
	if message != "" {
		fmt.Fprintf(os.Stderr, message)
		fmt.Fprintf(os.Stderr, "\n\n")
	}
	flag.Usage()
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(exitCode)
}
