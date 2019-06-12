package main

import (
	"io"
	"log"
	"net/http"

	"github.com/flume-cloud-services/database/middleware"
	"github.com/flume-cloud-services/file-storage/controllers"
)

func main() {
	http.HandleFunc("/signin", controllers.Signin)

	http.Handle("/welcome", middleware.Middleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			io.WriteString(w, "Welcome to you visitor !")
		}),
		middleware.AuthMiddleware,
	))

	log.Println("Starting server on port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
