package main

import (
	"./webmtv"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", webmtv.Home)
	http.HandleFunc("/v", webmtv.VideoPage)
	http.HandleFunc("/u", webmtv.UserPage)
	http.HandleFunc("/login", webmtv.Login)
	http.HandleFunc("/upload", webmtv.Upload)
	http.HandleFunc("/editvideo", webmtv.EditVideo)
	http.HandleFunc("/deletevideo", webmtv.DeleteVideo)
	http.HandleFunc("/comment", webmtv.HandleComment)
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))
	err := http.ListenAndServe(":8090", nil)
	if err != nil {
		fmt.Println(err)
	}
}
