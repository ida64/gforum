package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

/*
* GetParamterInt returns the integer value of the parameter with the given name.
* If the parameter is not an integer, 0 is returned.
 */
func GetParamterInt(c *gin.Context, name string) int {
	var param = c.Param(name)

	var paramInt, err = strconv.Atoi(param)
	if err != nil {
		return 0
	}

	return paramInt
}
