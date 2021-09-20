package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

var AllowedResources = map[string]bool{
	"books":   true,
	"authors": true,
	"genres":  true,
}

type Book struct {
	Titulo    string `json:"titulo"`
	Id_Autor  int    `json:"id_autor"`
	Id_gereno int    `json:"id_gereno"`
}

var Books = []Book{
	{
		Titulo:    "Lo que el viento se llevo",
		Id_Autor:  2,
		Id_gereno: 2,
	},
	{
		Titulo:    "El señor de los anillos",
		Id_Autor:  1,
		Id_gereno: 1,
	},
	{
		Titulo:    "La Odisea",
		Id_Autor:  1,
		Id_gereno: 3,
	},
}

func main() {
	http.Handle("/", http.HandlerFunc(ExampleHandler))
	http.ListenAndServe(":8080", nil)
}

func ExampleHandler(w http.ResponseWriter, r *http.Request) {
	// Get the resource from the request
	params := r.URL.Query()
	resource := params.Get("resource_type")
	resourceId := params.Get("resource_id")

	// Check if the resource is allowed
	if !AllowedResources[resource] {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "<h1>Not Found</h1>")
		return
	}

	// Handle the request
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		switch resource {
		case "books":

			if resourceId != "" {
				// Get a single book
				maxBooks := len(Books)

				bookid, err := strconv.Atoi(resourceId)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					fmt.Fprint(w, "<h1>Bad Request</h1>")
					return
				}

				if bookid >= maxBooks {
					w.WriteHeader(http.StatusNotFound)
					fmt.Fprint(w, "<h1>Not Found</h1>")
					return
				}

				book := Books[bookid]
				response, err := json.Marshal(book)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprint(w, "<h1>Internal Server Error</h1>")
					return
				}

				fmt.Fprint(w, string(response))

				return
			}

			respose, err := json.Marshal(Books)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, "<h1>Internal Server Error</h1>")
				return
			}

			fmt.Fprint(w, string(respose))
		}
		//case "authors":

	case http.MethodPost:
		w.Write([]byte("POST"))
	case http.MethodPut:
		w.Write([]byte("PUT"))
	case http.MethodDelete:
		w.Write([]byte("DELETE"))
	default:
		w.Header().Set("Allow", "GET, POST, PUT, DELETE")
		w.Write([]byte("UNKNOWN METHOD"))
	}

}
