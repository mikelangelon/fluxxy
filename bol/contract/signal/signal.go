package signal

import (
	"github.com/inconshreveable/log15"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var wg sync.WaitGroup

func SignalHandler() {

	channel := make(chan os.Signal)

	signal.Notify(channel, syscall.SIGINT, syscall.SIGTERM)

	for signal := range channel {
		switch signal {
		case syscall.SIGINT, syscall.SIGTERM:
			now := time.Now()
			log15.Info("Caught signal, exit....")
			close(channel)
			duration := time.Since(now)

			log15.Info("Shutdown in: ", "msec", duration.Nanoseconds()/1000000)
			os.Exit(0)
		default:
		}
	}

}

func complete(inner func()) {
	defer wg.Done()
	inner()
}
