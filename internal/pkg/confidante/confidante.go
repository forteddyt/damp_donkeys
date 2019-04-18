package confidante

import(
	"os"
	"log"
    "io/ioutil"
    "encoding/json"
	"github.com/damp_donkeys/configs/pathing"
)

func JWTSingingKey() (string, error){
	goPath, err := pathing.GoPath()

	if err != nil {
	    return "", err
	}

	jsonFile, err := os.Open(goPath + pathing.ProjectJWTFile)
	
	if err != nil { return "", err }
	defer jsonFile.Close() // Close json file at end

	byteValue , _ := ioutil.ReadAll(jsonFile) // Read json as []byte

	var result map[string]string // result will be a mapping of string to []byte

	json.Unmarshal(byteValue, &result)

	signingKey := result["signingKey"]
	return signingKey, nil
}

// Returns DBUsername, DBPassword, error
func DBCredentials() (string, string, error) {
	goPath, err := pathing.GoPath()

	if err != nil {
	    return "", "", err
	}

	jsonFile, err := os.Open(goPath + pathing.ProjectDBFile)

	if err != nil { return "", "", err }
	defer jsonFile.Close() // Close json file at end

	byteValue , _ := ioutil.ReadAll(jsonFile) // Read json as []byte
	var result map[string]string // result will be a mapping of string to a mapping of string to string

	json.Unmarshal(byteValue, &result)

	log.Println("Getting database username...")
	DBUsername := result["username"]
	DBPassword := result["password"]

	return DBUsername, DBPassword, nil
}