package pathing
// Stores all pathing information

import(
	"os"
	"errors"
)

// Obtained at runtime as the GoPath is a system-specific enviornment variable
var goPath = ""

const (
	//------ RELATIVE TO GOPATH ------//
	ProjectRepository 	= "src/github.com/damp_donkeys/"
	ProjectPackages 	= ProjectRepository + InternalPackages
	ProjectApps			= ProjectRepository + InternalApps
	
	ProjectPathings		= ProjectRepository + ConfigPathings
	ProjectResps		= ProjectRepository + ConfigResps
	ProjectSecrets		= ProjectRepository + ConfigSecrets

	ProjectMain			= ProjectRepository + Main

	ProjectJWTFile		= ProjectSecrets + JWTFile

	ProjectDBFile		= ProjectSecrets + DBFile

	ProjectEDIDScript	= ProjectPackages + EDIDScript
	ProjectEDIDPemFile	= ProjectSecrets + EDIDPemFile
	ProjectEDIDKeyFile	= ProjectSecrets + EDIDKeyFile
	ProjectEDIDCrtFile	= ProjectSecrets + EDIDCrtFile
	//--------------------------------//

	//------ RELATIVE TO PROJECT REPOSITORY ------//
	InternalPackages	= "internal/pkg/"
	InternalApps		= "internal/app/"

	ConfigPathings		= "configs/pathing/"
	ConfigResps			= "configs/resp/"
	ConfigSecrets		= "configs/secret/"

	Main				= "cmd/sever/"
	//--------------------------------------------//

	//------ RELATIVE TO SECRETS FOLDER ------//
	JWTFile				= "JWTKey.json" // All secret content regarding JWT should be in here
	DBFile				= "DBCredentials.json" // All secret content regarding DB should be in here
	EDIDPemFile			= "ed_id_ca.pem"
	EDIDKeyFile			= "private.key"
	EDIDCrtFile			= "public.crt"
	//----------------------------------------//

	//------ RELATIVE TO INTERNAL PKG FOLDER ------//
	EDIDScript			= "edidutil/edid.pl"
	//---------------------------------------------//
)

func GoPath() (string, error) {
	if goPath == ""{
	    goPath = os.Getenv("GOPATH")
	    if len(goPath) == 0 {
	    	return "", errors.New("GOPATH environment variable not set")
	    }

	    // If last character in path isn't "/", append it
	    if string(goPath[len(goPath)-1:]) != "/" {
	    	goPath = goPath + "/"
	    }
	}

	return goPath, nil
}