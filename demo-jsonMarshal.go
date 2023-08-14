package main

import (
	"encoding/json"
	"fmt"
)

type employee struct {
	ID           int
	EmployeeName string
	Tel          string
	Email        string
}

func main() {

	data, _ := json.Marshal(&employee{101, "Sirasit Boonklang", "0900000000", "sirasit@mail.com"})
	fmt.Println(string(data))

}
