package router

import(
	"log"

	"net/http"
	
	"github.com/damp_donkeys/internal/pkg/dbutil"
	"github.com/damp_donkeys/internal/pkg/confidante"
)

func DeleteCompany(w http.ResponseWriter, r *http.Request){
	params := r.URL.Query()

	// -> ERROR HANDLING
	if len(params["company_name"]) == 0 || params["company_name"][0] == "" ||
	   len(params["career_fair_name"]) == 0 || params["career_fair_name"][0] == "" ||
		len(params["jwt"]) == 0 || params["jwt"][0] == "" {
		log.Printf("Missing paramaters\n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	statusCode, new_jwt := tokenRefreshHelper(params["jwt"][0], "admin", JWTDuration)
	if statusCode != 0 {
		w.WriteHeader(statusCode)
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

	log.Printf("%s\n", new_jwt)
	
	interviews, err := dbutil.ShowStudents(dbconn, params["company_name"][0])
	if err != nil {
		log.Printf("Could not get students for given company: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// If admin tries to delete a company with interviews
	if len(interviews) > 0 {
		log.Printf("Cannot remove company that has held a interviews\n")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	deleted, err := dbutil.DeleteEmployer(dbconn, params["company_name"][0], params["career_fair_name"][0])
	if err != nil {
		log.Printf("Could not delete company from career fair: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// <- ERROR HANDLING
	
	if !deleted {
		log.Printf("Did not delete company from career fair\n")
	} else {
		log.Printf("Deleted company from career fair\n")
	}
	w.WriteHeader(http.StatusOK)
}