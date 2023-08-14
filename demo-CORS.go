package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
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

func findID(ID int) (*Course, int) {
	for i, course := range CourseList {
		if course.ID == ID {
			return &course, i
		}
	}
	return nil, 0
}

func courseHandler(w http.ResponseWriter, r *http.Request) {
	urlPathSegment := strings.Split(r.URL.Path, "course/")
	ID, err := strconv.Atoi(urlPathSegment[len(urlPathSegment)-1])
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	course, listItemIndex := findID(ID)
	if course == nil {
		http.Error(w, fmt.Sprintf("no course from this id %d", ID), http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		CourseJSON, err := json.Marshal(course)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(CourseJSON)

	case http.MethodPut:
		var updateCourse Course
		Bytebody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(Bytebody, &updateCourse)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if updateCourse.ID != ID {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		course = &updateCourse
		CourseList[listItemIndex] = *course
		w.WriteHeader(http.StatusOK)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

func coursesHandler(w http.ResponseWriter, r *http.Request) {
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

func enableCorsMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE")
		w.Header().Add("Access-Control-Allow-Headers", "Accept, Content-Type, X-Custom-Header, x-requested-with")
		handler.ServeHTTP(w, r)
	})
}

func main() {

	courseItemHandler := http.HandlerFunc(courseHandler)
	courseListHandler := http.HandlerFunc(coursesHandler)

	http.Handle("/course/", enableCorsMiddleware(courseItemHandler))
	http.Handle("/course", enableCorsMiddleware(courseListHandler))
	http.ListenAndServe(":5000", nil)
}
