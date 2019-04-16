package dbutil

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

//func main() {
//	db, err := OpenDB()
//	if err != nil {
//		fmt.Println("Failed to Open")
//	}
//	defer CloseDB(db)
//	name, err := CheckPasswordHash(db, "abc123")
//	if err != nil {
//		fmt.Println("Error")
//	}
//	fmt.Println(name)
//}

//Struct Type For Students
type Student struct {
	displayname string
	major string
	class string
	idnumber string
}

type Interview struct {
	student Student
	time string
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
	addEmployersToCurrentCareerFair(db, name)
	if isValidEmployer {
		return false, nil
	}
	stmt, err := db.Prepare("INSERT INTO Employers(name, password) VALUES(?, ?);")
	res, err := stmt.Exec(name, password)
	_ = res
	if err != nil {
		return false, err
	}
	//Add Employer to CareerFairsEmployers Table
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
	if err != nil {
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
	stmt, err := db.Prepare("INSERT INTO Interviews (StudentID, EmployerID) SELECT s.ID, e.ID FROM Students s, Employers e WHERE s.idnumber = ? AND e.Name = ?;")
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
	stmt, err := db.Prepare("SELECT displayname, major, class, DATE_FORMAT(CheckInTime, \"%h:%i\") FROM Interviews INNER JOIN StudentsCareerFairs USING (StudentID) WHERE EmployerID = (SELECT ID FROM Employers WHERE name = ?);")
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

func ShowEmployers(db *sql.DB, employername string) ([]string, error) {
        var (
                name string
        )
        stmt, err := db.Prepare("SELECT name FROM Employers WHERE ID IN (SELECT EmployerID FROM CareerFairsEmployers WHERE CareerFairID = (SELECT MAX(ID) FROM CareerFair));")
        if err != nil {
                return nil, err
        }
        defer stmt.Close()
        rows, err := stmt.Query(employername)
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
        stmt, err := db.Prepare("INSERT INTO StudentsCareerFairs(StudentID, CareerFairID, displayname, major, class) VALUES((SELECT ID FROM Students WHERE idnumber = ?), (SELECT MAX(ID) FROM CareerFairs), ?, ?, ?);")
        res, err := stmt.Exec(idnumber, displayname, major, class)
        _ = res
        if err != nil {
                return false, err;
        }
        return true, nil;
}

func addEmployersToCurrentCareerFair(db *sql.DB, name string) (bool, error) {
        stmt, err := db.Prepare("INSERT INTO EmployersCareerFairs(EmployerID, CareerFairID) SELECT e.ID, c.ID FROM Employers e, CareerFairs c WHERE c.ID = MAX(ID) AND e.name = ?;")
        res, err := stmt.Exec(name)
        _ = res
        if err != nil {
                return false, err;
        }
        return true, nil;
}

func getNumberOfStudents(db *sql.DB, fairname string) (int, error) {
	var (
                numStudents int
        )
        stmt, err := db.Prepare("SELECT COUNT(DISTINCT(StudentID)) FROM StudentsCareerFairs WHERE CareerFairID = (SELECT ID FROM CareerFairs WHERE name LIKE fairname);")
        if err != nil {
                return -1, err
        }
        defer stmt.Close()
        rows, err := stmt.Query()
        if err != nil {
                return -1, err
        }
        defer rows.Close()
        err = rows.Scan(&numStudents)
        if err != nil {
                return -1, err
        }
        if err = rows.Err(); err != nil {
                return -1 , err
        }
        return numStudents, nil
}

func getNumberOfInterviews(db *sql.DB, fairname string) (int, error) {
        var (
                numInterviews int
        )
        stmt, err := db.Prepare("SELECT COUNT(StudentID) FROM Interviews WHERE CareerFairID = (SELECT ID FROM CareerFairs WHERE name LIKE fairname);")
        if err != nil {
                return -1, err
        }
        defer stmt.Close()
        rows, err := stmt.Query()
        if err != nil {
                return -1, err
        }
        defer rows.Close()
        err = rows.Scan(&numInterviews)
        if err != nil {
                return -1, err
        }
        if err = rows.Err(); err != nil {
                return -1, err
        }
        return numInterviews, nil
}
