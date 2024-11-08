package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const listenAddres = "80"

type Config struct {
}

func main() {

	c := &Config{}

	log.Printf("server listeninig on port %s\n", listenAddres)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", listenAddres),
		Handler: c.routes(),
	}

	err := srv.ListenAndServe()

	if err != nil {
		log.Panic(err)
	}
}
