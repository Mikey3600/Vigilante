package main

import (
	"context"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/vigilante/internal/ai"
	"github.com/user/vigilante/internal/api"
	"github.com/user/vigilante/internal/grpc"
	"github.com/user/vigilante/internal/storage"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts the Vigilante HTTP and gRPC servers",
	Run: func(cmd *cobra.Command, args []string) {
		dbUrl := os.Getenv("DATABASE_URL")
		db, err := storage.NewDB(context.Background(), dbUrl)
		if err != nil {
			log.Fatalf("Failed to init db: %v", err)
		}

		aiClient, err := ai.NewClient(context.Background())
		if err != nil {
			log.Printf("Warning: Failed to init AI client: %v", err)
		}

		go func() {
			grpcPort := os.Getenv("GRPC_PORT")
			if grpcPort == "" {
				grpcPort = "50051"
			}
			log.Printf("Starting gRPC on :%s", grpcPort)
			if err := grpc.Start(grpcPort, db); err != nil {
				log.Fatalf("gRPC server failed: %v", err)
			}
		}()

		port := os.Getenv("PORT")
		if port == "" {
			port = "3000"
		}
		
		r := api.SetupRouter(db, aiClient)
		log.Printf("Starting HTTP on :%s", port)
		if err := r.Run(":" + port); err != nil {
			log.Fatalf("HTTP server failed: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
