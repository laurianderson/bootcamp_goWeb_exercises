package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	//Crea un router con gin 
	router := gin.Default()

	//Captura la solicitud get hello world
	router.GET("/hello-world", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello world!",
		})
    })

	//Corremos nuestro servidor sobre el puerto 8080
	router.Run()
}


