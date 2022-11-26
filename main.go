package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	_ "net/http"
)

var dB *sql.DB

func main() {
	var err error
	dB, err = dbConnection("mysql", "root:Aditi#2#@tcp(127.0.0.1:3306)/employee")
	if err != nil {
		log.Println(err)
		return
	}

	defer dB.Close()
	//h:= emp.Newhandle(emp.DB)
	//err = emp.DB.Ping()
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	router := mux.NewRouter()
	router.HandleFunc("/emp", getEmpData).Methods("GET")
	router.HandleFunc("/emp/{id}", getOneEmpData).Methods("GET")
	router.HandleFunc("/postempdata", postEmployeeData).Methods("POST")
	router.HandleFunc("/postdepdata", postDepartmentData).Methods("POST")
	router.HandleFunc("/dept", getDepData).Methods("GET")
	fmt.Println(("server at port 8000"))
	log.Fatal(http.ListenAndServe(":8000", router))
}
