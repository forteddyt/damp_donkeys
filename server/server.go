package main

import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/damp_donkeys/edidutil"
	"github.com/rs/cors"
)

// Company type information
type Company struct {
	Name string `json:"Name,omitempty"`
}

type Student struct{
	DisplayName string `json:"DisplayName,omitempty"`
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

	c := cors.New(cors.Options{
	    AllowedOrigins: []string{"http://csrcint.cs.vt.edu", "*"},
	    AllowCredentials: true,
	    // Enable Debugging for testing, consider disabling in production
	    Debug: true,
	})
	
	handler := cors.Default().Handler(router)
	handler = c.Handler(handler)
	http.ListenAndServe(":8080", handler)
}

func GetCompanyList(w http.ResponseWriter, r *http.Request){
	companyList = append(companyList, Company{Name: "Test Name"})


	log.Print("Serving companyList")
	json.NewEncoder(w).Encode(companyList)
}

func GetStudent(w http.ResponseWriter, r *http.Request){
	params := r.URL.Query()
	
	log.Printf("Origin Header: %s", r.Header.Get("Origin"))
	log.Printf("get_student api called with [%s]\n", params)
	studentInfo := edidutil.ObtainEdidInfo(params["VT_ID"][0]) // {} on failure

	// Client request error
	if(len(studentInfo) == 0){
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	
	json.NewEncoder(w).Encode(studentInfo)
}

func PutStudent(w http.ResponseWriter, r *http.Request){
	params := r.URL.Query()

	log.Print(params["VT_ID"])
}

