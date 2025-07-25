package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"rinhago/internal/api"
	"syscall"
	"time"
)

type HealthData struct {
	Failing         bool  `json:"failing"`
	MinResponseTime int64 `json:"minResponseTime"`
}

var ProcessorURI string = "payment-processor-default"

func main() {

	go func() {
		for {
			fmt.Println("Consultando service health default...")
			r, err := http.Get("http://payment-processor-default:8080/payments/service-health")
			if err != nil {
				fmt.Println("Erro ao consultar service health default:", err)
				return
			}
			defer r.Body.Close()

			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				log.Fatalf("Error reading response body: %v", err)
			}
			var data HealthData

			err = json.Unmarshal(bodyBytes, &data)
			if err != nil {
				log.Fatalf("Error unmarshalling response body: %v", err)
			}

			fmt.Println(data.Failing)
			fmt.Println(data.MinResponseTime)

			if data.Failing {
				// tbd: fazer get no health do service fallback e validar se ele está ok
				fmt.Println("Service default está com problemas, tentando fallback...")
				ProcessorURI = "payment-processor-fallback"
			} else {
				ProcessorURI = "payment-processor-default"
			}

			time.Sleep(time.Duration(5) * time.Second)
		}

	}()

	router := api.SetupRouter()

	srv := &http.Server{
		Addr:    ":8000",
		Handler: router,
	}

	// Canal para capturar sinais de interrupção
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	fmt.Println("Servidor iniciado em :8000")

	<-stop
	fmt.Println("Encerrando servidor...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Erro ao encerrar servidor: %v", err)
	}

	fmt.Println("Servidor finalizado com sucesso")

}
