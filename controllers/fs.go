package controllers

import (
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func generateName() string {
	var name string
	for i := 0; i < 30; i++ {
		name += strconv.Itoa(rand.Intn(9))
	}
	return name
}

// PostFile to the public directory
func PostFile(w http.ResponseWriter, r *http.Request) {
	var directory string

	r.ParseMultipartForm(0)
	needAuth := r.FormValue("need_auth")

	file, handler, err := r.FormFile("file")
	defer file.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	if needAuth == "true" {
		directory = "./private/"
	} else {
		directory = "./public/"
	}
	os.MkdirAll(directory, 0777)
	rand.Seed(time.Now().UnixNano())
	filename := generateName() + filepath.Ext(handler.Filename)
	f, err := os.OpenFile(directory+filename, os.O_WRONLY|os.O_CREATE, 0666)
	defer f.Close()
	io.Copy(f, file)
	json.NewEncoder(w).Encode(filename)
}
