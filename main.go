package main

import (
	"errors"
	"fmt"
	"github.com/fpiwowarczyk/simple-key-value-store/logging"
	"github.com/fpiwowarczyk/simple-key-value-store/store"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
)

var logger logging.TransactionLogger

func initializeTransactionLog() error {
	var err error

	//logger, err = logging.NewFileTransactionLogger("transaction.log")
	logger, err = logging.NewPostgresTransactionLogger("localhost", "kvs", "test", "pass1")
	if err != nil {
		return fmt.Errorf("failed to create event logger: %w", err)
	}

	events, errors := logger.ReadEvents()
	e, ok := logging.Event{}, true

	for ok && err == nil {
		select {
		case err, ok = <-errors:
		case e, ok = <-events:
			switch e.EventType {
			case logging.EventDelete:
				err = store.Delete(e.Key)
			case logging.EventPut:
				err = store.Put(e.Key, e.Value)
			}
		}
	}

	logger.Run()

	return err
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

	logger.WritePut(key, string(value))

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

func keyValueDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	err := store.Delete(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logger.WriteDelete(key)
}

func main() {
	router := mux.NewRouter()

	err := initializeTransactionLog()
	if err != nil {
		fmt.Printf("new error %v \n", err)
	}

	router.HandleFunc("/v1/{key}", keyValuePutHandler).Methods("PUT")
	router.HandleFunc("/v1/{key}", keyValueGetHandler).Methods("GET")
	router.HandleFunc("/v1/{key}", keyValueDeleteHandler).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))
}
