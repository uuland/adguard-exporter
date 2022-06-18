package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ebrianne/adguard-exporter/config"
	"github.com/ebrianne/adguard-exporter/internal/adguard"
	"github.com/ebrianne/adguard-exporter/internal/metrics"
	"github.com/ebrianne/adguard-exporter/internal/server"
	log "github.com/sirupsen/logrus"
)

const (
	name = "adguard-exporter"
)

var (
	s *server.Server
)

func main() {
	conf := config.Load()

	if level, err := log.ParseLevel(conf.LogLevel); err != nil {
		log.Fatal(err)
	} else {
		log.SetLevel(level)
	}

	metrics.Init()

	initAdguardClient(conf.AdguardProtocol, conf.AdguardHostname, conf.AdguardUsername, conf.AdguardPassword, conf.AdguardPort, conf.Interval, conf.LogLimit, conf.RDnsEnabled)
	initHttpServer(conf.ServerPort)

	handleExitSignal()
}

func initAdguardClient(protocol string, hostnames []string, username, password, port string, interval time.Duration, logLimit string, rdnsenabled bool) {
	for _, hostname := range hostnames {
		client := adguard.NewClient(protocol, hostname, username, password, port, interval, logLimit, rdnsenabled)
		go client.Scrape()
	}
}

func initHttpServer(port string) {
	s = server.NewServer(port)
	go s.ListenAndServe()
}

func handleExitSignal() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	s.Stop()
	fmt.Println(fmt.Sprintf("\n%s HTTP server stopped", name))
}
