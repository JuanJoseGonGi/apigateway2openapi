package main

import (
	"log"
	"os"

	"github.com/mauriciocm9/apigateway2openapi"
)

func main() {
	if len(os.Args) != 2 {
		log.Println("Missing openapi path")
		return
	}

	b, err := apigateway2openapi.Execute(os.Args[1])
	if err != nil {
		log.Println("Found error trying to read file:", err.Error())
		return
	}

	err = os.WriteFile("out.yaml", b, 0600)
	if err != nil {
		log.Println("Found error writing to file:", err.Error())
		return
	}
}
