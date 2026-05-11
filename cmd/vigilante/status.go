package main

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{Use:"status", RunE: func(cmd *cobra.Command, args []string) error {
	url := "http://localhost:"+getenv("HTTP_PORT","8080")+"/health"
	resp,err:=http.Get(url); if err!=nil { return err }
	defer resp.Body.Close(); fmt.Printf("HTTP\t%d\n", resp.StatusCode); return nil
}}

var versionCmd = &cobra.Command{Use:"version", Run: func(cmd *cobra.Command, args []string){ fmt.Println("vigilante 1.0.0") }}

func init(){ rootCmd.AddCommand(statusCmd, versionCmd) }
