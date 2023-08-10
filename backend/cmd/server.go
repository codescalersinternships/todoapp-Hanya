package main

import (
	"log"
	todo "github.com/codescalersinternships/todoapp-Hanya/internal"
)

func main() {
	dbPath, port, err := todo.ParseCommandLine()
	if err != nil {
		log.Fatal(err)
	}
	app := &todo.App{}
	err = app.NewApp(port,dbPath)
	if err != nil {
		log.Fatal(err)
	}
	err = app.Run(port)
	if err != nil {
		log.Fatal(err)
	}
}
