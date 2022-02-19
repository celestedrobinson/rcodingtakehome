/*
A simple package that creates a server for a coding exercise for a job application process


Given more time and a more refreshed memory when it comes to Golang, I would write formal unit tests. I would also set up some helpful logging. Another change would be adding
a response body if the name is missing on the post request.
I tested the code incrementally using Postman with a unit test inspired approach. First I created a function outline (get, post and delete functions). I started with get,
the first step was returning a foo (any foo) in json form. After, I created the global map, where I hardcoded one key-value pair for testing.
Then I added code to get the id value in the request. After I was satisfied with the functionality of the get function, I moved on to the delete. This one was straight
forward, following the same steps as before, just adding a line to delete the record if found and changing the status code.
For post, I needed to unmarshal the json into the foo. Then create and add the uid. After that return that in json form.

After each of these listed steps, I was sure to check that the output was as expected, running a through requests through Postman.
*/
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lithammer/shortuuid"
)

type Foo struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

type NoNameError struct {
	Body []byte
}

func (e *NoNameError) Error() string {
	return fmt.Sprintf("No name found in body: %v", e.Body)
}

var fooMap = map[string]Foo{}

// For a GET request, attempts to find a foo with the id in the request
func getFoo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r) // get the ID values from the variables in the address
	key := vars["id"]
	foo, exists := fooMap[key]
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	b, err := json.Marshal(foo) // convert the foo to JSON
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(b)
	return

}

// The steps to create and add a new foo
func setFoo(body []byte) (foo *Foo, err error) {
	if err := json.Unmarshal(body, &foo); err != nil {
		return nil, err
	}
	if foo.Name == "" {
		return nil, &NoNameError{Body: body}
	}
	id := shortuuid.New()
	foo.Id = id
	fooMap[foo.Id] = *foo
	return foo, nil
}

// On a POST request, takes steps to add a foo to fooMap using the "name" in the received JSON, returning the foo in the response
func postFoo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	foo, err := setFoo(body)
	if err != nil {
		switch err.(type) {
		case *NoNameError:
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
		return
	}

	b, err := json.Marshal(foo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

// On a DELETE request removes the record with the given ID, if found
func deleteFoo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	_, exists := fooMap[key]
	if exists {
		delete(fooMap, key)
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func main() {
	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/foo/{id}", getFoo).Methods(http.MethodGet)
	r.HandleFunc("/foo", postFoo).Methods(http.MethodPost)
	r.HandleFunc("/foo/{id}", deleteFoo).Methods(http.MethodDelete)
	log.Fatal(http.ListenAndServe(":8080", r))
}
