package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	gokitrest01 "github.com/aachi/gokit-rest01"
)

//comments

func main() {
	var (
		httpAddr = flag.String("http", ":8080", "http listen address")
	)
	flag.Parse()
	ctx := context.Background()
	// our napodate service
	srv := gokitrest01.NewService()
	errChan := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	// mapping endpoints
	endpoints := gokitrest01.Endpoints{
		GetEndpoint:      gokitrest01.MakeGetEndpoint(srv),
		StatusEndpoint:   gokitrest01.MakeStatusEndpoint(srv),
		ValidateEndpoint: gokitrest01.MakeValidateEndpoint(srv),
	}

	// HTTP transport
	go func() {
		log.Println("gokitrest01 is listening on port:", *httpAddr)
		handler := gokitrest01.NewHTTPServer(ctx, endpoints)
		errChan <- http.ListenAndServe(*httpAddr, handler)
	}()

	log.Fatalln(<-errChan)
}
