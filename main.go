package main

import (
	"context"
	"fmt"
	"os"

	"github.com/profx5/jordi/internal/app"
	"github.com/profx5/jordi/internal/config"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: jordi <url>")
		os.Exit(1)
	}
	url := os.Args[1]
	config := config.New(url)
	err := config.Validate()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	ctx := context.Background()
	app := app.New(config)
	if err := app.Run(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
