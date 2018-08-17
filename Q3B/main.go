package main

import (

	"os"
	"fmt"
	"encoding/json"
	"log"
	"net/http"
	"time"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	rs "../Q3B/utils" 
	"io/ioutil"
	_ "reflect"
	redisConn "../Q3B/model/redis"
	_ "github.com/garyburd/redigo/redis"

)

const urlPrefix string = "http://shorturl.com/"

// Model the record structa
type Record struct {
	URL string `json:"url"`
	ShortenURL string `json:"shorten_url"`
}


// var urllist []Record
var record Record


func (p *Record) GetAllOriginalURL(w http.ResponseWriter , r *http.Request){
	
		keysList, err := redisConn.GetKeys("*")
		if err != nil {
			fmt.Println("Error retrieve keys",err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(keysList)
	
}

// GetOriginalURL fetches the original URL for the given encoded(short) string
func (p *Record ) GetOriginalURL(w http.ResponseWriter, r *http.Request) {
	
	// var results []Record	
	// var res []byte
	vars := mux.Vars(r)
	req_url := urlPrefix + vars["urlhash"] 
	exist, _ := redisConn.Exists(req_url)
	if exist == false {
		w.Write([]byte("Not found this mapping domain"))
	} else {
			res, err := redisConn.Get(req_url)
			if err != nil {
				fmt.Println(err)
			}
			site := "http://" + string(res)
			http.Redirect(w, r, site, 301)
	}
	
}
// GenerateShortURL adds URL to Record dataset and gives back shortened string
func (p *Record) GenerateShortURL(w http.ResponseWriter, r *http.Request) {

	body, _ := ioutil.ReadAll(r.Body)
	record.ShortenURL = urlPrefix + rs.RandStringBytes(8)
	json.Unmarshal(body, &record)
	a1 := []byte(record.URL)
	// val := []byte(record.ShortenURL) // convert string to byte slice
	err := redisConn.SetEX(record.ShortenURL, a1)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(record)
	
}
func init(){
	_ = godotenv.Load(".env")

}

func main() {

	port := os.Getenv("PORT")
	host := os.Getenv("HOST")

	r := mux.NewRouter()
	// Attach an elegant path with handler

	r.Use(rs.RateLimitingMiddleware)
	r.HandleFunc("/all", record.GetAllOriginalURL).Methods("GET")
	r.HandleFunc("/{urlhash:[a-zA-Z0-9]{8}}", record.GetOriginalURL).Methods("GET")
	r.HandleFunc("/submit", record.GenerateShortURL).Methods("POST")
	
	srv := &http.Server{
		Handler: r,
		Addr:    host + ":" + port,
		// Good practice: enforce timeouts
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	fmt.Printf("Server is listening on: %s\n", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
