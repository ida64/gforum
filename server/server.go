package main

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getParamterInt(c *gin.Context, name string) int {
	var param = c.Param(name)
	var paramInt, err = strconv.Atoi(param)
	if err != nil {
		return 0
	}
	return paramInt
}

func main() {
	err := listenServer()
	if err != nil {
		log.Fatal(err)
	}
}
