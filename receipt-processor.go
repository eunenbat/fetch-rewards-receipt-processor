package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"net/http"
	// "strconv"
	"math"
	"time"
	"regexp"


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
	Total 			float64 `json:"total,string"`
	Items 			[]Item `json:"items"`
}

type Item struct {
	ShortDesc 	string `json:"shortDescription"`
	Price 		float64 `json:"price,string"`
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
	fmt.Println("COPY THIS: ", id)
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

func calcPoints(receipt Receipt) int {
	points := 0
	reg := regexp.MustCompile(`[^a-zA-Z0-9 ]+`)
	nonAlpha := reg.ReplaceAllString(receipt.Retailer, "")
	nonAlpha = strings.ReplaceAll(nonAlpha, " ", "")

	// Rule 1
	points += len(nonAlpha)
	fmt.Println("R1: ", points)
	// Rule 2
	if receipt.Total == float64(int(receipt.Total)) {
		points += 50
	}
	fmt.Println("R2: ",points)
	// Rule 3
	if int(receipt.Total * 100) % 25 == 0 {
		points += 25
	}
	fmt.Println("R3: ",points)
	// Rule 4
	points += 5 * ((len(receipt.Items) - (len(receipt.Items) % 2)) / 2)
	fmt.Println("R4: ",points)
	// Rule 5
	for _, item := range receipt.Items {
        descLen := len(strings.TrimSpace(item.ShortDesc))
        if descLen % 3 == 0 {
            points += int(math.Ceil(item.Price * 0.2))
        }
    }
	fmt.Println("R5: ",points)
	// Rule 6
	date, _ := time.Parse("2006-01-02", receipt.PurchaseDate)
	if date.Day() % 2 != 0 {
		points += 6
	}
	fmt.Println("R6: ",points)
	// Rule 7
	pTime, _ := time.Parse("15:04", receipt.PurchaseTime)
	
	twoPM := time.Date(pTime.Year(), pTime.Month(), pTime.Day(), 14, 0, 0, 0, time.UTC)
	fourPM := time.Date(pTime.Year(), pTime.Month(), pTime.Day(), 16, 0, 0, 0, time.UTC)
	if pTime.After(twoPM) && pTime.Before(fourPM) {
		points += 10
	}
	fmt.Println("R7 (Total): ",points)
	return points
}

func getPoints(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

 	id := params["id"]
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }
	fmt.Println("getPoints output (points): ", receipts[id])
	json.NewEncoder(w).Encode(map[string]int{"points": receipts[id]})
}

func randoms(w http.ResponseWriter, r *http.Request) {
	// w.Write([]byte("This is my home page"))
	json.NewEncoder(w).Encode(map[string]int{"points": 900})
}

func main() {
	router := mux.NewRouter()
	// handle requests
	router.HandleFunc("/", randoms).Methods("GET")
	router.HandleFunc("/receipts/process", processReceipts).Methods("POST")
	router.HandleFunc("/receipts/{id}/points", getPoints).Methods("GET")
	// router.HandleFunc("/hi", randoms).Methods("GET")

	fmt.Println("Server is running on port 5670")
	log.Fatal(http.ListenAndServe(":5670", router))
}
