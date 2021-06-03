package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gorilla/mux"
)

type ProductInfo struct {
	URL     string
	Product struct {
		Title, ImageURL, Price string
		Description            []string
		ReviewCount            int
	}
	LastUpdated time.Time
}

func getProductInfo(w http.ResponseWriter, r *http.Request) {
	var info ProductInfo
	var URL map[string]string
	json.NewDecoder(r.Body).Decode(&URL)
	var collyError error = nil
	var visited bool = false
	c := colly.NewCollector(
		colly.AllowedDomains("amazon.in", "www.amazon.in"),
		colly.MaxBodySize(0),
		colly.AllowURLRevisit(),
	).Clone()

	c.OnRequest(func(r *colly.Request) {
		visited = true
	})
	c.OnError(func(response *colly.Response, err error) {
		collyError = err
	})
	c.OnHTML("span[id=productTitle]", func(e *colly.HTMLElement) {
		info.Product.Title = strings.Trim(e.Text, "\n")
	})
	c.OnHTML("span[id=acrCustomerReviewText]", func(e *colly.HTMLElement) {
		cnt, _ := strconv.ParseInt(strings.Replace(strings.SplitN(e.Text, " ", 2)[0], ",", "", -1), 10, 64)
		info.Product.ReviewCount = int(cnt)
	})
	c.OnHTML("div[id=feature-bullets]", func(e *colly.HTMLElement) {
		e.ForEach("ul", func(_ int, ul *colly.HTMLElement) {
			ul.ForEach("li", func(_ int, li *colly.HTMLElement) {
				li.ForEach("span", func(_ int, span *colly.HTMLElement) {
					info.Product.Description = append(info.Product.Description, strings.Trim(span.Text, "\n"))
				})
			})
		})
	})
	c.OnHTML("span[id=priceblock_ourprice]", func(e *colly.HTMLElement) {
		info.Product.Price = e.Text
	})
	c.OnHTML("div[class=imgTagWrapper]", func(e *colly.HTMLElement) {
		e.ForEach("img[src]", func(_ int, i *colly.HTMLElement) {
			info.Product.ImageURL = i.Attr("src")
		})
	})
	c.Visit(URL["url"])
	if collyError != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error while scraping: " + collyError.Error()))
	} else if !visited {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Domain not allowed"))
	} else {
		info.URL = URL["url"]
		info.LastUpdated = time.Now()
		payload, _ := json.Marshal(info)
		res, err := http.Post("http://service2:8010/writeProductInfo", "application/json", bytes.NewReader(payload))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Something went wrong while writing the document: " + err.Error()))
		} else if res.StatusCode == http.StatusInternalServerError {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Something went wrong while writing the document"))
		} else {
			json.NewEncoder(w).Encode(info)
		}
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("made it"))
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/getProductInfo", getProductInfo).Methods("POST")
	router.HandleFunc("/", homePage).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))
}
