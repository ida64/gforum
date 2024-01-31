package main

import (
	"crypto/rand"
	"log"
	"math/big"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

/*
* generateRandomString generates a random string of the given length
 */
func generateRandomString(length int) string {
	const alphanumeric = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var builder strings.Builder
	builder.Grow(length)

	for i := 0; i < length; i++ {
		randomIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(alphanumeric))))
		builder.WriteByte(alphanumeric[randomIndex.Int64()])
	}

	return builder.String()
}

/*
* getParamterInt returns the integer value of the parameter with the given name.
* If the parameter is not an integer, 0 is returned.
 */
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
