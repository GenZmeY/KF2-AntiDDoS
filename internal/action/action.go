package action

import (
	"kf2-antiddos/internal/output"

	"os"
	"os/exec"
	"strings"
	"time"
)

type Action struct {
	ticker      *time.Ticker
	ips         map[string]bool // map[IP]readyToUnban
	allowAction string
	denyAction  string
	shell       string
	quit        chan struct{}
	banChan     *chan string
	resetChan   *chan string
	workerID    uint
}

func New(workerID uint, denyTime uint, shell, allowAction, denyAction string, banChan, resetChan *chan string) *Action {
	return &Action{
		ticker:      time.NewTicker(time.Duration(denyTime) * time.Second),
		ips:         make(map[string]bool),
		allowAction: allowAction,
		denyAction:  denyAction,
		shell:       shell,
		quit:        make(chan struct{}),
		banChan:     banChan,
		resetChan:   resetChan,
		workerID:    workerID,
	}
}

func (a *Action) Do() {
	go func() {
		for {
			select {
			case ip := <-*a.banChan:
				a.deny(ip)
			case <-a.ticker.C:
				a.allow(false)
			case <-a.quit:
				a.ticker.Stop()
				a.allow(true)
				return
			}
		}
	}()
}

func (a *Action) Stop() {
	close(a.quit)
}

func (a *Action) allow(unbanAll bool) {
	unban := make([]string, 0)

	for ip := range a.ips {
		if unbanAll || a.ips[ip] { // aka if readyToUnban
			unban = append(unban, ip)
		} else {
			a.ips[ip] = true // mark readyToUnban next time
		}
	}

	for _, ip := range unban {
		delete(a.ips, ip)
	}

	if len(unban) != 0 {
		for _, ip := range unban {
			*a.resetChan <- ip
		}
		output.Printf("Allow: %s", strings.Join(unban, ", "))

		if err := a.execCmd(a.allowAction, unban); err != nil {
			output.Error(err.Error())
			return
		}
	}
}

func (a *Action) deny(ip string) {
	a.ips[ip] = false

	output.Printf("Ban: %s", ip)

	if err := a.execCmd(a.denyAction, []string{ip}); err != nil {
		output.Error(err.Error())
		return
	}
}

func (a *Action) execCmd(command string, args []string) error {
	WorkingDir, err := os.Getwd()
	if err != nil {
		WorkingDir = ""
	}
	cmd := &exec.Cmd{
		Path:   a.shell,
		Args:   append([]string{a.shell, command}, args...),
		Stdout: output.StdoutWriter(),
		Stderr: output.StderrWriter(),
		Dir:    WorkingDir,
	}

	if err := cmd.Start(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}
