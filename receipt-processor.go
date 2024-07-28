package main

import (
	"encoding/json"
	"fmt"
	"log"
	// "strings"
	"net/http"
	// "strconv"
	"math"
	"time"


	"github.com/gorilla/mux"
	"github.com/google/uuid"
)

// type Receipts struct {
// 	Receipts []Receipt `json:"Receipts"`
// }
var receipts = make(map[string]int)

type Receipt struct {
	Retailer 		string `json:"retailer"`
	PurchaseDate 	string `json:"purchaseDate"`
	PurchaseTime 	string `json:"purchaseTime"`
	Total 			float64 `json:"total"`
	Items 			[]Item `json:"items"`
}

type Item struct {
	ShortDesc 	string `json:"shortDescription"`
	Price 		float64 `json:"price"`
}

func processReceipts(w http.ResponseWriter, r *http.Request) {
	var receipt Receipt
	err := json.NewDecoder(r.Body).Decode(&receipt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
 	}

	id := uuid.New().String()
	receipts[id] = calcPoints(receipt)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

func calcPoints(receipt Receipt) int {
	points := 0

	// Rule 1
	points += len(receipt.Retailer)
	// Rule 2
	if receipt.Total == float64(int(receipt.Total)) {
		points += 50
	}
	// Rule 3
	if int(receipt.Total * 100) % 25 == 0 {
		points += 25
	}
	// Rule 4
	points += 5 * ((len(receipt.Items) - (len(receipt.Items) % 2)) / 2)
	// Rule 5
	for _, item := range receipt.Items {
        descLen := len(item.ShortDesc)
        if descLen % 3 == 0 {
            points += int(math.Ceil(item.Price * 0.2))
        }
    }
	// Rule 6
	// const shortForm = "2006-01-02"
	date, _ := time.Parse("2006-01-02", receipt.PurchaseDate)
	if date.Day() % 2 != 0 {
		points += 6
	}
	// Rule 7
	pTime, _ := time.Parse("13:01", receipt.PurchaseTime)
	twoPM := time.Date(pTime.Year(), pTime.Month(), pTime.Day(), 14, 0, 0, 0, time.UTC)
	fourPM := time.Date(pTime.Year(), pTime.Month(), pTime.Day(), 16, 0, 0, 0, time.UTC)
	if pTime.After(twoPM) && pTime.Before(fourPM) {
		points += 10
	}

	return points
}

func getPoints(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

 	id, err := params["id"]
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(map[string]int{"points": receipts[id]})
}

func main() {
	router := mux.NewRouter()
	// handle requests
	router.HandleFunc("/receipts/process", processReceipts).Methods("POST")
	router.HandleFunc("/receipts/{id}/points", getPoints).Methods("GET")

	fmt.Println("Server is running on port 1000")
	log.Fatal(http.ListenAndServe(":1000", nil))
}