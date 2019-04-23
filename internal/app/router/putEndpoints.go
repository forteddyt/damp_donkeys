package router

import(
	"log"

	"net/http"
	"encoding/json"

	"github.com/damp_donkeys/internal/pkg/dbutil"
	"github.com/damp_donkeys/internal/pkg/confidante"

	"github.com/damp_donkeys/configs/resp"
)

func PutCareerFair(w http.ResponseWriter, r *http.Request){
	params := r.URL.Query()

	// -> ERROR HANDLING
	log.Printf("interview_check_in api called with [%s]\n", params)
	if len(params["career_fair_name"]) ==0  || params["career_fair_name"][0] == "" ||
	   len(params["comments"]) ==0  || params["comments"][0] == "" ||
	   len(params["jwt"]) == 0 || params["jwt"][0] == ""{
		log.Printf("Missing paramaters\n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	new_jwt := ""
	if DBName == "dev" { // Give free rein in development database
		new_jwt = params["jwt"][0]
	} else {
		statusCode, n_jwt := tokenRefreshHelper(params["jwt"][0], "admin", JWTDuration)
		new_jwt = n_jwt
		if statusCode != 0 {
			w.WriteHeader(statusCode)
			return
		}
	}

	DBUsername, DBPassword, err := confidante.DBCredentials()
	if err != nil {
		log.Printf("Could not obtain Database credentials: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	dbconn, err := dbutil.OpenDB(DBName, DBUsername, DBPassword)
	if err != nil {
		log.Printf("Database connection failed: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer dbutil.CloseDB(dbconn)

	_, err = dbutil.StartCareerFair(dbconn, params["career_fair_name"][0], params["comments"][0])
	if err != nil {
		log.Printf("Could not create new career fair: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// <- END ERROR HANDLING

	resp := &resp.PutCareerFair {
		CareerFairName: params["career_fair_name"][0],
		Comments: params["comments"][0],
		JWT: new_jwt,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func PutInterviewCheckIn(w http.ResponseWriter, r *http.Request){
	params := r.URL.Query()

	// -> ERROR HANDLING
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

	DBUsername, DBPassword, err := confidante.DBCredentials()
	if err != nil {
		log.Printf("Could not obtain Database credentials: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	dbconn, err := dbutil.OpenDB(DBName, DBUsername, DBPassword)
	if err != nil {
		log.Printf("Database connection failed: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer dbutil.CloseDB(dbconn)
	// <- ERROR HANDLING

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


func PutCompany(w http.ResponseWriter, r *http.Request){
	params := r.URL.Query()

	// -> ERROR HANDLING
	log.Printf("interview_check_in api called with [%s]\n", params)
	if len(params["company_name"]) == 0  || params["company_name"][0] == "" ||
		len(params["jwt"]) == 0 || params["jwt"][0] == "" {
		log.Printf("Missing paramaters\n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	new_jwt := ""
	if DBName == "dev" { // Give free rein in development database
		new_jwt = params["jwt"][0]
	} else {
		statusCode, n_jwt := tokenRefreshHelper(params["jwt"][0], "admin", JWTDuration)
		new_jwt = n_jwt
		if statusCode != 0 {
			w.WriteHeader(statusCode)
			return
		}
	}

	DBUsername, DBPassword, err := confidante.DBCredentials()
	if err != nil {
		log.Printf("Could not obtain Database credentials: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	dbconn, err := dbutil.OpenDB(DBName, DBUsername, DBPassword)
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
	for attempts = 0; attempts < 10 && !addedCompany; attempts++ {
		userCode = genCodeHelper(UserCodeLength)
		hashedUserCode := hashHelper(userCode)
		log.Printf("hashedUserCode: %x\n", hashedUserCode)

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

	resp := &resp.PutCompany {
		CompanyName: params["company_name"][0],
		UserCode: userCode,
		JWT: new_jwt,
	}

	// Comment out this print after testing
	log.Printf("Generated code \"%s\" for company \"%s\"", resp.UserCode, resp.CompanyName)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func PutResetCode(w http.ResponseWriter, r *http.Request){
	params := r.URL.Query()

	// -> ERROR HANDLING
	log.Printf("reset_code api called with [%s]\n", params)
	if len(params["company_name"]) == 0 || params["company_name"][0] == "" ||
		len(params["jwt"]) == 0 || params["jwt"][0] == "" {
		log.Printf("Missing paramaters\n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	company_name := params["company_name"][0]

	new_jwt := ""
	if DBName == "dev" { // Give free rein in development database
		new_jwt = params["jwt"][0]
	} else {
		statusCode, n_jwt := tokenRefreshHelper(params["jwt"][0], "admin", JWTDuration)
		new_jwt = n_jwt
		if statusCode != 0 {
			w.WriteHeader(statusCode)
			return
		}
	}

	DBUsername, DBPassword, err := confidante.DBCredentials()
	if err != nil {
		log.Printf("Could not obtain Database credentials: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	dbconn, err := dbutil.OpenDB(DBName, DBUsername, DBPassword)
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
		userCode = genCodeHelper(UserCodeLength)
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

	resp := &resp.PutResetCode {
		CompanyName: company_name,
		UserCode: userCode,
		JWT: new_jwt,
	}

	// Comment out this print after testing
	log.Printf("Generated code \"%s\" for company \"%s\"", resp.UserCode, resp.CompanyName)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

