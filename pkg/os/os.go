package os

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func CatchTermination(cancel context.CancelFunc) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("[WARN] Caught termination signal")
	cancel()
}
