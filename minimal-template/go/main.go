package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Authentication middleware
func authMiddleware() gin.HandlerFunc {
	sspAPIKey := os.Getenv("SSP_API_KEY")

	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")

		if apiKey == "" || apiKey != sspAPIKey {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   true,
				"message": "Unauthorized - Invalid or missing API key",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// Health check endpoint
func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"version":   "1.0.0",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// Capabilities endpoint
func capabilitiesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"supported_methods":    []string{"your_methods_here"},
		"supported_currencies": []string{"USD", "INR"},
		"features":             []string{"feature1", "feature2"},
	})
}

// Example endpoint with authentication
func yourEndpointHandler(c *gin.Context) {
	var requestBody map[string]interface{}

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid JSON request body",
		})
		return
	}

	// TODO: Implement your business logic here
	// Example: Process the request data
	// data := requestBody["data"]

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Request processed successfully",
		"data":    requestBody,
	})
}

// Error handler middleware
func errorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			log.Printf("Error: %v", err.Err)

			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   true,
				"message": "Internal server error",
			})
		}
	}
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Validate required environment variables
	if os.Getenv("SSP_API_KEY") == "" {
		log.Fatal("SSP_API_KEY environment variable is required")
	}

	// Set Gin mode
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router
	router := gin.Default()

	// Apply error handler
	router.Use(errorHandler())

	// Public routes
	router.GET("/health", healthHandler)
	router.GET("/capabilities", capabilitiesHandler)

	// Protected routes
	protected := router.Group("/")
	protected.Use(authMiddleware())
	{
		protected.POST("/your-endpoint", yourEndpointHandler)
	}

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Start server
	log.Printf("SSP Plugin Server starting on port %s", port)
	log.Printf("Health: http://localhost:%s/health", port)
	log.Printf("Capabilities: http://localhost:%s/capabilities", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
