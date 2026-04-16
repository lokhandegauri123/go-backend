package main

import (
	"fmt"
	"net/http"
	"log"
)


func main(){
	http.HandleFunc("/",func(w http.ResponseWriter , r *http.Request){
		fmt.Fprintln(w, "Hello from Go server!")
	})

	fmt.Println("starting server...")

	log.Fatal(http.ListenAndServe(":8080",nil))

	// if err != nil{
	// 	fmt.Println(err)
	// }

}