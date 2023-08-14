package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type employee struct {
	ID           int
	EmployeeName string
	Tel          string
	Email        string
}

func main() {
	e := employee{}

	err := json.Unmarshal([]byte(`{"ID": 101, "EmployeeName": "Siratsit Boonklan", "Tel": "0900000000", "Email": "sirasit@mail.com"}`), &e)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(e.EmployeeName)
}
