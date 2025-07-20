package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type Payment struct {
	Amount        decimal.Decimal `json:"amount"`
	CorrelationID string          `json:"correlationId"`
}

type PaymentTBProcessed struct {
	Amount        decimal.Decimal `json:"amount"`
	CorrelationID string          `json:"correlationId"`
	RequestedAt   string          `json:"requestedAt"`
}

func PaymentHandler(c *gin.Context) {
	fmt.Println("Processando pagamento...")

	var payment Payment
	if err := c.ShouldBindJSON(&payment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ptbp := PaymentTBProcessed{
		Amount:        payment.Amount,
		CorrelationID: payment.CorrelationID,
		RequestedAt:   time.Now().Format(time.RFC3339),
	}

	jsonData, err := json.Marshal(ptbp)
	if err != nil {
		fmt.Println("Erro ao converter para JSON:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process payment"})
		return
	}

	r, err := http.Post("http://payment-processor-default:8080/payments", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Erro ao enviar pagamento:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to send payment",
			"details": err.Error(),
		})
		return
	}
	defer r.Body.Close()

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}
	fmt.Println("Resposta do serviço de pagamento:", string(bodyBytes))

	c.JSON(http.StatusOK, gin.H{
		"message": "processado com sucesso",
	})
}
