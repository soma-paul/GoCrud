package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "test123"
	DB_NAME     = "test"
	TABLE_NAME  = "person"
	DB_DRIVER   = "postgres"
)

func checkError(err error) {
	if err != nil {
		log.Fatal(err.Error())
		//log.Panic(err)
	}
}

type Person struct {
	Id        int
	FirstName string
	LastName  string
	Gender    string
	Dob       time.Time
}

type HomePageVars struct {
	Title   string
	Persons []Person
}

func dbConn() (db *sql.DB) { //return a batabase with necessary info
	dbInfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open(DB_DRIVER, dbInfo)
	checkError(err)
	fmt.Print("Connection to test database is successful!")
	return db
}

func Index(w http.ResponseWriter, r *http.Request) {
	var persons []Person
	db := dbConn()

	defer db.Close()

	/*//inserting a row
	insertRow := fmt.Sprintf("INSERT INTO %s(first_name, last_name, gender, date_of_birth) VALUES('arpon', 'anaytul', 'male', date'1998-07-19')", TABLE_NAME)
	db.QueryRow(insertRow) */

	//QUERY to get data
	rows, err := db.Query("SELECT * FROM person ORDER BY id DESC")
	checkError(err)

	var (
		id int
		fn string
		ln string
		gn string
		dt time.Time
	)
	//dt.Format("02-01-2006")
	for rows.Next() {

		rows.Scan(&id, &fn, &ln, &gn, &dt)
		persons = append(persons, Person{id, fn, ln, gn, dt})
		//fmt.Printf("length=%v capacity=%v\n", len(persons), cap(persons))

	}
	//fmt.Printf("%v %T length=%v capacity=%v\n", persons, persons, len(persons), cap(persons))

	checkError(err)

	myPageVars := HomePageVars{
		Title:   "CRUD",
		Persons: persons,
	}

	temp, err := template.ParseFiles("homepage.html")
	checkError(err)
	err = temp.Execute(w, myPageVars)
	checkError(err)

}

func create(w http.ResponseWriter, r *http.Request) {
	var onePerson Person
	layout := "02-01-2006"
	//var dt time.Time
	r.ParseForm()
	fname := r.FormValue("first-name")
	lname := r.FormValue("last-name")
	dob := r.FormValue("date-of-birth")
	dt, err := time.Parse(layout, dob)
	gender := r.FormValue("gender")

	fmt.Println(fname)
	onePerson = Person{
		FirstName: fname,
		LastName:  lname,
		Gender:    gender,
		Dob:       dt,
	}

	temp, err := template.ParseFiles("create.html")
	checkError(err)

	err = temp.Execute(w, onePerson)
}

func main() {

	http.HandleFunc("/", Index)
	http.HandleFunc("/create", create)
	err := http.ListenAndServe(":9090", nil)
	checkError(err)
}
