package edidutil

import (
	"log"
	"strings"
	"os/exec" // For running perl script to get info from 90 number
	"github.com/damp_donkeys/configs/pathing"
)

// Returns populated string: string mapping on success, empty string: string mapping on fail
func ObtainEdidInfo(uid string) ([]string, error){
	goPath, err := pathing.GoPath()
	if err != nil {
		log.Printf("Cannot service request, returning empty mapping: %s\n", err)
		return nil, err
	}
	edidScript := goPath + pathing.ProjectEDIDScript
	edidPem := goPath + pathing.ProjectEDIDPemFile
	edidCrt := goPath + pathing.ProjectEDIDCrtFile
	edidKey := goPath + pathing.ProjectEDIDKeyFile

	out, err := exec.Command("/usr/bin/perl", edidScript, edidPem, edidCrt, edidKey, uid).Output()

	var s []string
	if err != nil {
		log.Printf("edid.pl failed: %s\n", err)
		return s, err
	}

	s = strings.Split(string(out), ";") //convert []byte to a string and split on ';'
	log.Printf("Requested edid.pl, got [%s]\n", out)
	return s, nil
}
