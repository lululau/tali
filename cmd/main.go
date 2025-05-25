package main

import (
	"fmt"
	"os"

	"aliyun-tui-viewer/internal/app"
)

func main() {
	application, err := app.New()
	if err != nil {
		fmt.Printf("Error initializing application: %v\n", err)
		os.Exit(1)
	}

	if err := application.Run(); err != nil {
		fmt.Printf("Error running application: %v\n", err)
		os.Exit(1)
	}
}
