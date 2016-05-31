package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {
	port := "8080"
	router := httprouter.New()
	router.ServeFiles("/public/*filepath", http.Dir("public"))
	router.GET("/", index)
	router.GET("/google", oauth)
	router.GET("/login", login)
	router.GET("/auth", auth)
	router.GET("/auth/callback", oauthCallback)
	log.Printf("Running on port: %v", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
