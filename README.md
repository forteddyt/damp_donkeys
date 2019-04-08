# The Damp Donkeys


## Setup Instructions
Install Golang version 1.12.

This server uses `gorilla/mux`, `rs/cors`, and `dgrijalva/jwt-go`. Obtain package by running `go get github.com/<package>`



Directory structure should mirror the following:
```
go/
  bin/
      ... // Executable files
  src/
      github.com/
          damp_donkeys/     // This repository MUST be in the 'github.com' directory
              .git/
              server/
              ...
```

## Notice
RESTful API based off https://www.codementor.io/codehakase/building-a-restful-api-with-golang-a6yivzqdo

