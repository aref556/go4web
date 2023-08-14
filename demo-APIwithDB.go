package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var Db *sql.DB

const coursePath = "courses"
const basePath = "/api"

type Course struct {
	CourseID   int     `json: "courseid"`
	Coursename string  `json: "coursename"`
	Price      float64 `json: "price"`
	ImageURL   string  `json: "image_url"`
}

// var courseList []Course

func getCourseList() ([]Course, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	results, err := Db.QueryContext(ctx, `SELECT
	courseid,
	coursename,
	price,
	image_url
	FROM courseonline`)

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	defer results.Close()
	courses := make([]Course, 0)
	for results.Next() {
		var course Course
		results.Scan(&course.CourseID,
			&course.Coursename,
			&course.ImageURL,
			&course.Price)

		courses = append(courses, course)
	}
	return courses, nil

}

func insertCourse(course Course) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result, err := Db.ExecContext(ctx, `INSERT INTO courseonline
	(courseid,
	coursename,
	price,
	image_url
	) VALUES (?, ?, ?, ?)`,
		course.CourseID,
		course.Coursename,
		course.Price,
		course.ImageURL)
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}

	insertID, err := result.LastInsertId()
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}

	return int(insertID), nil
}

func getCourse(courseid int) (*Course, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	row := Db.QueryRowContext(ctx, `SELECT
	courseid,
	coursename,
	price,
	image_url
	FROM courseonline
	WHERE courseid = ?`, courseid)

	course := &Course{}
	err := row.Scan(
		&course.CourseID,
		&course.Coursename,
		&course.Price,
		&course.ImageURL,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}

	return course, nil
}

func removeProduct(courseID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := Db.ExecContext(ctx, `DELETE FROM courseonline where courseid = ?`, courseID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func handleCourse(w http.ResponseWriter, r *http.Request) {
	urlPathSegment := strings.Split(r.URL.Path, fmt.Sprintf("%s/", coursePath))
	if len(urlPathSegment[1:]) > 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	courseID, err := strconv.Atoi(urlPathSegment[len(urlPathSegment)-1])
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		course, err := getCourse(courseID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if course == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		j, err := json.Marshal(course)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		_, err = w.Write(j)
		if err != nil {
			log.Fatal(err)
		}
	case http.MethodDelete:
		fmt.Println("CourseID : ", courseID)
		err := removeProduct(courseID)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handleCourses(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		courseList, err := getCourseList()

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		j, err := json.Marshal(courseList)
		if err != nil {
			log.Fatal(err)
		}

		_, err = w.Write(j)
		if err != nil {
			log.Fatal(err)
		}
	case http.MethodPost:
		var course Course
		err := json.NewDecoder(r.Body).Decode(&course)

		if err != nil {
			log.Print("1: ", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		CourseID, err := insertCourse(course)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprintf(`{"courseid" : %d}`, CourseID)))

	case http.MethodOptions:
		return

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

func corsMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE")
		w.Header().Add("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, X-Custom-Header, x-requested-with, Accept-Encoding")
		handler.ServeHTTP(w, r)
	})
}

func setupRoutes(apiBasePath string) {
	coursesHandler := http.HandlerFunc(handleCourses)
	courseHandler := http.HandlerFunc(handleCourse)

	http.Handle(fmt.Sprintf("%s/%s", apiBasePath, coursePath), corsMiddleware(coursesHandler))
	http.Handle(fmt.Sprintf("%s/%s/", apiBasePath, coursePath), corsMiddleware(courseHandler))

}

func setupDB() {
	var err error
	Db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/coursedb")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(Db)
	Db.SetConnMaxLifetime(time.Minute * 3)
	Db.SetMaxOpenConns(10)
	Db.SetMaxIdleConns(10)
}

func main() {

	setupDB()
	setupRoutes(basePath)
	log.Fatal(http.ListenAndServe(":5000", nil))

}
