/*
Vamos a crear un endpoint llamado /saludo. Con una pequeña estructura con nombre y apellido que al pegarle deberá responder en texto
“Hola + nombre + apellido”

El endpoint deberá ser de método POST
Se deberá usar el package JSON para resolver el ejercicio
La respuesta deberá seguir esta estructura: “Hola Andrea Rivas”
La estructura deberá ser como esta:
{
		“nombre”: “Andrea”,
		“apellido”: “Rivas”
}
*/

package main

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
)

type Saludo struct {
	// ` ` tags para operar con json y los campos me queden con minúscula
	Nombre   string `json:"nombre"`
	Apellido string `json:"apellido"`
}


func main() { 
// Utilizo ShouldBindJSON para pasar de estructura a un JSON

	var s Saludo

	router := gin.Default()

	router.POST("/saludar", func(c *gin.Context) {


		if err := c.ShouldBindJSON(&s); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		mensaje := fmt.Sprintf("Hola %v %v", s.Nombre, s.Apellido)

		c.JSON(http.StatusOK, gin.H{"mensaje": mensaje})
	})
	router.Run()
	
}