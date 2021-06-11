package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

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
	Dob       string
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
	persons = getAll()

	//insert data to send it to corresponding html file as an object
	myPageVars := HomePageVars{
		Title:   "CRUD",
		Persons: persons,
	}

	//\\for DELETE operation, 1 entry
	//parsing id of delete button on click through form submission

	if r.Method == "POST" {

		tobeDeleteID := r.FormValue("delete")
		idTobeUpdated := r.FormValue("update")
		if tobeDeleteID != "" {
			deleteId, err := strconv.Atoi(tobeDeleteID)
			checkError(err)
			fmt.Printf("to be deleted entry's id: %v %T", deleteId, deleteId)
			fmt.Println(delete(deleteId))
		} else {
			idUpdate, err := strconv.Atoi(idTobeUpdated)
			checkError(err)
			fmt.Printf("to be updated entry's id: %v %T", idUpdate, idUpdate)
		}

	}

	temp, err := template.ParseFiles("homepage.html")
	checkError(err)
	err = temp.Execute(w, myPageVars)
	checkError(err)

	/*


		affect, err := res.RowsAffected()
		fmt.Println(affect, "rows changed")
		//deletion of one entry is complete
	*/

	/*//\\UPDATE one entry

	 */

}

func create(w http.ResponseWriter, r *http.Request) {
	var onePerson Person
	if r.Method == "POST" {
		db := dbConn()
		defer db.Close()

		//layout := "2006-02-06"
		//var dt time.Time
		r.ParseForm()
		fname := r.FormValue("first-name")
		lname := r.FormValue("last-name")
		dob := r.FormValue("date-of-birth")
		fmt.Printf("value=%v type = %T of html date input", dob, dob)
		//dt, err := time.Parse(layout, dob)
		//checkError(err)
		gender := r.FormValue("gender")

		onePerson = Person{
			FirstName: fname,
			LastName:  lname,
			Gender:    gender,
			Dob:       dob,
		}
		var id int
		insertRow := `INSERT INTO person(first_name, last_name, gender, date_of_birth) VALUES($1, $2,$3, TO_DATE($4,'YYYYMMDD')) RETURNING id`
		err := db.QueryRow(insertRow, onePerson.FirstName, onePerson.LastName, onePerson.Gender, onePerson.Dob).Scan(&id)
		checkError(err)
	}

	temp, err := template.ParseFiles("create.html")
	checkError(err)
	err = temp.Execute(w, onePerson)
	checkError(err)
}

func insert(p Person) int {
	var id int

	return id //last inserted row's Id
}

func delete(id int) int64 {
	var effectedRows int64
	db := dbConn()
	defer db.Close()
	stmnt := `DELETE FROM person WHERE id=$1`
	res, err := db.Exec(stmnt, id)
	checkError(err)
	effectedRows, err = res.RowsAffected()
	return effectedRows
}

func update(id int, p Person) int64 { // takes an id and Person and update that rows
	var effectedRows int64
	db := dbConn()
	defer db.Close()
	stmnt := `UPDATE person SET first_name=$2, last_name=$3, gender=$4, date_of_birth=TO_DATE($5) WHERE id=$1`
	res, err := db.Exec(stmnt, id, p.FirstName, p.LastName, p.Gender, p.Dob)
	checkError(err)
	effectedRows, err = res.RowsAffected()
	checkError(err)
	return effectedRows
}

//get all the persons
func getAll() []Person {
	var persons []Person
	//db connection
	db := dbConn()
	defer db.Close()

	//QUERY to get data
	rows, err := db.Query("SELECT * FROM person ORDER BY id DESC")
	checkError(err)

	//get data from db and store it in struct Person
	var (
		id int
		fn string
		ln string
		gn string
		dt string
	)
	//dt.Format("02-01-2006")
	for rows.Next() {
		rows.Scan(&id, &fn, &ln, &gn, &dt)
		persons = append(persons, Person{id, fn, ln, gn, dt})
	}
	return persons
}
func main() {

	http.HandleFunc("/", Index)
	http.HandleFunc("/create", create)
	err := http.ListenAndServe(":9090", nil)
	checkError(err)
}
