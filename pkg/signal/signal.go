package signal

import (
	"context"
	"os"
	"os/signal"
)

var handleMap = make(map[string]handleSlice)

// Not a concurrent safe function.
func registerHandler(
	sig os.Signal,
	handler ...func(context.Context)) {

	if hslice, ok := handleMap[sig.String()]; ok {
		handleMap[sig.String()] = append(hslice, handler...)
		return
	}

	sigchan := make(chan os.Signal)
	signal.Notify(sigchan, sig)

	go func() {
		for {
			<-sigchan

			for h := range handleMap[sig.String()] {
				go h()
			}
		}
	}()
}
