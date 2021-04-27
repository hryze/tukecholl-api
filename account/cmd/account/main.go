package main

import (
	"log"

	"github.com/paypay3/tukecholl-api/account/infrastructure/server"
)

func main() {
	if err := server.Run(); err != nil {
		log.Fatalf("%+v", err)
	}
}
