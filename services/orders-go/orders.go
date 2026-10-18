// Package main is a tiny orders API: GET /total?qty=3&price=250 returns the order total in cents.
package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// Total returns qty * price, refusing negative inputs.
func Total(qty, price int) (int, error) {
	if qty < 0 || price < 0 {
		return 0, fmt.Errorf("qty and price must be non-negative")
	}
	return qty * price, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	qty, err1 := strconv.Atoi(r.URL.Query().Get("qty"))
	price, err2 := strconv.Atoi(r.URL.Query().Get("price"))
	if err1 != nil || err2 != nil {
		http.Error(w, "qty and price must be integers", http.StatusBadRequest)
		return
	}
	total, err := Total(qty, price)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "%d\n", total)
}

func main() {
	http.HandleFunc("/total", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
