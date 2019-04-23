package dbutil

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

//Struct Type For Students
type Student struct {
    DisplayName string `json:"display_name"`
    Major string `json:"major"`
    Class string `json:"class"`
    IdNumber string `json:"idnumber"`
}

type Interview struct {
    Student Student `json:"student"`
    Time string `json:"time"`
}

func OpenDB(databasename string, username string, password string) (*sql.DB, error) {
    //Reminder to Change Dev to CSRCINTERVIEW
    db, err := sql.Open("mysql",
        username + ":" + password + "@tcp(127.0.0.1:3306)/" + databasename)
    if err != nil {
        return nil, err
    }
    return db, nil
}

func CloseDB(db *sql.DB) {
    defer db.Close();
}

func CheckEmployer(db *sql.DB, employername string) (bool, error){
        stmt, err := db.Prepare("SELECT name FROM Employers WHERE name LIKE ?;")
        if err != nil {
                return false, err
        }
        defer stmt.Close()
        rows, err := stmt.Query(employername)
        if err != nil {
                return false, err
        }
        defer rows.Close()
        return rows.Next(), nil
}

func AddEmployer(db *sql.DB, name string, password string) (bool, error) {
    isValidEmployer, err := CheckEmployer(db, name)
    if(err != nil) {
        return false, err
    }
    if !isValidEmployer {
        stmt, err := db.Prepare("INSERT INTO Employers(name, password) VALUES(?, ?);")
        res, err := stmt.Exec(name, password)
        _ = res
        if err != nil {
            return false, err
        }
    }
    isAddedEmployer, err := addEmployersToCurrentCareerFair(db, name)
    if !isAddedEmployer || err != nil {
        return false, err
    }
    return true, nil
}

func CheckStudent(db *sql.DB, idnumber string) (bool, error){
    stmt, err := db.Prepare("SELECT ID FROM Students WHERE idnumber = ?;")
    if err != nil {
        return false, err
    }
    defer stmt.Close()
    rows, err := stmt.Query(idnumber)
    if err != nil {
        return false, err
    }
    defer rows.Close()
    return rows.Next(), nil
}

func AddStudent(db *sql.DB, displayname string, major string, class string, idnumber string) (bool, error){
    isValidStudent, err := CheckStudent(db, idnumber)
    if err !=  nil {
        return false, err;  
    }
    if !isValidStudent {
        stmt, err := db.Prepare("INSERT INTO Students(idnumber) VALUES(?);")
        res, err := stmt.Exec(idnumber)
        _ = res
        if err != nil {
            return false, err;
        }
    }
    isAddedStudent, err := addStudentToCurrentCareerFair(db, displayname, major, class, idnumber)
    if !isAddedStudent || err != nil {
        return false, err;
    }
    return true, nil;
}

func AddInterview(db *sql.DB, idnumber string, employername string) (bool, error)  {
    stmt, err := db.Prepare("INSERT INTO Interviews (StudentID, EmployerID, CareerFairID) SELECT s.ID, e.ID, c.ID FROM Students s, Employers e, CareerFairs c WHERE s.idnumber = ? AND e.Name = ? AND c.ID = (SELECT MAX(ID) FROM CareerFairs);")
    if err != nil {
        return false, err
    }
    res, err := stmt.Exec(idnumber, employername)
    _ = res
    if err != nil {
            return false, err
    }
    return true, nil
}
//Test Methods Below
func ShowStudents(db *sql.DB, employername string) ([]Interview, error) {
    var (
        displayname string
        major string
        class string
        date string
    )
    stmt, err := db.Prepare("SELECT displayname, major, class, DATE_FORMAT(CheckInTime, \"%h:%i\") FROM Interviews i INNER JOIN StudentsCareerFairs s USING (StudentID) WHERE EmployerID = (SELECT ID FROM Employers WHERE name = ?) AND i.CareerFairID = (SELECT MAX(ID) FROM CareerFairs);")
    if err != nil {
        return nil, err
    }
    defer stmt.Close()
    rows, err := stmt.Query(employername)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    students := make([]Interview, 0, 10)
    for rows.Next() {
        err := rows.Scan(&displayname, &major, &class, &date)
        if err != nil {
            return nil, err
        }
        students = append(students, Interview{Student{displayname, major, class, ""}, date})
    }
    if err = rows.Err(); err != nil {
        return nil, err
    }
    return students, nil
}
//Need to Test Methods Below
func CheckPasswordHash(db *sql.DB, passhash string) (string, error) {
    var name string
    stmt, err := db.Prepare("SELECT name FROM Employers WHERE password = ?")
    if err != nil {
        return "", err
    }
    defer stmt.Close()
    rows, err := stmt.Query(passhash)
    if err != nil {
        return "", err
    }
    defer rows.Close()
    if rows.Next() {
        err := rows.Scan(&name)
        if err != nil {
            return "", err
        }
    } else {
        return "", nil
    }
    if err = rows.Err(); err != nil {
        return "", err
    }
    return name, nil;
}

func ShowEmployersToStudents(db *sql.DB) ([]string, error) {
        var (
                name string
        )
        stmt, err := db.Prepare("SELECT name FROM Employers WHERE ID IN (SELECT EmployerID FROM CareerFairsEmployers WHERE CareerFairID = (SELECT MAX(ID) FROM CareerFairs)) ORDER BY name ASC;")
        if err != nil {
                return nil, err
        }
        defer stmt.Close()
        rows, err := stmt.Query()
        if err != nil {
                return nil, err
        }
        defer rows.Close()
        employers := make([]string, 0, 10)
        for rows.Next() {
                err := rows.Scan(&name)
                if err != nil {
                        return nil, err
                }
                employers = append(employers, name)
        }
        if err = rows.Err(); err != nil {
                return nil, err
        }
        return employers, nil
}

func ShowCareerFairsByName(db *sql.DB) ([]string, error) {
        var (
                name string
        )
        stmt, err := db.Prepare("SELECT name FROM CareerFairs ORDER BY startdate DESC;")
        if err != nil {
                return nil, err
        }
        defer stmt.Close()
        rows, err := stmt.Query()
        if err != nil {
                return nil, err
        }
        defer rows.Close()
        careerfairs := make([]string, 0, 10)
        for rows.Next() {
                err := rows.Scan(&name)
                if err != nil {
                        return nil, err
                }
                careerfairs = append(careerfairs, name)
        }
        if err = rows.Err(); err != nil {
                return nil, err
        }
        return careerfairs, nil
}

func UpdatePassword(db *sql.DB, name string, password string) (bool, error) {
        stmt, err := db.Prepare("UPDATE Employers SET password = ? WHERE name = ?")
        res, err := stmt.Exec(password, name)
        _ = res
        if err != nil {
                return false, err
        }
        return true, nil
}

func addStudentToCurrentCareerFair(db *sql.DB, displayname string, major string, class string, idnumber string) (bool, error) {
	isAdded, err := CheckStudentsCareerFairs(db, idnumber) 
	if err != nil {
		return false, err;
	}
	if isAdded {
		return false, nil;
	}
        stmt, err := db.Prepare("INSERT INTO StudentsCareerFairs(StudentID, CareerFairID, displayname, major, class) VALUES((SELECT ID FROM Students WHERE idnumber = ?), (SELECT MAX(ID) FROM CareerFairs), ?, ?, ?);")
        res, err := stmt.Exec(idnumber, displayname, major, class)
        _ = res
        if err != nil {
                return false, err;
        }
        return true, nil;
}

func addEmployersToCurrentCareerFair(db *sql.DB, name string) (bool, error) {
        stmt, err := db.Prepare("INSERT INTO CareerFairsEmployers(EmployerID, CareerFairID) SELECT e.ID, c.ID FROM Employers e, CareerFairs c WHERE e.name = ? AND c.ID = (SELECT MAX(ID) FROM CareerFairs);")
        res, err := stmt.Exec(name)
        _ = res
        if err != nil {
                return false, err;
        }
        return true, nil;
}

func GetNumberOfStudents(db *sql.DB, fairname string) (int, error) {
    var (
                numStudents int
        )
        stmt, err := db.Prepare("SELECT COUNT(DISTINCT(StudentID)) FROM StudentsCareerFairs WHERE CareerFairID = (SELECT ID FROM CareerFairs WHERE name LIKE ?);")
        if err != nil {
                return -1, err
        }
        defer stmt.Close()
        err = stmt.QueryRow(fairname).Scan(&numStudents)
        if err != nil {
                return -1, err
        }
        return numStudents, nil
}

func GetNumberOfInterviews(db *sql.DB, fairname string) (int, error) {
        var (
                numInterviews int
        )
        stmt, err := db.Prepare("SELECT COUNT(StudentID) FROM Interviews WHERE CareerFairID = (SELECT ID FROM CareerFairs WHERE name LIKE ?);")
        if err != nil {
                return -1, err
        }
        defer stmt.Close()
        err = stmt.QueryRow(fairname).Scan(&numInterviews)
        if err != nil {
                return -1, err
        }
        return numInterviews, nil
}

func StartCareerFair(db *sql.DB, name string, comments string) (bool, error) {
    stmt, err := db.Prepare("INSERT INTO CareerFairs (name, comments) VALUES( ?, ?);")
        if err != nil {
                return false, err
        }
        res, err := stmt.Exec(name, comments)
        _ = res
        if err != nil {
                return false, err
        }
        return true, nil
}

func ShowEmployersInterviewing(db *sql.DB, fairname string) ([]string, error) {
        var (
                name string
        )
        stmt, err := db.Prepare("SELECT name FROM Employers WHERE ID IN (SELECT DISTINCT(EmployerID) FROM Interviews WHERE CareerFairID = (SELECT ID FROM CareerFairs WHERE name LIKE ?)) ORDER BY name ASC;")
        if err != nil {
                return nil, err
        }
        defer stmt.Close()
        rows, err := stmt.Query(fairname)
        if err != nil {
                return nil, err
        }
        defer rows.Close()
        employers := make([]string, 0, 10)
        for rows.Next() {
                err := rows.Scan(&name)
                if err != nil {
                        return nil, err
                }
                employers = append(employers, name)
        }
        if err = rows.Err(); err != nil {
                return nil, err
        }
        return employers, nil
}

func UpdateEmployerName(db *sql.DB, oldname string, newname string) (bool, error) {
        stmt, err := db.Prepare("UPDATE Employers SET name = ? WHERE name = ?;")
        if err != nil {
        return false, err;
    }
    res, err := stmt.Exec(newname, oldname)
        _ = res
        if err != nil {
                return false, err
        }
        return true, nil
}

func DeleteEmployer(db *sql.DB, empname string, fairname string) (bool, error)  {
        stmt, err := db.Prepare("DELETE FROM CareerFairsEmployers WHERE EmployerID = (SELECT ID FROM Employers WHERE name = ?) AND CareerFairID = (SELECT ID FROM CareerFairs WHERE name = ?);")
        if err != nil {
                return false, err
        }
        res, err := stmt.Exec(empname, fairname)
        _ = res
        if err != nil {
                return false, err
        }
        return true, nil
}

func CheckStudentsCareerFairs(db *sql.DB, idnumber string) (bool, error){
    stmt, err := db.Prepare("SELECT StudentID FROM StudentsCareerFairs WHERE StudentID = (SELECT ID FROM Students WHERE idnumber = ?);")
    if err != nil {
        return false, err
    }
    defer stmt.Close()
    rows, err := stmt.Query(idnumber)
    if err != nil {
        return false, err
    }
    defer rows.Close()
    return rows.Next(), nil
}

func ShowAllEmployersByCareerFair(db *sql.DB, fairname string) ([]string, error) {
        var (
                name string
        )
        stmt, err := db.Prepare("SELECT name FROM Employers WHERE ID IN (SELECT EmployerID FROM CareerFairsEmployers WHERE CareerFairID = (SELECT ID FROM CareerFairs WHERE name LIKE ?)) ORDER BY name ASC;")
        if err != nil {
                return nil, err
        }
        defer stmt.Close()
        rows, err := stmt.Query(fairname)
        if err != nil {
                return nil, err
        }
        defer rows.Close()
        employers := make([]string, 0, 10)
        for rows.Next() {
                err := rows.Scan(&name)
                if err != nil {
                        return nil, err
                }
                employers = append(employers, name)
        }
        if err = rows.Err(); err != nil {
                return nil, err
        }
        return employers, nil
}
