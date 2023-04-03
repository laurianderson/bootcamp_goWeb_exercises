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
	//CAMBIAR ESTOOO!!!!!!
	"io/ioutil"
	"net/http"
	"github.com/gin-gonic/gin"
	"strconv"
	//"time"
	"log"
)


//Package service
type Product struct {
	ID         	 int64 `json:"id"`
    Name       	 string `json:"name"`
	Quantity   	 int `json:"quantity"`
	Code_Value 	 string `json:"code_value"`
	Is_Published bool `json:"is_published"`
	Expiration   string `json:"expiration"`
	Price        float64 `json:"price"`
}


//Declaramos de manera global el slide de Productos
var ProductSlide = readJson()

//Creamos una variable que guarde el valor del id y lo vaya incrementando
var lastID int

/*Generar un var de errores personalizados
var ( 
	ErrorProductIvalid = errors.New("Invalid Product")
	ErrorProductInternal = errors.New("Internal Error")
)
*/

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
	return nil

	

	return Products
	
	
}


//Package handler
func getPing() gin.HandlerFunc{
	return func(ctx *gin.Context){
		ctx.JSON(http.StatusOK, gin.H{"data": "pong"})
	}
}


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
		return
	}
}


var Token = "123456"

//Generamos una función que nos devuelva otra función con el context
func saveProduct() gin.HandlerFunc {
	type request struct {
		Name         string `json:"name" biding: "required"`
		Quantity     int `json:"quantity" biding: "required"`
		Code_Value 	 string `json:"code_value" biding: required`
		Is_Published bool `json:"is_published" biding:required`
		Expiration   string `json:"expiration" biding:required`
		Price        float64 `json:"price" biding:required`
	}
	
	return func(ctx *gin.Context) {
		var req request

		//auth
		token := ctx.GetHeader("token")
		if token != Token {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message":  "Error: invalid token"})
			return
		}

		//Hace el decoding + binding required (pasa la info del body al request)
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Error: Invalid request"})
			log.Print("Log:", err)
			return
		}
		//si la petición es correcta le agregamos un id a nuestro producto
		pr := Product{
			ID : int64(lastID) + 1,
			Name : req.Name,
            Quantity : req.Quantity,
            Code_Value : req.Code_Value,
            Is_Published : true,
            Expiration : req.Expiration,
			Price: req.Price,
		}
		for _,product := range ProductSlide{
			if pr.Code_Value == product.Code_Value{
				ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid request: code value must be unique", "data":nil})
				//log.Println(err)
				return
			}
		}
		
			ProductSlide = append(ProductSlide, pr)

			ctx.JSON(http.StatusOK, gin.H{"message": "Succes", "data": pr})
	}

}

/*
func (p *Product) Valid() error{
//Esto para validar en vez del biding en la estructura 

}
*/


func main() {
	
	router := gin.Default()
	router.GET("/ping", getPing())
	//Agrupo las rutas (añido)
	pr := router.Group("/products")
	{ 
		//Endpoint a cual pegarles, con las funciones que explican cómo hacer eso.
		pr.GET("/", getAll)
		pr.GET("/:id", getById)
		pr.GET("/search", searchByPriceMin)

		pr.POST("/", saveProduct())
	}

	router.Run()
}