package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	var apiRuntime = gin.Default()
	var apiRuntimeError = apiRuntime.Run()
	if nil != apiRuntimeError {
		log.Fatalf("[API] %s", apiRuntimeError.Error())
	}
}
