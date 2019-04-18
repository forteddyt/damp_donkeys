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

	// company_removed, err := dbutil.RemoveEmployer(dbconn, params["company_name"][0])
	// <- ERROR HANDLING

	// Remainder will be written once RemoveEmployer is implemented
}