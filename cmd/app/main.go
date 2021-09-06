package main

import (
	"context"
	"github.com/belliel/multiplexer/configs"
	"github.com/belliel/multiplexer/internal/transport"
	"github.com/belliel/multiplexer/pkg/os"
	"log"
)

func main() {
	log.Println("[INFO] Multiplexer inits...")

	mainCtx, cancelMainCtx := context.WithCancel(context.Background())
	defer cancelMainCtx()
	go os.CatchTermination(cancelMainCtx)

	config := configs.NewAppConfig()

	server := transport.NewTransportBuilder(mainCtx, transport.HTTP).
		WithAddr(config.ListenAddr).
		WithDebug(config.Debug).
		Build()

	if err := server.Listen(); err != nil && err.Error() != "http: Server closed" {
		log.Printf("[ERROR] Server cannot serve: %v", err)
		return
	}
	server.WaitForGracefulShutdown()
	log.Println("[WARN] Process terminated")
}
