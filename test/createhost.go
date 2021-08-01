package main

import (
	"net/http"
	_ "fmt"

)

func simple(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("Recieved on "+r.Host))
}

func main(){
	go func(){http.ListenAndServe(":50001",http.HandlerFunc(simple))}()
	go func(){http.ListenAndServe(":50002",http.HandlerFunc(simple))}()
	go func(){http.ListenAndServe(":50003",http.HandlerFunc(simple))}()
	go func(){http.ListenAndServe(":50004",http.HandlerFunc(simple))}()
	go func(){http.ListenAndServe(":50005",http.HandlerFunc(simple))}()

	for {  }
}