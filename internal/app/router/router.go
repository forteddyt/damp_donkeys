package router

import (
	"net/http"

	"github.com/rs/cors"
	"github.com/gorilla/mux"
)

func Setup() {
	router := mux.NewRouter()
	
	router.HandleFunc("/company_list", GetCompanyList).Methods("GET")
	router.HandleFunc("/company_interviewed_list", GetCompanyInterviewedList).Methods("GET")
	router.HandleFunc("/get_student", GetStudent).Methods("GET")
	router.HandleFunc("/company_check_ins", GetCompanyCheckIns).Methods("GET")
	router.HandleFunc("/login", GetLogin).Methods("GET")
	router.HandleFunc("/career_fair_list", GetCareerFairList).Methods("GET")
	router.HandleFunc("/career_fair_stats", GetCareerFairStats).Methods("GET")

	router.HandleFunc("/interview_check_in", PutInterviewCheckIn).Methods("PUT")
	router.HandleFunc("/add_company", PutCompany).Methods("PUT")
	router.HandleFunc("/reset_code", PutResetCode).Methods("PUT")
	router.HandleFunc("/career_fair", PutCareerFair).Methods("PUT")

	router.HandleFunc("/delete_company", DeleteCompany).Methods("DELETE")

	c := cors.New(cors.Options{
	    AllowedOrigins: []string{"*"},
	    AllowCredentials: true,
	    AllowedMethods: []string{"GET", "DELETE", "PUT", "POST"},
	    // Enable Debugging for testing, consider disabling in production
	    // Debug: true,
	})
	
	handler := cors.Default().Handler(router)
	handler = c.Handler(handler)
	http.ListenAndServe(":8080", handler)
}