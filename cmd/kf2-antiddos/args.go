package main

import (
	"github.com/juju/gnuflag"

	"kf2-antiddos/internal/config"
	"kf2-antiddos/internal/output"
)

func printHelp() {
	output.Println("Anti DDoS tool for kf2 servers")
	output.Println("")
	output.Printf("Usage: <kf2_logs_output> | %s [option]... <shell> <deny_script> <allow_script>", AppName)
	output.Println("kf2_logs_output            KF2 logs to redirect to stdin")
	output.Println("shell                      shell to run deny_script and allow_script")
	output.Println("deny_script                firewall deny script (takes IP as argument)")
	output.Println("allow_script               firewall allow script (takes IPs as arguments)")
	output.Println("")
	output.Println("Options:")
	output.Println("  -j, --jobs N             allow N jobs at once")
	output.Println("  -o, --output MODE        self|proxy|all|quiet")
	output.Println("  -t, --deny-time TIME     minimum ip deny TIME (seconds)")
	output.Println("  -c, --max-connections N  Skip N connections before run deny script")
	output.Println("  -v, --version            Show version")
	output.Println("  -h, --help               Show help")
}

func printVersion() {
	output.Printf("%s %s", AppName, AppVersion)
}

func parseArgs() config.Config {
	rawCfg := config.Config{}

	gnuflag.UintVar(&rawCfg.Jobs, "j", 0, "")
	gnuflag.UintVar(&rawCfg.Jobs, "jobs", 0, "")

	gnuflag.StringVar(&rawCfg.OutputMode, "o", "", "")
	gnuflag.StringVar(&rawCfg.OutputMode, "output", "", "")

	gnuflag.UintVar(&rawCfg.DenyTime, "t", 0, "")
	gnuflag.UintVar(&rawCfg.DenyTime, "deny-time", 0, "")

	gnuflag.UintVar(&rawCfg.MaxConn, "c", 0, "")
	gnuflag.UintVar(&rawCfg.MaxConn, "max-connections", 0, "")

	gnuflag.BoolVar(&rawCfg.ShowVersion, "v", false, "")
	gnuflag.BoolVar(&rawCfg.ShowVersion, "version", false, "")

	gnuflag.BoolVar(&rawCfg.ShowHelp, "h", false, "")
	gnuflag.BoolVar(&rawCfg.ShowHelp, "help", false, "")

	gnuflag.Parse(true)

	for i := 0; i < 3 && i < gnuflag.NArg(); i++ {
		switch i {
		case 0:
			rawCfg.Shell = gnuflag.Arg(i)
		case 1:
			rawCfg.DenyAction = gnuflag.Arg(i)
		case 2:
			rawCfg.AllowAction = gnuflag.Arg(i)
		}
	}

	return rawCfg
}
