# The Damp Donkeys


## Setup Instructions
Install Golang version 1.12.

This server uses `mux` to route requests, install the package by executing `go get github.com/gorilla/mux`


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

