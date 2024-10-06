package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kamarajugadda-pavan-kumar/booking-service-GOLANG/internal/config"
)

func main() {

	// load configuration
	cfg := config.MustGetConfig()

	// database setup

	// setup router
	router := http.NewServeMux()
	router.HandleFunc("GET /api/v1/bookings", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"message": "Booking created successfully"}`))
	})

	// setup http server
	httpServer := &http.Server{
		Addr:    cfg.HTTPServer.Address + ":" + cfg.HTTPServer.Port,
		Handler: router,
	}

	// start server
	fmt.Printf("server starting on http://localhost:%s", cfg.HTTPServer.Port)
	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := httpServer.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	<-done

	slog.Info("Shutting down the server", slog.String("Address",
		cfg.HTTPServer.Address+":"+cfg.HTTPServer.Port))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := httpServer.Shutdown(ctx)
	if err != nil {
		slog.Error("Failed to shutdown server", slog.String("error", err.Error()))
	}
	slog.Info("Server shutdown sucessfully")
}
