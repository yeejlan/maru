package maru

import (
	"fmt"
	"log"
	"time"
	"net/http"
	"os"
	"os/signal"
	"context"
)

//start http server
func StartHttpServer(host string, port int) {
	var wait time.Duration = time.Second * 30

	router := NewRouter()
	addr := fmt.Sprintf("%s:%d", host, port)

	srv := &http.Server{
		Addr: addr,
		WriteTimeout: time.Second * 15,
		ReadTimeout: time.Second * 15,
		IdleTimeout: time.Second * 60,
		Handler: router,
	}

	log.Printf("listen on %s", addr);
	// Run server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// will accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until receive a signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("http server shutting down")	
}