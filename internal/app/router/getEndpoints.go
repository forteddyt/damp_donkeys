package router

import(
	"log"

	"net/http"
	"encoding/json"
	
	"github.com/damp_donkeys/internal/pkg/edidutil"
	"github.com/damp_donkeys/internal/pkg/jwtutil"
	"github.com/damp_donkeys/internal/pkg/dbutil"
	"github.com/damp_donkeys/internal/pkg/confidante"

	"github.com/damp_donkeys/configs/resp"
)

func getCompanyListForCareerFair(w http.ResponseWriter, r *http.Request){
	params := r.URL.Query()

	if len(params) != 2 || len(params["jwt"]) == 0 || params["jwt"][0] == "" ||
		len(params["career_fair_name"]) == 0 || params["career_fair_name"][0] == ""{
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
	companies, err := dbutil.ShowAllEmployersByCareerFair(dbconn, params["career_fair_name"][0])

	// Database request error
	if err != nil {
		log.Printf("Could not get interviewing companies from database: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := &resp.GetCompanyListForCareerFair {
		CompanyList: companies,
		JWT: new_jwt,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func GetCompanyList(w http.ResponseWriter, r *http.Request){
	params := r.URL.Query()

	// -> ERROR HANDLING
	log.Printf("company_list api called with [%s]\n", params)
	if len(params) > 0 {
		getCompanyListForCareerFair(w, r)
		return
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
	companies, err := dbutil.ShowEmployersToStudents(dbconn)

	// Database request error
	if err != nil {
		log.Printf("Could not get interviewing companies from database: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := &resp.GetCompanyList {
		CompanyList: companies,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func GetStudent(w http.ResponseWriter, r *http.Request){
	params := r.URL.Query()
	

	// -> ERROR HANDLING
	log.Printf("get_student api called with [%s]\n", params)
	if len(params["VT_ID"]) == 0 || params["VT_ID"][0] == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	student, err := edidutil.ObtainEdidInfo(params["VT_ID"][0]) // {} on failure

	// Client request error
	if err != nil {
		if student == nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		return
	}
	// <- ERROR HANDLING
	
	resp := &resp.GetStudent {
		DisplayName: student[0],
		Major: student[1],
		Class: student[2],
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func GetCompanyCheckIns(w http.ResponseWriter, r *http.Request){
	params := r.URL.Query()
	
	// -> ERROR HANDLING
	log.Printf("company_check_ins api called with [%s]\n", params)
	if len(params["jwt"]) == 0 || params["jwt"][0] == "" ||
		len(params["company_name"]) == 0 || params["company_name"][0] == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	company_name := params["company_name"][0]

	statusCode, new_jwt := tokenRefreshHelper(params["jwt"][0], company_name, JWTDuration)
	if statusCode != 0 {
		w.WriteHeader(statusCode)
		return
	}
	// <- ERROR HANDLING

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
	students, err := dbutil.ShowStudents(dbconn, company_name)

	// Database request error
	if err != nil {
		log.Printf("Could not get checked in students from database: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := &resp.GetCompanyCheckIn {
		CompanyName: company_name,
		Students: students,
		JWT: new_jwt,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func GetLogin(w http.ResponseWriter, r *http.Request){
	params := r.URL.Query()
	log.Printf("login api called with [%s]\n", params)

	// Client request error
	if len(params["code"]) == 0 || params["code"][0] == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	code := params["code"][0]
	hashedCode := hashHelper(code)

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
	user, err := dbutil.CheckPasswordHash(dbconn, hashedCode)

	// Database request error
	if err != nil {
		log.Printf("Could not get user from database: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// No user with that password hash was found
	if user == "" {
		log.Printf("No user with code hash %x found\n", hashedCode)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	jwt, err := jwtutil.CreateToken(user, JWTDuration)
	if err != nil {
		log.Printf("JWT creation failed: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := &resp.GetLogin {
		JWT: jwt,
	}
	log.Printf("JWT created successfully\n")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
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

	career_fair_list, err := dbutil.ShowCareerFairsByName(dbconn)

	if err != nil {
		log.Printf("Could not get career fairs from database: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := &resp.GetCareerFairList {
		CareerFairList: career_fair_list,
		JWT: new_jwt,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

