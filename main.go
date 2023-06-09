package main

import (
	"fmt"
	_ "go-web/database"
	"go-web/file"
	"go-web/session"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
)

var globalSessions *session.Manager

func init() {

	log.Println("init manager")
	globalSessions, _ = session.NewManager("memory", "gosessionid", 3600)
	globalSessions.GC() // 定期清除过期的session
	log.Println(globalSessions)
}

func sessionLogin(w http.ResponseWriter, r *http.Request) {
	sess := globalSessions.SessionStart(w, r)
	r.ParseForm()
	if r.Method == "GET" {
		log.Println("GET REQUEST")
		t, _ := template.ParseFiles("login.gtpl")
		w.Header().Set("Content-Type", "text/html")
		t.Execute(w, sess.Get("username"))
	} else {
		sess.Set("username", r.Form["username"])
		http.Redirect(w, r, "/", 302)
	}

}

func main() {
	http.HandleFunc("/sayHelloName", sayHello)
	http.HandleFunc("/sessionLogin", sessionLogin)
	http.HandleFunc("/login", login)
	http.HandleFunc("/upload", file.Upload)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("listenAndServer:", err)
	}
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

func count(w http.ResponseWriter, r *http.Request) {
	session := globalSessions.SessionStart(w, r)
	createTime := session.Get("create_time")
	if createTime == nil {
		session.Set("create_time", time.Now().Unix())
	} else if (createTime.(int64) + 360) < (time.Now().Unix()) {

	}
}
