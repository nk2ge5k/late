package shutdown

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const GracePeriod = 1 * time.Minute

var shutdown struct {
	mu       sync.Mutex
	once     sync.Once
	sequence []func()
}

func init() {
	// Close the listener when a shutdown event happens.
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, os.Interrupt)
	go func() {
		sig := <-c
		logf("shutdown: process received signal %v", sig)
		Now(0)
	}()
}

// Handle adds new termination handler that must be called before exit.
func Handle(onShutdown func()) {
	shutdown.mu.Lock()
	defer shutdown.mu.Unlock()

	shutdown.sequence = append(shutdown.sequence, onShutdown)
}

// Now stops entire program, in case if termination was not reached after
// GracePeriod, forces shutdown.
func Now(code int) {
	shutdown.once.Do(func() {
		// Ensure we terminate after a fixed amount of time.
		go func() {
			time.Sleep(GracePeriod)
			logf("shutdown: %v elapsed since shutdown requested - exiting forcefully",
				GracePeriod)
			os.Exit(1)
		}()

		shutdown.mu.Lock() // No need to ever unlock.
		for i := len(shutdown.sequence) - 1; i >= 0; i-- {
			shutdown.sequence[i]()
		}

		os.Exit(code)
	})
}

func logf(format string, args ...any) {
	_, _ = fmt.Fprintf(os.Stdout, format, args...)
}
