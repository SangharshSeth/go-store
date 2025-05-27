package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

// sampleProducts is a static list of products to replace database queries
var sampleProducts = []Product{
	{ID: 1, Name: "iPhone 13 Pro", Price: 999},
	{ID: 2, Name: "MacBook Air M1", Price: 1299},
	{ID: 3, Name: "Samsung Galaxy S21", Price: 799},
	{ID: 4, Name: "Dell XPS 13", Price: 1199},
	{ID: 5, Name: "Sony WH-1000XM4", Price: 349},
	{ID: 6, Name: "iPad Pro", Price: 799},
	{ID: 7, Name: "Nintendo Switch", Price: 299},
	{ID: 8, Name: "Bose QuietComfort Earbuds", Price: 279},
}

func main() {
	// Load environment variables if available but not required anymore
	_ = godotenv.Load(".env.development")

	router := gin.Default()

	// Load HTML templates
	router.LoadHTMLGlob("templates/*")

	// Add custom template functions
	router.SetFuncMap(template.FuncMap{
		"len": func(v interface{}) int {
			switch val := v.(type) {
			case []Product:
				return len(val)
			default:
				return 0
			}
		},
	})

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"products": sampleProducts,
		})
	})

	port := "5000"
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("error running app: %v", err)
	}
}
