package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

// Sample product data
var products = []Product{
	{ID: 1, Name: "Laptop", Price: 999},
	{ID: 2, Name: "Smartphone", Price: 499},
	{ID: 3, Name: "Headphones", Price: 99},
	{ID: 4, Name: "Monitor", Price: 299},
	{ID: 5, Name: "Keyboard", Price: 59},
}

func getSignedURL(bucket string, key string, s3Client *s3.Client) (string, error) {
	presignedClient := s3.NewPresignClient(s3Client)

	presignDuration := 10 * time.Minute

	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	presignedURL, err := presignedClient.PresignGetObject(context.TODO(), input, s3.WithPresignExpires(presignDuration))
	if err != nil {
		return "", err
	}
	return presignedURL.URL, nil
}

func main() {
	err := godotenv.Load(".env.development")
	if err != nil {
		log.Println("app is running in development mode")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Printf("unable to load SDK config, %v", err)
	}

	s3Client := s3.NewFromConfig(cfg)
	log.Println(s3Client)

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
		// Get price range filter parameters
		minPriceStr := c.Query("min_price")
		maxPriceStr := c.Query("max_price")

		// Filter products based on price range
		filteredProducts := make([]Product, 0)
		for _, p := range products {
			include := true

			if minPriceStr != "" {
				minPrice, err := strconv.Atoi(minPriceStr)
				if err == nil && p.Price < minPrice {
					include = false
				}
			}

			if maxPriceStr != "" {
				maxPrice, err := strconv.Atoi(maxPriceStr)
				if err == nil && p.Price > maxPrice {
					include = false
				}
			}

			if include {
				filteredProducts = append(filteredProducts, p)
			}
		}

		c.HTML(http.StatusOK, "index.html", gin.H{
			"products":  filteredProducts,
			"min_price": minPriceStr,
			"max_price": maxPriceStr,
		})
	})

	port := "5000"
	err = router.Run(":" + port)
	if err != nil {
		log.Fatalf("error running app: %v", err)
	}
}
