package main

import (
	"database/sql"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Department struct {
	DeptID   string `json:"dept_id"`
	DeptName string `json:"dept_name"`
}

type Employee struct {
	ID      string     `json:"id"`
	Name    string     `json:"name"`
	PhoneNo string     `json:"phoneNo"`
	Dept    Department `json:"dept"`
}

func GetEmployeeData(db *sql.DB) ([]Employee, error) {
	rows, err := db.Query("Select e.ID,e.Name,e.PhoneNo,e.DeptID,d.Name from emp e join dept d on e.DeptID=d.DeptID	")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var employees []Employee
	for rows.Next() {
		var e Employee
		err = rows.Scan(&e.ID, &e.Name, &e.PhoneNo, &e.Dept.DeptID, &e.Dept.DeptName)
		if err != nil {
			return nil, err
		}

		employees = append(employees, e)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	//respBody,_:=json.Marshal(employees)
	//w.Write(respBody)
	return employees, nil
}

func getEmpData(w http.ResponseWriter, r *http.Request) {
	val, err := GetEmployeeData(dB)
	if err != nil {
		log.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	respBody, err := json.Marshal(val)
	w.Write(respBody)
	if err != nil {
		log.Println(err)
	}
}
func GetOneEmployeeData(db *sql.DB, id string) (Employee, error) {
	var e Employee
	row := db.QueryRow("Select e.ID,e.Name,e.PhoneNo,e.DeptID,d.Name from emp e join dept d on e.DeptID=d.DeptID WHERE e.id=?", id)
	err := row.Scan(&e.ID, &e.Name, &e.PhoneNo, &e.Dept.DeptID, &e.Dept.DeptName)
	if err != nil {
		return Employee{}, err
	}
	return e, nil
}

func getOneEmpData(w http.ResponseWriter, r *http.Request) {
	empID := mux.Vars(r)["id"]
	val, err := GetOneEmployeeData(dB, empID)
	if err != nil {
		log.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	respBody, err := json.Marshal(val)
	w.Write(respBody)
	if err != nil {
		log.Println(err)
	}
}
func postEmployeeData(w http.ResponseWriter, r *http.Request) {
	var emp Employee
	w.Header().Set("Content-Type", "application/json")
	req, _ := ioutil.ReadAll(r.Body)
	_ = json.Unmarshal(req, &emp)
	newUUID := uuid.New()
	//newUUID, err := exec.Command("uuidgen").Output()
	nUUID := strings.TrimSpace(newUUID.String())
	_, err := dB.Exec("insert into emp values (?,?,?,?)", nUUID, emp.Name, emp.Dept.DeptID, emp.PhoneNo)
	if err != nil {
		log.Println(err)
	}
	w.WriteHeader(http.StatusCreated)
	emp.ID = nUUID
	respBody, err := json.Marshal(emp)
	w.Write(respBody)
	//json.NewEncoder(w).Encode(emp)
}
func postDepartmentData(w http.ResponseWriter, r *http.Request) {
	var dep Department
	w.Header().Set("Content-Type", "application/json")
	req, _ := ioutil.ReadAll(r.Body)
	_ = json.Unmarshal(req, &dep)
	_, err := dB.Exec("insert into dept values (uuid(),?)", dep.DeptName)
	if err != nil {
		log.Println(err)
	}
	w.WriteHeader(http.StatusCreated)

}
func GetDeptData(db *sql.DB) ([]Department, error) {
	rows, err := db.Query("Select * from dept")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var departments []Department
	for rows.Next() {
		var d Department
		err = rows.Scan(&d.DeptID, &d.DeptName)
		if err != nil {
			return nil, err
		}
		departments = append(departments, d)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return departments, nil
}

func getDepData(w http.ResponseWriter, r *http.Request) {
	val, err := GetDeptData(dB)
	if err != nil {
		log.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	respBody, err := json.Marshal(val)
	w.Write(respBody)
	if err != nil {
		log.Println(err)
	}
}
