package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	// initEvents()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/image/{id}", getImage).Methods("GET")
	router.HandleFunc("/image", saveImage).Methods("POST")
	fmt.Println("Server listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

func getImage(w http.ResponseWriter, r *http.Request) {
	imageName := mux.Vars(r)["id"]
	fileBytes, err := ioutil.ReadFile("storage/" + imageName)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(fileBytes)
	return
}

func saveImage(w http.ResponseWriter, r *http.Request) {
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error while reading file")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer file.Close()

	// Create the storage folder if it doesn't
	// already exist
	err = os.MkdirAll("./storage", os.ModePerm)
	if err != nil {
		fmt.Println("Error making directory")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create a new file in the storage directory
	dst, err := os.Create(fmt.Sprintf("./storage/%d%s", time.Now().UnixNano(), filepath.Ext(fileHeader.Filename)))
	if err != nil {
		fmt.Println("Error while creating file")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer dst.Close()

	// Copy the uploaded file to the filesystem
	// at the specified destination
	_, err = io.Copy(dst, file)
	if err != nil {
		fmt.Println("Error while copying file")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, dst.Name())
}
