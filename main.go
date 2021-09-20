package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	http.HandleFunc("/books", books)
	http.HandleFunc("/book/", book)

	http.ListenAndServe(":8080", nil)
}

func book(w http.ResponseWriter, r *http.Request) {
	idstr := r.URL.Path[len("/book/"):]
	id, err := strconv.Atoi(idstr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "<h1>Bad Request</h1>")
		return
	}

	switch r.Method {
	case http.MethodGet:
		getBook(w, r, id)
	case http.MethodPut:
		putBook(w, r, id)
	case http.MethodDelete:
		deleteBook(w, r, id)
	}
}

func getBook(w http.ResponseWriter, r *http.Request, id int) {
	// Get a single book
	maxBooks := len(Books)

	if id >= maxBooks {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "<h1>Not Found</h1>")
		return
	}

	book := Books[id]
	response, err := json.Marshal(book)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "<h1>Internal Server Error</h1>")
		return
	}

	fmt.Fprint(w, string(response))
}

func putBook(w http.ResponseWriter, r *http.Request, id int) {
	var book Book

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "<h1>Bad Request</h1>")
		return
	}

	if err := json.Unmarshal(body, &book); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "<h1>Bad Request</h1>")
		return
	}

	if id >= len(Books) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "<h1>Not Found</h1>")
		return
	}

	Books[id] = book
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

func deleteBook(w http.ResponseWriter, r *http.Request, id int) {
	if id >= len(Books) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "<h1>Not Found</h1>")
		return
	}

	Books = append(Books[:id], Books[id+1:]...)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

func books(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getBooks(w, r)
	case "POST":
		postBooks(w, r)
	}
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	respose, err := json.Marshal(Books)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "<h1>Internal Server Error</h1>")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(respose))
}

func postBooks(w http.ResponseWriter, r *http.Request) {
	var book Book

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "<h1>Bad Request</h1>")
		return
	}

	if err := json.Unmarshal(body, &book); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "<h1>Bad Request, error on unmarshal</h1>")
		return
	}

	Books = append(Books, book)
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, len(Books)-1)
}

// This is ussing url/?resource=books
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
