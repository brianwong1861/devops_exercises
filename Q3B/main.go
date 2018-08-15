package main

import (
	"strings"
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

)

const urlPrefix string = "http://shorturl.com/"

// Model the record structa
type Record struct {
	URL string `json:"url"`
	ShortenURL string `json:"shorten_url"`
	StartTime int64 `json:"start_time"`
	IPAddr string `json:"ip_addr"`
}


var urllist []Record
var record Record


func (p *Record) GetAllOriginalURL(w http.ResponseWriter , r *http.Request){
			
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(urllist)
}

// GetOriginalURL fetches the original URL for the given encoded(short) string
func (p *Record ) GetOriginalURL(w http.ResponseWriter, r *http.Request) {
	
	var results []Record	
	vars := mux.Vars(r)
	req_url := urlPrefix + vars["urlhash"] 
	w.Header().Set("Content-Type", "application/json")
	for _, v := range urllist {
		if req_url == v.ShortenURL {
			results = append(results, v)
		} 
	}
	site := "http://" + results[0].URL
	http.Redirect(w, r, site, 301)
	
}
// GenerateShortURL adds URL to Record dataset and gives back shortened string
func (p *Record) GenerateShortURL(w http.ResponseWriter, r *http.Request) {

	body, _ := ioutil.ReadAll(r.Body)
	record.ShortenURL = urlPrefix + rs.RandStringBytes(8)
	json.Unmarshal(body, &record)
	urllist = append(urllist, record)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(record)
	
}
func RateLimiting(w http.ResponseWriter, r *http.Request) bool {

	ipAddr_trimmed := strings.Split(r.RemoteAddr, ":")[0] //Split Remote IP ADDR into IP and Port
	// message = make(map[string]string)

	if ( record.IPAddr == "" ) && ( record.StartTime == 0 ) {
		record.IPAddr = ipAddr_trimmed
		record.StartTime = time.Now().Unix()
	} else if ( time.Now().Unix() - record.StartTime ) <= 5 {
		w.Write([]byte("Exceeded rate limit\n"))
		record.StartTime = time.Now().Unix() 
		return false
	}
	record.StartTime = time.Now().Unix() 
	return true
}

func rateLimitingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		val := RateLimiting(w, r)
		if val == true {
			next.ServeHTTP(w, r)
		}
        // Call the next handler, which can be another middleware in the chain, or the final handler.
    })
}

func init(){
	_ = godotenv.Load(".env")
}

func main() {

	port := os.Getenv("PORT")
	host := os.Getenv("HOST")

	r := mux.NewRouter()
	// Attach an elegant path with handler

	r.Use(rateLimitingMiddleware)
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
