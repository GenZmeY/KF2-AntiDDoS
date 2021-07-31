package parser

import (
	"kf2-antiddos/internal/common"
	"regexp"
)

const (
	ngConnectIP     = "ConnectIP"
	ngPlayerStartIP = "PlayerStartIP"
	ngPlayerEndIP   = "PlayerEndIP"
	rxIP            = `\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`
	rxConnect       = `NetComeGo:\sOpen\sTheWorld\s+(?P<` + ngConnectIP + `>` + rxIP + `){1}`
	rxPlayerStart   = `DevOnline:\sVerifyClientAuthSession:\sClientIP:\s(?P<` + ngPlayerStartIP + `>` + rxIP + `){1}`
	rxPlayerEnd     = `DevOnline:\sEndRemoteClientAuthSession:\sClientAddr:\s(?P<` + ngPlayerEndIP + `>` + rxIP + `){1}`
	rxValue         = rxConnect + `|` + rxPlayerStart + `|` + rxPlayerEnd
)

var (
	rxKFLog *regexp.Regexp = regexp.MustCompile(rxValue)
)

type Parser struct {
	quit       chan struct{}
	inputChan  *chan common.RawEvent
	outputChan *chan common.Event
	workerID   uint
}

func New(workerID uint, inputChan *chan common.RawEvent, outputChan *chan common.Event) *Parser {
	return &Parser{
		inputChan:  inputChan,
		outputChan: outputChan,
		quit:       make(chan struct{}),
		workerID:   workerID,
	}
}

func (p *Parser) Do() {
	go func() {
		for {
			select {
			case rawEvent := <-*p.inputChan:
				*p.outputChan <- p.parse(rawEvent)
			case <-p.quit:
				return
			}
		}
	}()
}

func (p *Parser) Stop() {
	close(p.quit)
}

func (p *Parser) parse(rawEvent common.RawEvent) common.Event {
	res := common.Event{
		LineNum: rawEvent.LineNum,
	}

	match := rxKFLog.FindStringSubmatch(rawEvent.Text)
	for i, name := range rxKFLog.SubexpNames() {
		if i != 0 && name != "" && i <= len(match) && match[i] != "" {
			switch name {
			case ngConnectIP:
				res.ConnectIP = match[i]
			case ngPlayerStartIP:
				res.PlayerStartIP = match[i]
			case ngPlayerEndIP:
				res.PlayerEndIP = match[i]
			}
		}
	}
	return res
}
