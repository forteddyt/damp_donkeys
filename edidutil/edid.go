package edidutil

import (
    "log"
    "os/exec" // For running perl script to get info from 90 number
    "strings"
)

func ObtainEdidInfo(uid string) map[string] string{
    student := map[string]string{"dispName": "", "major": "", "class": ""}
    out, err := exec.Command("/usr/bin/perl", "edidutil/edid.pl", uid).Output()
    if err != nil {
        log.Printf("ERROR: [%s]\n", err)
	return student
    } else {
        log.Printf("Requested edid.pl, got [%s]\n", out)
	s := strings.Split(string(out), ";") //convert []byte to a string and split on ';'
	
	student["dispName"] = s[0]
	student["major"] = s[1]
	student["class"] = s[2]
	
	return student
    }
}
