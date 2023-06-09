package main

import (
	"fmt"
	_ "go-web/database"
	"go-web/file"
	"html/template"
	"log"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/sayHelloName", sayHello)
	http.HandleFunc("/login", login)
	http.HandleFunc("/upload", file.Upload)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("listenAndServer:", err)
	}

	file.Simulation()
}

func sayHello(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()       //解析参数，默认是不会解析的
	fmt.Println(r.Form) //这些信息是输出到服务器端的打印信息
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello astaxie!") //这个写入到w的是输出到客户端的
}

func login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	log.Println("login")
	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		files, _ := template.ParseFiles("login.gtpl")
		log.Println(files.Execute(w, nil))
	} else {
		fmt.Println("username:", r.Form["username"])
		fmt.Println("password:", r.Form["password"])
	}

}
