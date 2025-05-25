package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
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

	// Connect to PostgresSQL
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := "postgres" // Default database name
	dbPort := "5432"     // Default PostgresSQL port

	// Create a connection string
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	// Connect to the database
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer conn.Close(context.Background())

	log.Println("Connected to PostgresSQL database")

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
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

		// Build the SQL query with optional price range filter
		query := "SELECT id, name, price FROM products"
		var args []interface{}
		var conditions []string

		// Add price range conditions if provided
		if minPriceStr != "" {
			// Convert string to int
			var minPrice int
			_, err := fmt.Sscanf(minPriceStr, "%d", &minPrice)
			if err == nil {
				conditions = append(conditions, "price >= $"+fmt.Sprintf("%d", len(args)+1))
				args = append(args, minPrice)
			} else {
				log.Printf("Error parsing min_price: %v", err)
			}
		}

		if maxPriceStr != "" {
			// Convert string to int
			var maxPrice int
			_, err := fmt.Sscanf(maxPriceStr, "%d", &maxPrice)
			if err == nil {
				conditions = append(conditions, "price <= $"+fmt.Sprintf("%d", len(args)+1))
				args = append(args, maxPrice)
			} else {
				log.Printf("Error parsing max_price: %v", err)
			}
		}

		// Add WHERE clause if we have conditions
		if len(conditions) > 0 {
			query += " WHERE " + conditions[0]
			for i := 1; i < len(conditions); i++ {
				query += " AND " + conditions[i]
			}
		}

		// Execute the query
		rows, err := conn.Query(context.Background(), query, args...)
		if err != nil {
			log.Printf("Error querying products: %v", err)
			c.HTML(http.StatusInternalServerError, "index.html", gin.H{
				"error": "Failed to fetch products",
			})
			return
		}
		defer rows.Close()

		var products []Product
		for rows.Next() {
			var p Product
			err := rows.Scan(&p.ID, &p.Name, &p.Price)
			if err != nil {
				log.Printf("Error scanning product row: %v", err)
				continue
			}
			products = append(products, p)
		}

		if err := rows.Err(); err != nil {
			log.Printf("Error iterating product rows: %v", err)
		}

		c.HTML(http.StatusOK, "index.html", gin.H{
			"products":  products,
			"min_price": minPriceStr,
			"max_price": maxPriceStr,
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}

	err = router.Run(":" + port)
	if err != nil {
		log.Fatalf("error running app: %v", err)
	}
}
