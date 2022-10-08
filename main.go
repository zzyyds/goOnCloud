package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	go startListen()
	logrus.Info("start routing...")
	{
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
		<-sig
	}
}

func startListen() {
	r := NewRouter()
	s := &http.Server{
		Addr:           ":9876",
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	err := s.ListenAndServe()
	if err != nil {
		panic(err)
	}
	go func() {
		// close gracefully
		gracefullyShutdown(s)
	}()
}

func gracefullyShutdown(server *http.Server) {

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done

	logrus.Info("closing http server gracefully ...")

	err := server.Shutdown(context.Background())
	if err != nil {
		log.Panic("closing http server gracefully failed: ", err)
	}

}
