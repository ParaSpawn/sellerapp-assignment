package main

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	scribble "github.com/nanobox-io/golang-scribble"
)

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

func getHash(val []byte) string {
	hasher := sha1.New()
	hasher.Write(val)
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return sha
}

func writeProduct(w http.ResponseWriter, r *http.Request) {
	var product ProductInfo
	json.NewDecoder(r.Body).Decode(&product)
	if err := db.Write("Product", getHash([]byte(product.Product.Title)), product); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func getAllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := db.ReadAll("Product")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		payload := make([]ProductInfo, len(products))
		for i, s := range products {
			json.Unmarshal([]byte(s), &payload[i])
		}
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
