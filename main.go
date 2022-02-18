package main

import (
	"encoding/json"
	"log"
	"net/http"

	// "github.com/lithammer/shortuuid"
	"github.com/gorilla/mux"
)

type Foo struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

var idMap = map[string]Foo{}

func getFoo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// now get the id
	vars := mux.Vars(r)
	key := vars["id"]
	foo, exists := idMap[key]
	if exists {
		b, err := json.Marshal(foo)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write(b)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func postFoo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "post called"}`))
}

func deleteFoo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// now get the id
	vars := mux.Vars(r)
	key := vars["id"]
	delete(idMap, key)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "delete called"}`))
}

func main() {
	r := mux.NewRouter().StrictSlash(true)

	idMap["a"] = Foo{
		Id:   "a",
		Name: "Al",
	}
	r.HandleFunc("/foo/{id}", getFoo).Methods(http.MethodGet)
	r.HandleFunc("/foo", postFoo).Methods(http.MethodPost)
	r.HandleFunc("/foo/{id}", deleteFoo).Methods(http.MethodDelete)
	log.Fatal(http.ListenAndServe(":8080", r))
}
