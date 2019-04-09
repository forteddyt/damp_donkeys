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
        stmt, err := db.Prepare("SELECT name FROM Employers WHERE name = ?;")
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
	if isValidEmployer {
		return false, nil
	}
	stmt, err := db.Prepare("INSERT INTO Employers(name, password) VALUES(?, ?);")
	res, err := stmt.Exec(name, password)
	_ = res
	if err != nil {
		return false, err
	}
	return true, nil
}

func CheckStudent(db *sql.DB, idnumber string) (bool, error){
	stmt, err := db.Prepare("SELECT displayname FROM Students WHERE idnumber = ?;")
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
	if isValidStudent {
		return false, nil;
	}	
        stmt, err := db.Prepare("INSERT INTO Students(displayname, major, class, idnumber) VALUES(?, ?, ?, ?);")
        res, err := stmt.Exec(displayname, major, class, idnumber)
	_ = res
        if err != nil {
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

func ShowStudents(db *sql.DB, employername string) ([]Student, error) {
	var (
		displayname string
		major string
		class string
	)
	stmt, err := db.Prepare("SELECT displayname, major, class FROM Students WHERE ID IN (SELECT StudentID FROM Interviews WHERE EmployerID = (SELECT ID FROM Employers WHERE name = ?));")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(employername)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	students := make([]Student, 0, 10)
	for rows.Next() {
		err := rows.Scan(&displayname, &major, &class)
		if err != nil {
			return nil, err
		}
		students = append(students, Student{displayname, major, class, ""})
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return students, nil
}

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

