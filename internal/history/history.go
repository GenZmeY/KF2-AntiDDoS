package history

import (
	"kf2-antiddos/internal/common"
)

type History struct {
	quit      chan struct{}
	eventChan *chan common.Event
	banChan   *chan string
	resetChan *chan string
	head      byte
	history   map[byte]common.Event
	ips       map[string]uint // map[ip]conn_count
	whitelist map[string]struct{}
	banned    map[string]struct{}
	maxConn   uint
	workerID  uint
}

func New(workerID uint, eventChan *chan common.Event, banChan *chan string, resetChan *chan string, maxConn uint) *History {
	return &History{
		quit:      make(chan struct{}),
		ips:       make(map[string]uint, 0),
		history:   make(map[byte]common.Event, 0),
		whitelist: make(map[string]struct{}, 0),
		banned:    make(map[string]struct{}, 0),
		eventChan: eventChan,
		banChan:   banChan,
		resetChan: resetChan,
		head:      0,
		maxConn:   maxConn,
		workerID:  workerID,
	}
}

func (h *History) Do() {
	go func() {
		for {
			select {
			case event := <-*h.eventChan:
				h.registerEvent(event)
			case ip := <-*h.resetChan:
				h.resetIp(ip)
			case <-h.quit:
				return
			}
		}
	}()
}

func (h *History) Stop() {
	close(h.quit)
}

func (h *History) registerEvent(e common.Event) {
	h.history[e.LineNum] = e

	for {
		nextEvent, nextEventExists := h.history[h.head+1]
		if nextEventExists {
			switch {
			case nextEvent.ConnectIP != "":
				h.registerConnect(nextEvent.ConnectIP)
			case nextEvent.PlayerEndIP != "":
				h.registerEndPlayer(nextEvent.PlayerEndIP)
			case nextEvent.PlayerStartIP != "":
				h.registerNewPlayer(nextEvent.PlayerEndIP)
			}
			delete(h.history, h.head+1)
			h.head++
		} else {
			break
		}
	}
}

func (h *History) registerConnect(ip string) {
	h.ips[ip]++
	if h.ips[ip] > h.maxConn {
		_, whitelisted := h.whitelist[ip]
		_, banned := h.banned[ip]
		if !whitelisted && !banned {
			h.banned[ip] = struct{}{}
			*h.banChan <- ip
		}
	}
}

func (h *History) registerNewPlayer(ip string) {
	h.whitelist[ip] = struct{}{}
}

func (h *History) registerEndPlayer(ip string) {
	delete(h.whitelist, ip)
	delete(h.ips, ip)
	delete(h.banned, ip)
}

func (h *History) resetIp(ip string) {
	delete(h.ips, ip)
	delete(h.banned, ip)
}
