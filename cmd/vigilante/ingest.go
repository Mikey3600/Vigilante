package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var inputFile string
var ingestCmd = &cobra.Command{Use:"ingest", RunE: func(cmd *cobra.Command, args []string) error {
	if inputFile=="" { return fmt.Errorf("--file required") }
	data,err:=os.ReadFile(inputFile); if err!=nil { return err }
	url := "http://localhost:"+getenv("HTTP_PORT","8080")+"/api/v1/logs"
	resp,err:=http.Post(url,"application/json",bytes.NewBuffer(data)); if err!=nil { return err }
	defer resp.Body.Close(); b,_:=io.ReadAll(resp.Body); fmt.Printf("%d %s\n", resp.StatusCode, string(b)); return nil
}}
func init(){ ingestCmd.Flags().StringVar(&inputFile,"file","","json file"); rootCmd.AddCommand(ingestCmd)}
