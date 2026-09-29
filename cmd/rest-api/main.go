package main

import (
	"fmt"

	"github.com/MHNahib/rest-api/internal/config"
)

func main() {
	fmt.Println("Bismillah")

	appConfig := config.MountConfig()

	fmt.Println(appConfig)
}
