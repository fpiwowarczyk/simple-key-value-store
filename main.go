package main

import (
	"errors"
	"github.com/fpiwowarczyk/simple-key-value-store/store"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
)

func helloGoHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello gorilla/mux!\n"))
}

func keyValuePutHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	value, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w,
			err.Error(),
			http.StatusInternalServerError)
		return
	}
	err = store.Put(key, string(value))
	if err != nil {
		http.Error(w,
			err.Error(),
			http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func keyValueGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	value, err := store.Get(key)
	if errors.Is(err, store.ErrorNoSuchKey) {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(value))
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/v1/{key}", keyValuePutHandler).Methods("PUT")
	router.HandleFunc("/v1/{key}", keyValueGetHandler).Methods("GET")
	router.HandleFunc("/", helloGoHandler)

	log.Fatal(http.ListenAndServe(":8080", router))
}
