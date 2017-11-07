package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"htest/model"

	"github.com/gorilla/mux"
)

const (
	MaxUint32 uint32 = 1<<32 - 1
	MaxUint64 uint64 = 1<<64 - 1
)

type DTO struct {
	BigNumber int64  `json:"bigNumber"`
	Text      string `json:"text"`
}

var (
	Auth map[string]string
	Work map[string]int32
)

func main() {
	route := mux.NewRouter()

	route.HandleFunc("/", mainPage)
	route.HandleFunc("/login", login).Methods("POST")
	route.HandleFunc("/login/pass", changePass).Methods("POST")
	route.HandleFunc("/task", doWork).Methods("POST")

	fmt.Println("Start server on port 4000")
	http.ListenAndServe(":4000", route)
}

func init() {
	Auth = make(map[string]string, 0)
	Work = make(map[string]int32, 1)
	Work["admin"] = 1000000
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`<!DOCTYPE html>
		<html>
		<head>
		<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<meta name="theme-color" content="#375EAB">

			<title>main page</title>
		</head>
		<body>
			Page body and some more content
		</body>
		</html>`))
}

func login(w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("login")
	pass := r.FormValue("pass")

	if Auth[login] == pass {
		w.WriteHeader(http.StatusOK)
	}

	user := &model.User{}
	err := user.Get(login, pass)

	if err == nil {
		Auth[login] = pass
		Work[login] = user.WorkNumber
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func changePass(w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("login")
	pass := r.FormValue("pass")
	newPass := r.FormValue("newPass")

	if Auth[login] != pass {
		w.WriteHeader(http.StatusBadRequest)
	}

	user := &model.User{}
	err := user.Get(login, pass)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		user.Pass = newPass
		err = user.Save()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func doWork(w http.ResponseWriter, r *http.Request) {
	value := DTO{}
	login := r.FormValue("login")
	if Work[login] <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	json.Unmarshal([]byte(r.FormValue("value")), &value)

	v := reflect.ValueOf(value)
	for i := 0; i < v.NumField(); i++ {
		w.Write(reverse(v.Field(i)))
	}
}

func reverse(val reflect.Value) []byte {
	switch val.Kind().String() {
	case "int64":
		result := make([]byte, 8)
		binary.LittleEndian.PutUint64(result, MaxUint64-uint64(val.Interface().(int64)))
		return result
	case "int32":
		result := make([]byte, 4)
		binary.LittleEndian.PutUint32(result, MaxUint32-uint32(val.Interface().(int32)))
		return result
	case "string":
		var result string
		for _, v := range val.Interface().(string) {
			result = string(v) + result
		}
		return []byte(result)
	}
	return nil
}
