package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"geolocation-resolver/internal/geoloc_resolver/config"
	"geolocation-resolver/internal/geoloc_resolver/http"
	"geolocation-resolver/internal/geoloc_resolver/resolver"
	"geolocation-resolver/pkg/storage/hashmap"
	"geolocation-resolver/pkg/storage/scylla"
)

var (
	configPath = flag.String("config", "./config.yml", "path to config file")
)

func main() {
	flag.Parse()

	cfg, err := config.NewConfig(*configPath)
	if err != nil {
		log.Fatalf("[F] main: failed to create config: %s", err)
	}

	*cfg = cfg.WithDefaults()
	cfg.Print()

	db, shutdownDB := ChoseAndInitStorage(cfg)
	defer shutdownDB()

	r := resolver.NewResolver(&cfg.Resolver, db)
	if err = r.ImportCSV(cfg.PathToCSV); err != nil {
		log.Printf("[E] main: failed to parse CSV: %s\n", err)
	}

	httpServerShutdown := InitHTTPServer(cfg, r)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("[I] main: got SIGINT")

	httpServerShutdown()
}

func ChoseAndInitStorage(cfg *config.Config) (db resolver.GeoDB, shutdown func()) {
	switch cfg.StorageOption {
	case config.StorageOptionHashmap:
		return hashmap.NewStorage(), func() {}
	case config.StorageOptionScylla:
		log.Println("[I] main: starting scylla connector")

		scylla, err := scylla.NewConnector(&cfg.Scylla)
		if err != nil {
			log.Fatalf("[F] main: failed to create scylla connector: %s", err)
		}

		scylla.Init()
		log.Println("[I] main: scylla connector was initialised successfully")

		return scylla, scylla.Shutdown
	default:
		log.Fatalln("[F] main: unknown storage option")
	}

	return
}

func InitHTTPServer(cfg *config.Config, r *resolver.Resolver) func() {
	log.Println("[I] main: starting HTTP server")

	httpServer := http.NewServer(http.NewHandler(r), cfg.HTTP.Addr)
	httpServer.Init()

	return httpServer.Shutdown
}
