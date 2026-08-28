package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"example.com/graduation-showcase/internal/api"
	"example.com/graduation-showcase/internal/flow014"
	"example.com/graduation-showcase/internal/store"
)

func main() {
	serve := flag.Bool("serve", false, "start local HTTP service")
	path := flag.String("db", filepath.Join(os.TempDir(), "graduation-showcase.db"), "bolt database path")
	flag.Parse()
	persistence, err := store.Open(*path)
	if err != nil {
		log.Fatal(err)
	}
	defer persistence.Close()
	service := &api.Service{Store: persistence, Clock: api.FixedClock{Value: 202501010900}}
	if *serve {
		server := &http.Server{Addr: "127.0.0.1:8080", Handler: api.NewHandler(service)}
		log.Printf("graduation showcase listening on %s", server.Addr)
		log.Fatal(server.ListenAndServe())
	}
	if err := runDemo(service, persistence); err != nil {
		log.Fatal(err)
	}
	f := flow014.New(service, persistence)
	label, err := f.CurrentLabel("demo-1")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("demo label: %s\n", label)
}
