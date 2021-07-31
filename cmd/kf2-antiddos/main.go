package main

import (
	"kf2-antiddos/internal/action"
	"kf2-antiddos/internal/common"
	"kf2-antiddos/internal/config"
	"kf2-antiddos/internal/history"
	"kf2-antiddos/internal/output"
	"kf2-antiddos/internal/parser"
	"kf2-antiddos/internal/reader"

	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
)

const (
	ExitSuccess  int = 0
	ExitArgError int = 1
)

const (
	AppName = "kf2-antiddos"
)

var (
	AppVersion string = "dev"
)

func main() {
	cfg := parseArgs()

	switch {
	case cfg.ShowHelp:
		printHelp()
		os.Exit(ExitSuccess)
	case cfg.ShowVersion:
		printVersion()
		os.Exit(ExitSuccess)
	}

	if cfg.IsValid() {
		cfg.SetEmptyArgs()
	} else {
		os.Exit(ExitArgError)
	}

	switch cfg.OutputMode {
	case config.OT_All:
		output.AllMode()
	case config.OT_Proxy:
		output.ProxyMode()
	case config.OT_Quiet:
		output.QuietMode()
	case config.OT_Self:
		output.SelfMode()
	}

	runtime.GOMAXPROCS(int(cfg.Jobs))

	Workers := make([]common.Worker, 0, cfg.Jobs+3)

	wg := sync.WaitGroup{}

	// Data flow:
	banChan := make(chan string, cfg.Jobs)
	inputChan := make(chan common.RawEvent, cfg.Jobs)
	eventChan := make(chan common.Event, cfg.Jobs)
	resetChan := make(chan string, cfg.Jobs)

	// Reader worker
	Workers = append(Workers,
		reader.New(
			uint(len(Workers)),
			&inputChan,
		))

	// parser workers
	for i := uint(0); i < cfg.Jobs; i++ {
		Workers = append(Workers,
			parser.New(
				uint(len(Workers)),
				&inputChan,
				&eventChan,
			))
	}

	// History worker
	Workers = append(Workers,
		history.New(
			uint(len(Workers)),
			&eventChan,
			&banChan,
			&resetChan,
			cfg.MaxConn,
		))

	// Action worker
	Workers = append(Workers,
		action.New(
			uint(len(Workers)),
			cfg.DenyTime,
			cfg.Shell,
			cfg.AllowAction,
			cfg.DenyAction,
			&banChan,
			&resetChan,
		))

	wg.Add(len(Workers))

	closeHandler(Workers, &wg)

	for i := range Workers {
		Workers[i].Do()
	}

	output.Println("started")

	wg.Wait()

	output.Println("exit")

	os.Exit(ExitSuccess)
}

func closeHandler(Workers []common.Worker, wg *sync.WaitGroup) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	go func() {
		<-interrupt
		output.Println("interrupt")
		for _, worker := range Workers {
			worker.Stop()
			wg.Done()
		}
	}()
}
