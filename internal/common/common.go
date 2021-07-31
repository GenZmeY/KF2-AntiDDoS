package common

type Worker interface {
	Do()
	Stop()
}

type RawEvent struct {
	LineNum byte
	Text    string
}

type Event struct {
	LineNum       byte
	ConnectIP     string
	PlayerStartIP string
	PlayerEndIP   string
}
