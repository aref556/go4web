package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Course struct {
	ID         int     `json: "id"`
	Name       string  `json: "name"`
	Price      float64 `json: "price"`
	Instructor string  `json: "instructor"`
}

var CourseList []Course

func init() {
	// fmt.Println("Hello from init")
	CourseJSON := `[
		{
			"id": 1,
			"name": "Python",
			"price": 2590,
			"instructor": "BorntoDev"
		},
		{
			"id": 2,
			"name": "JavaScript",
			"price": 0,
			"instructor": "BorntoDev"
		},
		{
			"id": 3,
			"name": "SQL",
			"price": 0,
			"instructor": "BorntoDev"
		}
	]`

	err := json.Unmarshal([]byte(CourseJSON), &CourseList)

	if err != nil {
		log.Fatal(err)
	}
}

func getNextID() int {
	highestID := -1

	for _, course := range CourseList {
		if highestID < course.ID {
			highestID = course.ID
		}
	}
	return highestID + 1

}
func courseHandler(w http.ResponseWriter, r *http.Request) {
	courseJSON, err := json.Marshal(CourseList)
	// fmt.Println("courseJSON = ", string(courseJSON))

	switch r.Method {
	case http.MethodGet:
		fmt.Println("Hi From Get Method")

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(courseJSON)

	case http.MethodPost:
		fmt.Println("Hi From Post Method")
		var newCourse Course

		Bodybyte, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// fmt.Println("Bodybyte = ", string(Bodybyte))

		err = json.Unmarshal(Bodybyte, &newCourse)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// fmt.Println("newCourse.ID = ", newCourse.ID)

		if newCourse.ID != 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		newCourse.ID = getNextID()
		CourseList = append(CourseList, newCourse)
		w.WriteHeader(http.StatusCreated)
		return

	}

}

func main() {
	http.HandleFunc("/course", courseHandler)
	http.ListenAndServe(":5000", nil)
}
