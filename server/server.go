package main

import (
	"os"
	"io/ioutil"
	"log"
	"errors"
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/damp_donkeys/edidutil"
	"github.com/damp_donkeys/jwtutil"
	"github.com/damp_donkeys/dbutil"
)

const(
	JWT_DURATION = 30 // In minutes. How long the token will stay valid between requests
)

// Credentials obtained from .gitignored files on server startup
var DBPassword, DBUsername string

type CheckInResp struct {
	CompanyName string `json:"company_name"`
	Students []dbutil.Interview `json:"students"`
	JWT string `json:"jwt"`
}

func main() {
	log.Print("Starting server...")

	err := credentialSetup()
	if err != nil { log.Fatal("Credential setup failed: %s", err) }

	router := mux.NewRouter()
	
	router.HandleFunc("/company_list", GetCompanyList).Methods("GET")
	router.HandleFunc("/get_student", GetStudent).Methods("GET")
	router.HandleFunc("/company_check_ins", CompanyCheckIns).Methods("GET")
	router.HandleFunc("/login", Login).Methods("GET")

	c := cors.New(cors.Options{
	    AllowedOrigins: []string{"https://csrcint.cs.vt.edu"},
	    AllowCredentials: true,
	    // Enable Debugging for testing, consider disabling in production
	    Debug: true,
	})
	
	handler := cors.Default().Handler(router)
	handler = c.Handler(handler)
	http.ListenAndServe(":8080", handler)
}

func credentialSetup() error{
	gopath := os.Getenv("GOPATH")
	if len(gopath) == 0 {
		return errors.New("GOPATH environment variable not set; Cannot obtain credentials")
	}

	log.Println("Opening server's secret.json...")
	jsonFile, err := os.Open(gopath + "/src/github.com/damp_donkeys/server/secret.json")

	if err != nil { return err }
	defer jsonFile.Close() // Close json file at end
	defer log.Println("Closing server's secret.json") // Then log the close

	byteValue , _ := ioutil.ReadAll(jsonFile) // Read json as []byte
	var result map[string]map[string]string // result will be a mapping of string to a mapping of string to string

	json.Unmarshal(byteValue, &result)

	log.Println("Getting database username...")
	DBUsername = result["database"]["username"]
	log.Println("Database username obtained!")
	log.Println("Getting database password...")
	DBPassword = result["database"]["password"]
	log.Println("Database password obtained!")

	return nil
}

func CompanyCheckIns(w http.ResponseWriter, r *http.Request){
	params := r.URL.Query()
	
	log.Printf("company_check_ins api called with [%s]\n", params)
	if len(params["jwt"]) == 0 || params["jwt"][0] == "" ||
		len(params["company_name"]) == 0 || params["company_name"][0] == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{})
		return
	}

	old_jwt := params["jwt"][0]
	company_name := params["company_name"][0]

	// ERROR HANDLING

	// If jwt has expired, deny access
	is_valid, err := jwtutil.IsValidToken(old_jwt)
	if !is_valid {
		log.Printf("JWT Token invalid: %s\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{})
		return
	}

	claims, err := jwtutil.ParseClaims(old_jwt)

	// Something went wrong internally
	if err != nil {
		log.Printf("ParseClaims error: \n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{})
		return
	}

	jwt_user := claims.User
	// If JWT and requested company are invalid
	if jwt_user != company_name {
		log.Printf("JWT invalid for requested company [%s != %s]\n", jwt_user, company_name)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{})
		return
	}

	new_jwt, err := jwtutil.RefreshToken(old_jwt, JWT_DURATION)

	if err != nil {
		log.Printf("RefreshToken error: \n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{})
		return
	}

	// END ERROR HANDLING

	dbconn, err := dbutil.OpenDB("dev", DBUsername, DBPassword)
	if err != nil {
		log.Printf("Database connection failed: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{})
		return
	}
	defer dbutil.CloseDB(dbconn)
	students, err := dbutil.ShowStudents(dbconn, company_name)

	// Database request error
	if err != nil {
		log.Printf("Could not get checked in students from database: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{})
		return
	}

	resp := &CheckInResp{
		CompanyName: company_name,
		Students: students,
		JWT: new_jwt,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func GetCompanyList(w http.ResponseWriter, r *http.Request){
	log.Print("Serving companyList")
	
	// Temporary until db function is written
	temp := [3]string{"Company 1", "Company 2", "Long named company example"}
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(temp)
}

func GetStudent(w http.ResponseWriter, r *http.Request){
	params := r.URL.Query()
	
	log.Printf("Origin Header: %s\n", r.Header.Get("Origin"))
	log.Printf("get_student api called with [%s]\n", params)
	if len(params["VT_ID"]) == 0 || params["VT_ID"][0] == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{})
		return
	}

	studentInfo := edidutil.ObtainEdidInfo(params["VT_ID"][0]) // {} on failure

	// Client request error
	if len(studentInfo) == 0 {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	
	json.NewEncoder(w).Encode(studentInfo)
}

func Login(w http.ResponseWriter, r *http.Request){
	params := r.URL.Query()
	log.Printf("login api called with [%s]\n", params)

	// Client request error
	if len(params["password_hash"]) == 0 || params["password_hash"][0] == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{})
	} else {
		givenPwdHash := params["password_hash"][0]
		dbconn, err := dbutil.OpenDB("dev", DBUsername, DBPassword)
		if err != nil {
			log.Printf("Database connection failed: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{})
			return
		}
		defer dbutil.CloseDB(dbconn)
		user, err := dbutil.CheckPasswordHash(dbconn, givenPwdHash) // dbutil.getUser(givenPwdHash)

		// Database request error
		if err != nil {
			log.Printf("Could not get user from database: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{})
			return
		}

		// No user with that password hash was found
		if user == "" {
			log.Printf("No user with password hash '%s' found\n", givenPwdHash)
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{})
			return
		}

		jwt, err := jwtutil.CreateToken(user, JWT_DURATION)
		if err != nil {
			log.Printf("JWT Token creation failed: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{})
			return
		}

		m := map[string]string{"jwt": jwt}
		log.Printf("JWT Token created successfully, valid for %d minutes\n", JWT_DURATION)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(m)
	}
}

