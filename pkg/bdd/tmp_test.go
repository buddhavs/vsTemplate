package signal

import (
	// "github.com/smartystreets/assertions/should"
	"os"
	"syscall"
	"testing"

	"github.com/smartystreets/gunit"
)

type (
	signalType struct {
		*gunit.Fixture
		sigs []syscall.Signal
		pid  int
	}
)

func TestSignalHandler(t *testing.T) {
	gunit.Run(new(signalType), t)
}

func (st *signalType) SetupTest() {
	st.sigs = []syscall.Signal{
		syscall.SIGINT,
		syscall.SIGHUP,
	}

	st.pid = os.Getpid()
}

func (st *signalType) TeardownTest() {
}

func raise(sig os.Signal) error {
	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		return err
	}
	return p.Signal(sig)
}
