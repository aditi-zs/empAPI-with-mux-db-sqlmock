package main

import (
	"bytes"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var mock sqlmock.Sqlmock
var err error

func TestGetEmpData(t *testing.T) {
	dB, mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer dB.Close()
	rows := sqlmock.NewRows([]string{"id", "name", "phoneNo", "dept_Id", "dept_name"}).
		AddRow("71bbdbb9-6bde-11ed-aaff-64bc589051b4", "Monika Jaiswal", "6377873927", "1fa46d13-6a50-11ed-90d1-64bc589051b4", "HR Department")
	mock.ExpectQuery("Select e.ID,e.Name,e.PhoneNo,e.DeptID,d.Name from emp e join dept d on e.DeptID=d.DeptID").WillReturnRows(rows)
	tests := []struct {
		description string
		input       string
		expRes      []Employee
		statusCode  int
	}{
		{"All entries are present",
			"",
			[]Employee{
				{"71bbdbb9-6bde-11ed-aaff-64bc589051b4", "Monika Jaiswal",
					"6377873927",
					Department{
						"1fa46d13-6a50-11ed-90d1-64bc589051b4",
						"HR Department",
					},
				},
			},
			200,
		},
	}
	for _, tc := range tests {
		req, err := http.NewRequest("GET", "/emp", nil)
		if err != nil {
			t.Errorf(err.Error()) //err.Error() will return a string
		}
		resRec := httptest.NewRecorder()
		getEmpData(resRec, req)
		var val []Employee
		_ = json.Unmarshal(resRec.Body.Bytes(), &val) //json to go

		assert.Equal(t, tc.statusCode, resRec.Code)
		assert.Equal(t, tc.expRes, val)
	}
}
func TestGetOneEmpData(t *testing.T) {
	dB, mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an erroe '#{err}' was not expected when opening a stub database connection")
	}
	defer dB.Close()

	row := sqlmock.NewRows([]string{"id", "name", "phoneNo", "dept_id", "dept_name"}).
		AddRow("71bbdbb9-6bde-11ed-aaff-64bc589051b4", "Monika Jaiswal", "6377873927", "1fa46d13-6a50-11ed-90d1-64bc589051b4", "HR Department")
	mock.ExpectQuery("Select e.ID,e.Name,e.PhoneNo,e.DeptID,d.Name from emp e join dept d on e.DeptID=d.DeptID WHERE e.id=?").WithArgs("71bbdbb9-6bde-11ed-aaff-64bc589051b4").WillReturnRows(row)

	tests := []struct {
		description string
		input       string
		expRes      Employee
		statusCode  int
	}{
		{"All entries are present",
			"",
			Employee{
				"71bbdbb9-6bde-11ed-aaff-64bc589051b4", "Monika Jaiswal",
				"6377873927",
				Department{
					"1fa46d13-6a50-11ed-90d1-64bc589051b4",
					"HR Department",
				},
			},
			200,
		},
	}
	for _, tc := range tests {
		req, err := http.NewRequest("GET", "/emp/{id}", nil)
		if err != nil {
			t.Errorf(err.Error()) //err.Error() will return a string
		}
		resRec := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{"id": tc.expRes.ID})
		getOneEmpData(resRec, req)

		var val Employee
		_ = json.Unmarshal(resRec.Body.Bytes(), &val) //json to go

		assert.Equal(t, tc.statusCode, resRec.Code)
		assert.Equal(t, tc.expRes, val)
	}

}
func TestGetDepData(t *testing.T) {
	dB, mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer dB.Close()
	rows := sqlmock.NewRows([]string{"dept_Id", "dept_name"}).
		AddRow("1fa46d13-6a50-11ed-90d1-64bc589051b4", "HR Department")
	mock.ExpectQuery("Select * from dept").WillReturnRows(rows)
	tests := []struct {
		description string
		input       string
		expRes      []Department
		statusCode  int
	}{
		{"All entries are present",
			"",
			[]Department{
				{"1fa46d13-6a50-11ed-90d1-64bc589051b4",
					"HR Department"},
			},
			//`[{"dept_id":"1fa46d13-6a50-11ed-90d1-64bc589051b4","dept_name":"HR Department"}]`,
			200,
		},
	}
	for _, tc := range tests {
		req, err := http.NewRequest("GET", "/dept", nil)
		if err != nil {
			t.Errorf(err.Error()) //err.Error() will return a string
		}
		resRec := httptest.NewRecorder()
		getDepData(resRec, req)
		var val []Department
		_ = json.Unmarshal(resRec.Body.Bytes(), &val) //json to go

		assert.Equal(t, tc.statusCode, resRec.Code)
		assert.Equal(t, tc.expRes, val)

	}
}

func TestPostEmployeeData(t *testing.T) {
	dB, mock, err = sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer dB.Close()
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO emp").WithArgs(sqlmock.AnyArg(), "Aditi Verma", "1fa46d13-6a50-11ed-90d1-64bc589051b4", "6388768118").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	tests := []struct {
		description string
		input       Employee
		expRes      Employee
		statusCode  int
	}{
		{"All entries are present",
			Employee{
				"", "Aditi Verma",
				"6388768118",
				Department{
					"1fa46d13-6a50-11ed-90d1-64bc589051b4",
					"",
				},
			},
			Employee{
				"", "Aditi Verma",
				"6388768118",
				Department{
					"1fa46d13-6a50-11ed-90d1-64bc589051b4",
					"",
				},
			},
			201,
		},
	}

	for _, tc := range tests {
		val, _ := json.Marshal(tc.input) //go to json
		req, err := http.NewRequest("POST", "/postempdata", bytes.NewReader(val))
		if err != nil {
			t.Errorf(err.Error())
		}
		//response recorder
		resRec := httptest.NewRecorder()
		postEmployeeData(resRec, req)
		var actRes Employee
		_ = json.Unmarshal(resRec.Body.Bytes(), &actRes) //json to go
		assert.Equal(t, tc.statusCode, resRec.Code)
		//assert.Equal(t, tc.expRes, actRes)
	}
}
