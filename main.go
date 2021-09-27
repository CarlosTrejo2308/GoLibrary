package main

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

// AllowedRoles is a map of resources that are allowed to be served from API.
var AllowedResources = map[string]bool{
	"books":   true,
	"authors": true,
	"genres":  true,
}

// Book is a book.
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
	// Create a new HTTP server.
	http.Handle("/", http.HandlerFunc(ExampleHandler))
	http.HandleFunc("/books", basicAuth(books))
	http.HandleFunc("/book/", basicAuth(book))

	http.ListenAndServe(":8080", nil)
}

// basicAuth is a middleware that checks for a valid username and password.
func basicAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the username and password from the request
		// Authorization header. If no Authentication header is present
		// or the header value is invalid, then the 'ok' return value
		// will be false.
		username, password, ok := r.BasicAuth()
		if ok {
			// Calculate SHA-256 hashes for the provided and expected
			// usernames and passwords.
			usernameHash := sha256.Sum256([]byte(username))
			passwordHash := sha256.Sum256([]byte(password))
			expectedUsernameHash := sha256.Sum256([]byte("admin"))
			expectedPasswordHash := sha256.Sum256([]byte("pwd1"))

			// Use the subtle.ConstantTimeCompare() function to check if
			// the provided username and password hashes equal the
			// expected username and password hashes. ConstantTimeCompare
			// will return 1 if the values are equal, or 0 otherwise.
			// Importantly, we should to do the work to evaluate both the
			// username and password before checking the return values to
			// avoid leaking information.
			usernameMatch := (subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1)
			passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1)

			// If the username and password are correct, then call
			// the next handler in the chain. Make sure to return
			// afterwards, so that none of the code below is run.
			if usernameMatch && passwordMatch {
				next.ServeHTTP(w, r)
				return
			}
		}

		// If the Authentication header is not present, is invalid, or the
		// username or password is wrong, then set a WWW-Authenticate
		// header to inform the client that we expect them to use basic
		// authentication and send a 401 Unauthorized response.
		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}

// book is a handler for a single book.
func book(w http.ResponseWriter, r *http.Request) {

	// Get the book id from the URL.
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

// getBook returns a book.
func getBook(w http.ResponseWriter, r *http.Request, id int) {

	// Check if the book id is valid.
	maxBooks := len(Books)
	if id >= maxBooks {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "<h1>Not Found</h1>")
		return
	}

	// Get the book.
	book := Books[id]

	// Marshal the book into JSON.
	response, err := json.Marshal(book)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "<h1>Internal Server Error</h1>")
		return
	}

	// Send the JSON response.
	fmt.Fprint(w, string(response))
}

// putBook updates a book.
func putBook(w http.ResponseWriter, r *http.Request, id int) {

	// Check if the book id is valid.
	if id >= len(Books) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "<h1>Not Found</h1>")
		return
	}

	var book Book

	// Read the request body.
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "<h1>Bad Request</h1>")
		return
	}

	// Unmarshal the JSON into the book struct.
	if err := json.Unmarshal(body, &book); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "<h1>Bad Request</h1>")
		return
	}

	// Update the book.
	Books[id] = book
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

// deleteBook deletes a book.
func deleteBook(w http.ResponseWriter, r *http.Request, id int) {
	// Check if the book id is valid.
	if id >= len(Books) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "<h1>Not Found</h1>")
		return
	}

	// Delete the book.
	Books = append(Books[:id], Books[id+1:]...)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

// books is a handler for all books.
func books(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getBooks(w, r)
	case "POST":
		postBooks(w, r)
	}
}

// getBooks returns all books.
func getBooks(w http.ResponseWriter, r *http.Request) {
	// Marshal the books into JSON.
	respose, err := json.Marshal(Books)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "<h1>Internal Server Error</h1>")
		return
	}

	// Send the JSON response.
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(respose))
}

// postBooks adds a book.
func postBooks(w http.ResponseWriter, r *http.Request) {
	var book Book

	// Read the request body.
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "<h1>Bad Request</h1>")
		return
	}

	// Unmarshal the JSON into the book struct.
	if err := json.Unmarshal(body, &book); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "<h1>Bad Request, error on unmarshal</h1>")
		return
	}

	// Add the book.
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
