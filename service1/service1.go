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
	// Product info will be scrapped and stored in this variable
	var info ProductInfo
	var URL map[string]string
	// Getting url from request
	json.NewDecoder(r.Body).Decode(&URL)
	// Variable that tracks whether colly's OnError callback was executed
	var collyError error = nil
	// Variable that tracks whether the site was visited and scrapped or not
	var visited bool = false

	// Initializing collector
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

	// Getting product title using the 'id' attribute of 'span' tag
	c.OnHTML("span[id=productTitle]", func(e *colly.HTMLElement) {
		info.Product.Title = strings.Trim(e.Text, "\n")
	})

	// Getting the review count similarly
	c.OnHTML("span[id=acrCustomerReviewText]", func(e *colly.HTMLElement) {
		cnt, _ := strconv.ParseInt(strings.Replace(strings.SplitN(e.Text, " ", 2)[0], ",", "", -1), 10, 64)
		info.Product.ReviewCount = int(cnt)
	})

	// Getting product description and storing it as array of strings
	c.OnHTML("div[id=feature-bullets]", func(e *colly.HTMLElement) {
		e.ForEach("ul", func(_ int, ul *colly.HTMLElement) {
			ul.ForEach("li", func(_ int, li *colly.HTMLElement) {
				li.ForEach("span", func(_ int, span *colly.HTMLElement) {
					info.Product.Description = append(info.Product.Description, strings.Trim(span.Text, "\n"))
				})
			})
		})
	})

	// Getting product price
	c.OnHTML("span[id=priceblock_ourprice]", func(e *colly.HTMLElement) {
		info.Product.Price = e.Text
	})

	// Getting image url
	c.OnHTML("div[class=imgTagWrapper]", func(e *colly.HTMLElement) {
		e.ForEach("img[src]", func(_ int, i *colly.HTMLElement) {
			info.Product.ImageURL = i.Attr("src")
		})
	})
	c.Visit(URL["url"])

	// Checking for errors from colly
	if collyError != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error while scraping: " + collyError.Error()))
	} else if !visited {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Domain not allowed"))
	} else {
		// Storing URL and create/update time in info
		info.URL = URL["url"]
		info.LastUpdated = time.Now()

		// Preparing payload
		payload, _ := json.Marshal(info)

		// Sending payload to second service
		res, err := http.Post("http://service2:8010/writeProductInfo", "application/json", bytes.NewReader(payload))

		// Checking for errors from service 2
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Something went wrong while writing the document: " + err.Error()))
		} else if res.StatusCode == http.StatusInternalServerError {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Something went wrong while writing the document"))
		} else {
			// Encoding the scrapped info in responsewriter
			json.NewEncoder(w).Encode(info)
		}
	}
}

func main() {
	router := mux.NewRouter()
	// Defining endpoints
	router.HandleFunc("/getProductInfo", getProductInfo).Methods("POST")
	// Hosting
	log.Fatal(http.ListenAndServe(":8000", router))
}
