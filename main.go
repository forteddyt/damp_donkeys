package main

import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/gorilla/mux"
)

type Company struct {
	Name string `json:"Name,omitempty"`
}

var companyList []Company

func main() {
	log.Print("Starting server...")

	router := mux.NewRouter()
	
	router.HandleFunc("/company_list", GetCompanyList).Methods("GET")
	
	log.Fatal(http.ListenAndServe(":8080", router))
}

func GetCompanyList(w http.ResponseWriter, r *http.Request){
	companyList = append(companyList, Company{Name: "Test Name"})


	log.Print("Serving companyList")
	json.NewEncoder(w).Encode(companyList)
}
