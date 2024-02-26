package app

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"VEEEKTOR_api/internal/service"
	"VEEEKTOR_api/pkg/database/pgsql"
)

var apiPrefix = "/api"

func Start() {
	log.Printf("VEEEKTOR_api is starting...")

	defer pgsql.DB.Close()

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	mux := getMultiplexer()
	server := &http.Server{Addr: ":8080", Handler: mux}
	defer server.Close()

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit
}

func getMultiplexer() *http.ServeMux {
	mux := http.NewServeMux()

	// Auth free
	mux.HandleFunc(apiPrefix+"/users/signin", service.UsersSignInHandler)
	mux.HandleFunc(apiPrefix+"/users/signup", service.UsersSignUpHandler)

	mux.HandleFunc(apiPrefix+"/auth/refresh", service.UpdateToken)

	return mux
}
