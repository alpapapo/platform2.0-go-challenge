package common

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

func LoadDotenv(env string) error {

	if "" == env {
		env = "development"
	}

	envLocalError := godotenv.Overload(".env." + env + ".local")
	fmt.Println(os.Getenv(""))
	if envLocalError != nil {
		fmt.Print(envLocalError)
	}

	if "test" != env {
		defaultLocalError := godotenv.Load(".env.local")
		if defaultLocalError != nil {
			fmt.Println(defaultLocalError)
		}

	}
	envError := godotenv.Load(".env." + env)
	if envError != nil {
		fmt.Println(envError)
	}
	e := godotenv.Load() // The Original .env

	return e
}