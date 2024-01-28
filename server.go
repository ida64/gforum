package main

import (
	"log"
)

func main() {
	err := listenServer()
	if err != nil {
		log.Fatal(err)
	}
}
