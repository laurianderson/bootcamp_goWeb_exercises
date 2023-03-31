/*
Vamos a levantar un servidor utilizando el paquete gin en el puerto 8080. Para probar nuestros endpoints haremos uso de postman.
Crear una ruta /ping que debe respondernos con un string que contenga pong con el status 200 OK.
Crear una ruta /products que nos devuelva la lista de todos los productos en la slice.
Crear una ruta /products/:id que nos devuelva un producto por su id.
Crear una ruta /products/search que nos permita buscar por parámetro los productos cuyo precio sean mayor a un valor priceGt.
*/

package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"github.com/gin-gonic/gin"
	"strconv"
)

//Creamos una estructura producto siguiendo los campos que tiene el products.json
type Product struct {
	ID         	 int64 `json:"id"`
    Name       	 string `json:"name"`
	Quantity   	 int `json:"quantity"`
	Code_Value 	 string `json:"code_value"`
	Is_Published bool `json:"is_published"`
	Expiration   string `json:"expiration"`
	Price        float64 `json:"price"`
}


//Creamos una función que lea el archivo products.json y me devuelva un slide de productos donde los va a guardar
func readJson() []Product {
	var Products []Product

	//lee el json y se almacena la info en bytes
	bytes, err := ioutil.ReadFile("./products.json") 
    if err!= nil {
        panic(err)
    }

	//Con la función unmarshal pasamos de byte a estructura y le indicamos que el valor(&) lo guarde en el slide Products
	json.Unmarshal([]byte(bytes), &Products)

	return Products
}


//Declaramos de manera global el slide de Productos
var ProductSlide = readJson()


//....................................................................................................................
//Funciones que actuan como el controlador
func getAll(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, ProductSlide)
}


func getById(ctx *gin.Context) {
	id := ctx.Param("id")


	for _, product := range ProductSlide {
		//El id que recibimos por párametro es del tipo string. Y el dato del tipo int. Vamos a parsear el dato recibido por párametro
		//para comparar dos datos del tipo
		value, _ := strconv.ParseInt(id, 10, 64)
		if product.ID == value {
            ctx.JSON(http.StatusOK, product)
            return
        } 
	}
	ctx.JSON(http.StatusNotFound, gin.H{"message":  "The product with ID: " + id + " does not exist"})
}

func searchByPriceMin(ctx *gin.Context) {
	//Los parámetros del endpoint deben tener el valor pasado por queryParams
	priceMin := ctx.Query("priceGt")

	//Creamos un slide para ir guardando los productos que va enontrando y así después imprimir el slide
	var productsFilter []Product

	for _, product := range ProductSlide {
		//el priceMin, lo recibimos como un string, lo parseamos a float para compararlo por el tipo de dato
		value, _ := strconv.ParseFloat(priceMin, 64)
		if product.Price > value {
			productsFilter = append(productsFilter, product)
			ctx.JSON(http.StatusOK, productsFilter)
        }
	}
	//Acá evaluamos si el largo del slide es igual a 0, error
	if len(productsFilter) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"message":  "Not found product with price: " + priceMin})
	}
}
func main() {
	
	router := gin.Default()

	//Endpoint a cual pegarles, con las funciones que explican cómo hacer eso.
	router.GET("products", getAll)
	router.GET("products/:id", getById)
	router.GET("products/search", searchByPriceMin)

	router.Run()
}