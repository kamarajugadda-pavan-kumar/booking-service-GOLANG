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

	"github.com/gorilla/mux"
	"github.com/kamarajugadda-pavan-kumar/booking-service-GOLANG/internal/config"
	"github.com/kamarajugadda-pavan-kumar/booking-service-GOLANG/internal/cronjob"
	"github.com/kamarajugadda-pavan-kumar/booking-service-GOLANG/internal/db"
	handlers_v1 "github.com/kamarajugadda-pavan-kumar/booking-service-GOLANG/internal/http/handlers/v1"
)

func main() {

	// load configuration
	cfg := config.MustGetConfig()

	// database setup
	db.InitDBCredentials(&cfg.Database)
	db := db.GetDB()
	defer db.Close()

	// setup router
	router := mux.NewRouter()
	handlers_v1.RegisterV1Routes(router)

	// setup http server
	httpServer := &http.Server{
		Addr:    "0.0.0.0:" + cfg.HTTPServer.Port,
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

	//================================================================
	//  Start the cron job to clean expired bookings
	go cronjob.CleanExpiredBookings()
	//================================================================

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
