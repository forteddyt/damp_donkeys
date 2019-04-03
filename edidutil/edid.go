package edidutil

import (
    "log"
    "os" // To get GOPATH environment variable
    "os/exec" // For running perl script to get info from 90 number
    "strings"
)

// Returns populated string: string mapping on success, empty string: string mapping on fail
func ObtainEdidInfo(uid string) map[string] string{
    gopath := os.Getenv("GOPATH")
    if len(gopath) == 0 {
	log.Printf("GOPATH environment variable not set; Cannot service request, returning empty mapping.")
	return map[string]string{}
    }
    edidpath := gopath + "/src/github.com/damp_donkeys/edidutil/edid.pl"

    out, err := exec.Command("/usr/bin/perl", edidpath, uid).Output()
    
    if err != nil {
	log.Printf("Requested edid.pl, got [%s]\n", err)
	return map[string]string{}
    } else {
	student := map[string]string{"dispName": "", "major": "", "class": ""}
	s := strings.Split(string(out), ";") //convert []byte to a string and split on ';'
	
	student["dispName"] = s[0]
	student["major"] = s[1]
	student["class"] = s[2]
	
	log.Printf("Requested edid.pl, got [%s]\n", out)
	return student
    }
}
