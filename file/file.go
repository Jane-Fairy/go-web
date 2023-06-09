package file

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"
)

//http.HandleFunc("/upload", upload)

func Upload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //获取请求的方法
	if r.Method == "GET" {
		unix := time.Now().Unix()
		hash := md5.New()
		_, err := io.WriteString(hash, strconv.FormatInt(unix, 10))
		if err != nil {
			fmt.Println("server err:", err)
		}
		token := fmt.Sprintf("%x", hash.Sum(nil))
		files, _ := template.ParseFiles("login.gtpl")
		files.Execute(w, token)
	} else {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		fmt.Fprintf(w, "%v", handler.Header)
		f, err := os.OpenFile("./filesystemtest/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
	}
}

//本地模拟客货端
func postFile(filename string, targetUrl string) error {
	bodybuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodybuf)
	fileWriter, err := bodyWriter.CreateFormFile("uploadfile", filename)
	if err != nil {
		fmt.Println("error writing to buffer")
		return err
	}

	fh, err := os.Open(filename)

	if err != nil {
		fmt.Println("error opening file")
		return err
	}
	defer fh.Close()

	io.Copy(fileWriter, fh)

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(targetUrl, contentType, bodybuf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(resp.Status)
	fmt.Println(string(resp_body))
	return nil
}

func Simulation() {
	target_url := "http://localhost:8080/upload"
	filename := "./filesystemtest/8.png"
	postFile(filename, target_url)
}
