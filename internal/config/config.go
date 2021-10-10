package config

import (
	"fmt"
	"os"
	"runtime"

	"kf2-antiddos/internal/output"
)

const (
	OT_Proxy = "proxy"
	OT_Self  = "self"
	OT_All   = "all"
	OT_Quiet = "quiet"
)

type Config struct {
	Shell       string
	DenyAction  string
	AllowAction string
	Jobs        uint
	OutputMode  string
	DenyTime    uint
	AllowTime   uint
	MaxConn     uint

	ShowVersion bool
	ShowHelp    bool
}

func (cfg Config) IsValid() bool {
	errs := make([]string, 0)

	if cfg.Shell == "" {
		errs = append(errs, "shell can not be empty")
	} else if _, err := os.Stat(cfg.Shell); os.IsNotExist(err) {
		errs = append(errs, fmt.Sprintf("shell %s not found", cfg.Shell))
	}

	if cfg.AllowAction == "" {
		errs = append(errs, "allow_action can not be empty")
	} else if _, err := os.Stat(cfg.AllowAction); os.IsNotExist(err) {
		errs = append(errs, fmt.Sprintf("allow_action file %s not found", cfg.AllowAction))
	}

	if cfg.DenyAction == "" {
		errs = append(errs, "deny_action can not be empty")
	} else if _, err := os.Stat(cfg.DenyAction); os.IsNotExist(err) {
		errs = append(errs, fmt.Sprintf("deny_action file %s not found", cfg.DenyAction))
	}

	switch cfg.OutputMode {
	case OT_Proxy:
	case OT_Self:
	case OT_All:
	case OT_Quiet:
	case "":
	default:
		errs = append(errs, fmt.Sprintf("Unknown output_type: %s", cfg.OutputMode))
	}

	for _, err := range errs {
		output.Errorln(err)
	}

	return len(errs) == 0
}

func (cfg *Config) SetEmptyArgs() {
	if cfg.Jobs == 0 {
		cfg.Jobs = uint(runtime.NumCPU())
	}
	if cfg.MaxConn == 0 {
		cfg.MaxConn = 10
	}
	if cfg.OutputMode == "" {
		cfg.OutputMode = OT_Self
	}
	if cfg.DenyTime == 0 {
		cfg.DenyTime = 20 * 60
	}
	if cfg.AllowTime == 0 {
		cfg.AllowTime = 20 * 60
	}
}
