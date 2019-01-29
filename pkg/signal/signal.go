package signal

import (
	"os"
	"os/signal"
)

var (
	// key: signal name, value: channel
	registered = make(map[string]cancelChan)
)

func registerHandler(handler func(), sig os.Signal) {

	if cancel, ok := registered[sig.String()]; ok {
		close(cancel)
	}

	cancelchan := make(cancelChan)
	registered[sig.String()] = cancelchan

	sigchan := make(chan os.Signal)
	signal.Notify(sigchan, sig)

	go func() {
		for {
			select {
			case <-cancelchan:
				return
			case <-sigchan:
				handler()
			}
		}
	}()
}
