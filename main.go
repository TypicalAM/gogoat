package main

import (
	"fmt"
	"log"

	"github.com/TypicalAM/gogoat/client"
	"github.com/TypicalAM/gogoat/config"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	client := client.NewCaller(*cfg)
	th, err := client.GetTotalHits()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(th.PrettyPrint())
}
