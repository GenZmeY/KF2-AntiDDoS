package reader

import (
	"kf2-antiddos/internal/common"
	"kf2-antiddos/internal/output"

	"bufio"
	"os"
)

type Reader struct {
	quit       chan struct{}
	outputChan *chan common.RawEvent
	workerID   uint
}

func New(workerID uint, outputChan *chan common.RawEvent) *Reader {
	return &Reader{
		outputChan: outputChan,
		quit:       make(chan struct{}),
		workerID:   workerID,
	}
}

func (r *Reader) Do() {
	go func() {
		stdin := bufio.NewScanner(os.Stdin)
		stdin.Split(bufio.ScanLines)
		for {
			select {
			case <-r.quit: // check quit if there are no input
				return
			default:
				for lineNum := byte(1); stdin.Scan(); lineNum++ { // byte overflow it's not a bug, but a feature
					select {
					case <-r.quit: // check quit if there are input
						return
					default:
					}

					text := stdin.Text()
					output.Proxyln(text)
					*r.outputChan <- common.RawEvent{
						LineNum: lineNum,
						Text:    text,
					}
				}

				if err := stdin.Err(); err != nil {
					output.Errorln(err.Error())
				}
			}
		}
	}()
}

func (r *Reader) Stop() {
	close(r.quit)
}
