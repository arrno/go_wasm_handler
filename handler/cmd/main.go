package main

import (
	"log"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	_ "github.com/arrno/go_wasm_handler"
)

func main() {
  port := "8080"
  if envPort := os.Getenv("PORT"); envPort != "" {
    port = envPort
  }
  if err := funcframework.Start(port); err != nil {
    log.Fatalf("funcframework.Start: %v\n", err)
  }
}