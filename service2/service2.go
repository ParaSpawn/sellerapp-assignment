package main

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	scribble "github.com/nanobox-io/golang-scribble"
)

// Using a simple document store that stores the JSONs as files
var db *scribble.Driver

type ProductInfo struct {
	URL     string
	Product struct {
		Title, ImageURL, Price string
		Description            []string
		ReviewCount            int
	}
	LastUpdated time.Time
}

// Hash function for titles. The hashed value will be used as id
func getHash(val []byte) string {
	hasher := sha1.New()
	hasher.Write(val)
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return sha
}

func writeProduct(w http.ResponseWriter, r *http.Request) {
	var product ProductInfo
	// Decoding information
	json.NewDecoder(r.Body).Decode(&product)
	// Attempting to store information
	if err := db.Write("Product", getHash([]byte(product.Product.Title)), product); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		// OK response on successful storage
		w.WriteHeader(http.StatusOK)
	}
}

func getAllProducts(w http.ResponseWriter, r *http.Request) {
	// Retrieving all product information
	products, err := db.ReadAll("Product")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		// Preparing payload
		payload := make([]ProductInfo, len(products))
		for i, s := range products {
			json.Unmarshal([]byte(s), &payload[i])
		}
		// Encoding in responsewriter
		json.NewEncoder(w).Encode(payload)
	}
}

func main() {
	db, _ = scribble.New("./db", nil)
	router := mux.NewRouter()
	router.HandleFunc("/writeProductInfo", writeProduct).Methods("POST")
	router.HandleFunc("/getAllProducts", getAllProducts).Methods("GET")
	log.Fatal(http.ListenAndServe(":8010", router))
}
