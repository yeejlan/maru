package maru

import (
	"fmt"
	"log"
	"time"
	"net/http"
	"os"
	"os/signal"
	"context"
	"strings"
	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
)

//start http server
func StartHttpServer(router *Router, host string, port int) {
	var wait time.Duration = time.Second * 15

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
			Shutdown()
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

	go func() {
		// Doesn't block if no connections, but will otherwise wait
		// until the timeout deadline.		
		srv.Shutdown(ctx)
	}()

	log.Println("http server shutting down")

	<-ctx.Done()
}

//shutdown function, close redis and db connections
func Shutdown() {
	for k, v := range Registry.GetMap() {
		if(strings.HasPrefix(k, "redis.")){
			redis := v.(*redis.Client)
			redis.Close()
		}
		if(strings.HasPrefix(k, "db.")){
			db := v.(*sqlx.DB)
			db.DB.Close()
		}
	}
}
