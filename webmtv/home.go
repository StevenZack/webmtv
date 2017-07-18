package webmtv

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"html/template"
	"net/http"
	"time"
)

type HomeData struct {
	Loggedin bool
	New      []Video
	Me       string
}

func Home(rw http.ResponseWriter, req *http.Request) {
	hd := HomeData{}

	s, e := mgo.Dial("127.0.0.1")
	if e != nil {
		go RestartMongodb()
		fmt.Fprint(rw, e)
		return
	}
	vc := s.DB("webmtv").C("videos")
	vc.Find(nil).Limit(30).Sort("-uploadtime").All(&hd.New)
	s.Close()
	sid, err := req.Cookie("WEBMTV-SESSION-ID")
	if err == nil {
		user, e := CheckOutSessionID(sid)
		if e != nil {
			http.SetCookie(rw, &http.Cookie{Name: "WEBMTV-SESSION-ID", Value: "", Expires: time.Now()})
		} else {
			hd.Me = user.ID
			hd.Loggedin = true
		}
	}
	t, _ := template.ParseFiles("./html/home.html")
	t.Execute(rw, hd)
}
