package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/feliperyan/ai_preso_coach/go-backend/api"
)

func main() {
	fmt.Println("this is main")

	server := api.NewServer()
	router := http.NewServeMux()

	handler := api.HandlerFromMux(server, router)

	s := &http.Server{
		Handler: handler,
		Addr:    "0.0.0.0:8090",
	}

	log.Fatal(s.ListenAndServe())
}
