package main

import (
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/flume-cloud-services/file-storage/controllers"
	"github.com/flume-cloud-services/file-storage/middleware"
)

func main() {
	http.HandleFunc("/signin", controllers.Signin)

	http.Handle("/welcome", middleware.Middleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			io.WriteString(w, "Welcome to you visitor !")
		}),
		middleware.AuthMiddleware,
	))

	http.Handle("/post", middleware.Middleware(
		http.HandlerFunc(controllers.PostFile),
		middleware.AuthMiddleware,
	))

	http.Handle("/private/", middleware.Middleware(
		http.StripPrefix(strings.TrimRight("/private/", "/"),
			http.FileServer(http.Dir("private"))),
		middleware.AuthMiddleware,
	))
	http.Handle("/public/", http.StripPrefix(strings.TrimRight("/public/", "/"), http.FileServer(http.Dir("public"))))

	log.Println("Starting server on port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
