package main

import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	// "exec" // For running perl script to get info from 90 number
)

// Company type information
type Company struct {
	Name string `json:"Name,omitempty"`
}

type Student struct{
	FirstName string `json:"FirstName,omitempty"`
	LastName string `json:"LastName,omitempty"`
	Class string `json:"Class,omitempty"` // Freshman, Sophomore, Junior, Senior
	Major string `json:"Major,omitempty"`
}

var companyList []Company

func main() {
	log.Print("Starting server...")

	router := mux.NewRouter()
	
	router.HandleFunc("/company_list", GetCompanyList).Methods("GET")
	router.HandleFunc("/get_student", GetStudent).Methods("GET")
	router.HandleFunc("/put_student", PutStudent).Methods("POST")
	
	log.Fatal(http.ListenAndServe(":8080", router))
}

func GetCompanyList(w http.ResponseWriter, r *http.Request){
	companyList = append(companyList, Company{Name: "Test Name"})


	log.Print("Serving companyList")
	json.NewEncoder(w).Encode(companyList)
}

func GetStudent(w http.ResponseWriter, r *http.Request){
	params := r.URL.Query()

	log.Print(params["VT_ID"])

	
}

func PutStudent(w http.ResponseWriter, r *http.Request){
	params := r.URL.Query()

	log.Print(params["VT_ID"])

	// exec.Command(/* shell script */, /* 90-number */)
}
