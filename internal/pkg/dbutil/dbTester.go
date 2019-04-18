package dbutil

// import (
// 	"fmt"
// )

//func main() {
    //db, err := OpenDB("STUFF", "stuff", "StUfF")
    //if err != nil {
    //  fmt.Println("Failed to Open")
    //}
    //defer CloseDB(db)
    //Testing
    //Admin
    //isStarted, err := StartCareerFair(db, "Overall Test 1", "NO COMMENTS")
    //if !isStarted || err != nil {
    //  fmt.Println(err)
    //}
    //isAddedEmployer, err := AddEmployer(db, "OverallEmployer", "7AB7FEA73C1F4768E8DF6B7C944C50E9607A4A1B")
    //if !isAddedEmployer || err != nil {
    //  fmt.Println(err)
    //}
    //isUpdated, err := UpdatePassword(db, "OverallEmployer", "5FB8319E3C71287EA6DF060B6B597BED555643F9")
    //if !isUpdated || err != nil {
    //  fmt.Println(err)
    //}
    //employer, err := CheckPasswordHash(db, "5FB8319E3C71287EA6DF060B6B597BED555643F9")    
    //if err != nil {
    //  fmt.Println(err)
    //}
    //fmt.Println(employer)
    //fairs, err := ShowCareerFairsByName(db)
    //if err != nil {
    //  fmt.Println(err)
    //}
    //students, err := GetNumberOfStudents(db, "TestMethods")
    //if err != nil {
    //  fmt.Println(err)
    //}
    //fmt.Println(students)
    //interviews, err := GetNumberOfInterviews(db, "Overall Test 1")
        //if err != nil {
        //        fmt.Println(err)
        //}
        //fmt.Println(interviews)
    //fmt.Println(fairs)
    //Student
    //employers, err := ShowEmployersToStudents(db) 
    //if err != nil {
    //  fmt.Println(err)
    //}
    //fmt.Println(employers)
    //isAdded, err := AddInterview(db, "573CCF70469A7FDFDFF2AF27BAA552A25EA1793F", "OverallEmployer")
    //if !isAdded || err != nil {
    //  fmt.Println(err)
    //}
    //Employer
    //employers, err := ShowEmployersInterviewing(db, "Overall Test 1") 
    //if err != nil {
    //  fmt.Println(err)
    //}
    //fmt.Println(employers)
    //isUpdated, err := UpdateEmployerName(db, "IBM", "USA")
    //if !isUpdated || err != nil {
    //  fmt.Println(err)
    //}
//}