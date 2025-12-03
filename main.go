package main

import (
	"fmt"
	"os"

	"github.com/cbwinslow/cbwsh/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
