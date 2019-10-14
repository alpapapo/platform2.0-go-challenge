package main

import (
	"fmt"
	"github.com/GlobalWebIndex/platform2.0-go-challenge/app"
	"os"
)


func main() {
	projectEnv := os.Getenv("PROJECT_ENV")
	fmt.Println(projectEnv)
	app := app.App{}
	app.Initialize(projectEnv)
	defer app.DB.Close()
	app.Run()
}
