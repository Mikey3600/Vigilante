package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/user/vigilante/internal/api"
	igrpc "github.com/user/vigilante/internal/grpc"
	"github.com/user/vigilante/internal/storage"
	"golang.org/x/sync/errgroup"
)

func RunServe(ctx context.Context) error {
	fmt.Println("Vigilante starting...")

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}

	db, err := storage.NewDB(ctx, databaseURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()
	fmt.Println("Connected to database")

	if err := db.RunMigrations(ctx); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	fmt.Println("Migrations complete")

	r := api.SetupRouter(db)
	httpAddr := ":" + getenv("HTTP_PORT", "8080")
	grpcPort := getenv("GRPC_PORT", "50051")

	httpSrv := &http.Server{Addr: httpAddr, Handler: r}

	serveCtx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	g, gctx := errgroup.WithContext(serveCtx)

	g.Go(func() error {
		return igrpc.Start(gctx, grpcPort, db)
	})

	g.Go(func() error {
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	})

	g.Go(func() error {
		<-gctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		return httpSrv.Shutdown(shutdownCtx)
	})

	fmt.Printf("HTTP server listening on %s\n", httpAddr)
	fmt.Printf("gRPC server listening on :%s\n", grpcPort)
	fmt.Println("Vigilante is ready")

	return g.Wait()
}

func getenv(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
