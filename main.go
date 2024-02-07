package main

import (
	controller "gforum/controller"
)

func main() {
	err := controller.ListenServer()
	if err != nil {
		panic(err)
	}

}
