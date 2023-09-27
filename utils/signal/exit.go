package signal

import (
	"os"
	"os/signal"
	"syscall"
)

func OnExit() chan os.Signal {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP)
	return sig
}
