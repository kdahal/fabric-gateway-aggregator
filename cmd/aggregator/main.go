package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fabric-gateway-aggregator/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
)

func main() {
	// 1. Initialize Distributed Messaging (NATS)
	// In a real scenario, this would point to a GCP/AWS managed NATS cluster
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = nats.DefaultURL
	}

	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer nc.Close()

	// 2. Initialize the Aggregator Logic
	aggregator := service.NewAggregatorService(nc)

	// 3. Set up the Web Server (Implementing the OpenAPI Spec)
	router := gin.Default()

	v1 := router.Group("/v1")
	{
		v1.POST("/ports/provision", func(c *gin.Context) {
			var req struct {
				PortID     string `json:"portId" binding:"required"`
				VlanID     int    `json:"vlanId" binding:"required"`
				FabricZone string `json:"fabricZone" binding:"required"`
			}

			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body per OpenAPI spec"})
				return
			}

			// Dispatch to our asynchronous service
			ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
			defer cancel()

			if err := aggregator.ProvisionPort(ctx, req.PortID, req.VlanID, req.FabricZone); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to dispatch task"})
				return
			}

			c.JSON(http.StatusAccepted, gin.H{"status": "Provisioning initiated"})
		})
	}

	// 4. Graceful Shutdown
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		log.Printf("Fabric Aggregator listening on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down Fabric Aggregator...")
	srv.Shutdown(context.Background())
}