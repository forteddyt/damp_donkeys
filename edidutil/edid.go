package edidutil

import (
    "fmt"
    "log"
    "os/exec" // For running perl script to get info from 90 number
)

func ObtainEdidInfo(uid string) {
    out, err := exec.Command("/usr/bin/perl", "~/go/src/server/edidutil/edid.pl", uid).Output()
    if err != nil {
        log.Printf("ERROR: [%s]\n", err)
    } else {
        fmt.Printf("%s\n", out)
    }
}
