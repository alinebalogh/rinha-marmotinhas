package api

import "github.com/gin-gonic/gin"

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Define a rota para o endpoint de pagamento
	router.POST("/payments", PaymentHandler)

	// Rota de teste
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	return router
}
