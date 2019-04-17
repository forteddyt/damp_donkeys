package main

import (
	"os"
	"log"
	"errors"

	"net/http"
	"math/rand"
	"io/ioutil"
	"crypto/sha1"
	"encoding/json"
	"encoding/base64"
	
	"github.com/rs/cors"
	"github.com/gorilla/mux"
	"github.com/damp_donkeys/edidutil"
	"github.com/damp_donkeys/jwtutil"
	"github.com/damp_donkeys/dbutil"
)

const(
	JWT_DURATION = 30 // In minutes. How long the token will stay valid between requests
	USER_CODE_LENGTH = 5 // How long user / company unique codes should be
	USER_CODE_CHARS = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" // The characters allowed in the unique codes
)

// Credentials obtained from .gitignored files on server startup
var DBPassword, DBUsername string

type CheckInResp struct {
	CompanyName string `json:"company_name"`
	Students []dbutil.Interview `json:"students"`
	JWT string `json:"jwt"`
}

type AddCompanyResp struct {
	CompanyName string `json:"company_name`
	UserCode string `json:"user_code"`
	JWT string `json:"jwt"`
}

type PutResetCodeResp AddCompanyResp // Same as AddCompanyResp, just different name

type CareerFairListResp struct {
	CareerFairList []string `json:"career_fair_list"`
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
	router.HandleFunc("/career_fair_list", GetCareerFairList).Methods("GET")

	router.HandleFunc("/interview_check_in", InterviewCheckIn).Methods("PUT")
	router.HandleFunc("/add_company", AddCompany).Methods("PUT")
	router.HandleFunc("/reset_code", PutResetCode).Methods("PUT")

	router.HandleFunc("/delete_company", DeleteCompany).Methods("DELETE")

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

func hashHelper(str string) string {
	bv := []byte(str)
	hasher := sha1.New()
	hasher.Write(bv)
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

func genCode(codeLength int) string {
	// 26 lowercase + 26 uppercase + 10 digit possibilities
	bv := make([]byte, codeLength)

	for i := range bv {
		bv[i] = USER_CODE_CHARS[rand.Int63() % int64(len(USER_CODE_CHARS))]
	}

	return string(bv)
}

func PutResetCode(w http.ResponseWriter, r *http.Request){
	params := r.URL.Query()

	// -> ERROR HANDLING
	if len(params["company_name"]) == 0 || params["company_name"][0] == "" ||
		len(params["jwt"]) == 0 || params["jwt"][0] == "" {
		log.Printf("Missing paramaters\n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// -> JWT ERROR HANDLING
	old_jwt := params["jwt"][0]
	is_valid, err := jwtutil.IsValidToken(old_jwt)
	if !is_valid {
		log.Printf("JWT Token invalid: %s\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	claims, err := jwtutil.ParseClaims(old_jwt)

	// Something went wrong internally
	if err != nil {
		log.Printf("ParseClaims error: \n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jwt_user := claims.User
	// Only admins should be able to add a company
	if jwt_user != "admin" {
		// Note: Hard coded "admin" could (/should) eventually be replaced with a cross check to some 'Admins' Table in the db 
		log.Printf("JWT invalid for requested user [%s != %s]\n", jwt_user, "admin")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	new_jwt, err := jwtutil.RefreshToken(old_jwt, JWT_DURATION)

	if err != nil {
		log.Printf("RefreshToken error: \n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// <- END JWT ERROR HANDLING

	dbconn, err := dbutil.OpenDB("dev", DBUsername, DBPassword)
	if err != nil {
		log.Printf("Database connection failed: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer dbutil.CloseDB(dbconn)

	// Try to generate a unique code for the company. Try up to 10 times
	updated := false
	err = nil
	var userCode string
	var attempts int
	for attempts = 0; attempts < 10 && err == nil && !updated; attempts++ {
		userCode = genCode(USER_CODE_LENGTH)
		hashedUserCode := hashHelper(userCode)

		updated, err = dbutil.UpdatePassword(dbconn, params["company_name"][0], hashedUserCode)
	}
	if attempts >= 10 {
		log.Printf("Could not add company to database: Exceeded code generation attempts\n")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err != nil {
		log.Printf("Could not add company to database: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// <- END ERROR HANDLING

	resp := &PutResetCodeResp {
		CompanyName: params["company_name"][0],
		UserCode: userCode,
		JWT: new_jwt,
	}

	// Comment out this print after testing
	log.Printf("Generated code \"%s\" for company \"%s\"", resp.UserCode, resp.CompanyName)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func DeleteCompany(w http.ResponseWriter, r *http.Request){
	params := r.URL.Query()

	// -> ERROR HANDLING
	if len(params["company_name"]) == 0 || params["company_name"][0] == "" ||
		len(params["jwt"]) == 0 || params["jwt"][0] == "" {
		log.Printf("Missing paramaters\n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// -> JWT ERROR HANDLING
	old_jwt := params["jwt"][0]
	is_valid, err := jwtutil.IsValidToken(old_jwt)
	if !is_valid {
		log.Printf("JWT Token invalid: %s\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	claims, err := jwtutil.ParseClaims(old_jwt)

	// Something went wrong internally
	if err != nil {
		log.Printf("ParseClaims error: \n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jwt_user := claims.User
	// Only admins should be able to add a company
	if jwt_user != "admin" {
		// Note: Hard coded "admin" could (/should) eventually be replaced with a cross check to some 'Admins' Table in the db 
		log.Printf("JWT invalid for requested user [%s != %s]\n", jwt_user, "admin")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	new_jwt, err := jwtutil.RefreshToken(old_jwt, JWT_DURATION)

	if err != nil {
		log.Printf("RefreshToken error: \n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// <- END JWT ERROR HANDLING

	dbconn, err := dbutil.OpenDB("dev", DBUsername, DBPassword)
	if err != nil {
		log.Printf("Database connection failed: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer dbutil.CloseDB(dbconn)

	log.Printf("%s\n", new_jwt)

	// company_removed, err := dbutil.RemoveEmployer(dbconn, params["company_name"][0])
	// <- ERROR HANDLING

	// Remainder will be written once RemoveEmployer is implemented
}

func GetCareerFairList(w http.ResponseWriter, r *http.Request){
	params := r.URL.Query()

	// -> ERROR HANDLING
	log.Printf("interview_check_in api called with [%s]\n", params)
	if len(params["jwt"]) == 0 || params["jwt"][0] == "" {
		log.Printf("Missing paramaters\n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// -> JWT ERROR HANDLING
	old_jwt := params["jwt"][0]
	is_valid, err := jwtutil.IsValidToken(old_jwt)
	if !is_valid {
		log.Printf("JWT Token invalid: %s\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	claims, err := jwtutil.ParseClaims(old_jwt)

	// Something went wrong internally
	if err != nil {
		log.Printf("ParseClaims error: \n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jwt_user := claims.User
	// Only admins should be able to add a company
	if jwt_user != "admin" {
		// Note: Hard coded "admin" could (/should) eventually be replaced with a cross check to some 'Admins' Table in the db 
		log.Printf("JWT invalid for requested user [%s != %s]\n", jwt_user, "admin")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	new_jwt, err := jwtutil.RefreshToken(old_jwt, JWT_DURATION)

	if err != nil {
		log.Printf("RefreshToken error: \n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// <- END JWT ERROR HANDLING

	dbconn, err := dbutil.OpenDB("dev", DBUsername, DBPassword)
	if err != nil {
		log.Printf("Database connection failed: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer dbutil.CloseDB(dbconn)

	career_fair_list, err := dbutil.ShowCareerFairsByName(dbconn)

	if err != nil {
		log.Printf("Could not get career fairs from database: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := &CareerFairListResp {
		CareerFairList: career_fair_list,
		JWT: new_jwt,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func AddCompany(w http.ResponseWriter, r *http.Request){
	params := r.URL.Query()

	// -> ERROR HANDLING
	log.Printf("interview_check_in api called with [%s]\n", params)
	if len(params["company_name"]) == 0  || params["company_name"][0] == "" ||
		len(params["jwt"]) == 0 || params["jwt"][0] == "" {
		log.Printf("Missing paramaters\n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// -> JWT ERROR HANDLING
	old_jwt := params["jwt"][0]
	is_valid, err := jwtutil.IsValidToken(old_jwt)
	if !is_valid {
		log.Printf("JWT Token invalid: %s\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	claims, err := jwtutil.ParseClaims(old_jwt)

	// Something went wrong internally
	if err != nil {
		log.Printf("ParseClaims error: \n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jwt_user := claims.User
	// Only admins should be able to add a company
	if jwt_user != "admin" {
		// Note: Hard coded "admin" could (/should) eventually be replaced with a cross check to some 'Admins' Table in the db 
		log.Printf("JWT invalid for requested user [%s != %s]\n", jwt_user, "admin")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	new_jwt, err := jwtutil.RefreshToken(old_jwt, JWT_DURATION)

	if err != nil {
		log.Printf("RefreshToken error: \n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// <- END JWT ERROR HANDLING

	dbconn, err := dbutil.OpenDB("dev", DBUsername, DBPassword)
	if err != nil {
		log.Printf("Database connection failed: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer dbutil.CloseDB(dbconn)

	// Try to generate a unique code for the company. Try up to 10 times
	addedCompany := false
	err = nil
	var userCode string
	var attempts int
	for attempts = 0; attempts < 10 && err == nil && !addedCompany; attempts++ {
		userCode = genCode(USER_CODE_LENGTH)
		hashedUserCode := hashHelper(userCode)

		addedCompany, err = dbutil.AddEmployer(dbconn, params["company_name"][0], hashedUserCode)
	}
	if attempts >= 10 {
		log.Printf("Could not add company to database: Exceeded code generation attempts\n")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err != nil {
		log.Printf("Could not add company to database: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// <- END ERROR HANDLING

	resp := &AddCompanyResp {
		CompanyName: params["company_name"][0],
		UserCode: userCode,
		JWT: new_jwt,
	}

	// Comment out this print after testing
	log.Printf("Generated code \"%s\" for company \"%s\"", resp.UserCode, resp.CompanyName)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func InterviewCheckIn(w http.ResponseWriter, r *http.Request){
	params := r.URL.Query()

	// ERROR HANDLING

	log.Printf("interview_check_in api called with [%s]\n", params)
	if len(params["company_name"]) ==0  || params["company_name"][0] == "" ||
	   len(params["display_name"]) == 0 || params["display_name"][0] == "" ||
	   len(params["major"]) == 0 /*|| params["major"][0] == ""*/ ||
	   len(params["class"]) == 0 /*|| params["class"][0] == ""*/ ||
	   len(params["VT_ID"]) == 0 || params["VT_ID"][0] == "" {
		log.Printf("Missing paramaters\n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// First thing's first, hash the VT_ID (90- number)
	hashedId := hashHelper(params["VT_ID"][0])

	dbconn, err := dbutil.OpenDB("dev", DBUsername, DBPassword)
	if err != nil {
		log.Printf("Database connection failed: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer dbutil.CloseDB(dbconn)

	// END ERROR HANDLING

	addedStudent, err := dbutil.AddStudent(dbconn, params["display_name"][0], params["major"][0], params["class"][0], hashedId)
	if err != nil {
		log.Printf("Could not add student to database: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if addedStudent {
		log.Printf("New student added to database\n")
	} else  {
		log.Printf("Student already exists in database\n")
	}

	addedInterview, err := dbutil.AddInterview(dbconn, hashedId, params["company_name"][0])
	if err != nil {
		log.Printf("Could not add interview to database: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if addedInterview {
		log.Printf("New interview added to database\n")
	} else  {
		log.Printf("Interview already exists in database (THIS SHOULD NOT HAPPEN)\n")
	}

	w.WriteHeader(http.StatusOK)
}

func CompanyCheckIns(w http.ResponseWriter, r *http.Request){
	params := r.URL.Query()
	
	log.Printf("company_check_ins api called with [%s]\n", params)
	if len(params["jwt"]) == 0 || params["jwt"][0] == "" ||
		len(params["company_name"]) == 0 || params["company_name"][0] == "" {
		w.WriteHeader(http.StatusBadRequest)
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
		return
	}

	claims, err := jwtutil.ParseClaims(old_jwt)

	// Something went wrong internally
	if err != nil {
		log.Printf("ParseClaims error: \n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jwt_user := claims.User
	// If JWT and requested company are invalid
	if jwt_user != company_name {
		log.Printf("JWT invalid for requested company [%s != %s]\n", jwt_user, company_name)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	new_jwt, err := jwtutil.RefreshToken(old_jwt, JWT_DURATION)

	if err != nil {
		log.Printf("RefreshToken error: \n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// END ERROR HANDLING

	dbconn, err := dbutil.OpenDB("dev", DBUsername, DBPassword)
	if err != nil {
		log.Printf("Database connection failed: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer dbutil.CloseDB(dbconn)
	students, err := dbutil.ShowStudents(dbconn, company_name)

	// Database request error
	if err != nil {
		log.Printf("Could not get checked in students from database: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
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
	} else {
		givenPwdHash := params["password_hash"][0]
		dbconn, err := dbutil.OpenDB("dev", DBUsername, DBPassword)
		if err != nil {
			log.Printf("Database connection failed: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer dbutil.CloseDB(dbconn)
		user, err := dbutil.CheckPasswordHash(dbconn, givenPwdHash) // dbutil.getUser(givenPwdHash)

		// Database request error
		if err != nil {
			log.Printf("Could not get user from database: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// No user with that password hash was found
		if user == "" {
			log.Printf("No user with password hash '%s' found\n", givenPwdHash)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		jwt, err := jwtutil.CreateToken(user, JWT_DURATION)
		if err != nil {
			log.Printf("JWT Token creation failed: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		m := map[string]string{"jwt": jwt}
		log.Printf("JWT Token created successfully, valid for %d minutes\n", JWT_DURATION)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(m)
	}
}

