package edidutil

import (
    "fmt"
    "log"
    "os/exec" // For running perl script to get info from 90 number
)

func ObtainEdidInfo(uid string) {
    out, err := exec.Command("/bin/sh", "test.sh", uid).Output()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("The date is %s\n", out)
}
