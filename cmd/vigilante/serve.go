package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/user/vigilante/internal/api"
	igrpc "github.com/user/vigilante/internal/grpc"
	"github.com/user/vigilante/internal/storage"
	"golang.org/x/sync/errgroup"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start HTTP and gRPC servers",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		db, err := storage.NewDB(ctx, os.Getenv("DATABASE_URL"))
		if err != nil {
			return err
		}
		defer db.Close()

		r := api.SetupRouter(db)
		httpSrv := &http.Server{
			Addr:    ":" + getenv("HTTP_PORT", "8080"),
			Handler: r,
		}

		g, gctx := errgroup.WithContext(ctx)

		g.Go(func() error {
			return igrpc.Start(gctx, getenv("GRPC_PORT", "50051"), db)
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

		slog.Info("servers_started", "http", httpSrv.Addr)
		return g.Wait()
	},
}

func getenv(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}

func init() {
	rootCmd.AddCommand(serveCmd)
	_ = fmt.Sprintf
}
