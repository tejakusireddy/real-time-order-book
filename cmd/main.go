package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tejakusireddy/real-time-order-book/internal/engine"
	websocketTransport "github.com/tejakusireddy/real-time-order-book/internal/transport/websocket"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	hub := websocketTransport.NewHub()
	go hub.Run()

	orderBook := engine.NewOrderBook(hub)

	// HTTP server (placeholder — WebSocket comes later)
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	mux.HandleFunc("/ws", hub.HandleWS(orderBook))

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		log.Println("🚀 Server starting on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server crashed: %v", err)
		}
	}()

	// Graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
	log.Println("🛑 Shutdown signal received...")

	ctxTimeout, cancelTimeout := context.WithTimeout(ctx, 5*time.Second)
	defer cancelTimeout()

	if err := srv.Shutdown(ctxTimeout); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}

	log.Println("✅ Server exited cleanly")
}
