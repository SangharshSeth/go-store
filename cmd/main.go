package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

func main() {
	err := godotenv.Load(".env.development")
	if err != nil {
		log.Println("app is running in development mode")
	}

	router := gin.Default()

	// Load HTML templates
	router.LoadHTMLGlob("templates/*")

	router.GET("/", func(c *gin.Context) {
		products := []Product{
			{ID: 1, Name: "IPhone 16 Pro", Price: 999},
			{ID: 2, Name: "Macbook Pro 16", Price: 2999},
			{ID: 3, Name: "IPad Pro 13", Price: 1999},
		}

		c.HTML(http.StatusOK, "index.html", gin.H{
			"products": products,
		})
	})

	router.GET("/products", func(c *gin.Context) {
		var products []Product
		products = append(products, Product{Name: "IPhone 16 Pro", Price: 999})
		products = append(products, Product{Name: "Macbook Pro 16", Price: 2999})
		products = append(products, Product{Name: "IPad Pro 13", Price: 1999})
		c.JSON(200, products)
	})

	err = router.Run(":8080")
	if err != nil {
		log.Fatalf("error running app: %v", err)
	}
}
