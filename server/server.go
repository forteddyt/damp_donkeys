package main

import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/damp_donkeys/edidutil"
	"github.com/damp_donkeys/jwtutil"
	"github.com/damp_donkeys/dbutil"
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
	router.HandleFunc("/login", Login).Methods("GET")

	c := cors.New(cors.Options{
	    AllowedOrigins: []string{"https://csrcint.cs.vt.edu"},
	    AllowCredentials: true,
	    // Enable Debugging for testing, consider disabling in production
	    Debug: true,
	})
	
	handler := cors.Default().Handler(router)
	// handler = c.Handler(handler)
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

func Login(w http.ResponseWriter, r *http.Request){
	log.Printf("login api called with [%s]\n", params)

	params := r.URL.Query()
	givenPwdHash := params["password_hash"]

	// Client request error
	if(givenPwdHash == nil){
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode()
	} else {
		dbconn, err := dbutil.OpenDB("")
		user, err = "", nil // dbutil.getUser(givenPwdHash)

		// Database request error
		if err != nil {
			log.Printf("Could not get password hash from database\n")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode()
			return
		}

		// No user with that password hash was found
		if user == "" {
			log.Printf("No user with password hash '%s' found\n", givenPwdHash)
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode()
		}

	}
}

