package webmtv

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"net/http"
)

type UserPageData struct {
	Me       User
	MyVideos []Video
	IsMyPage bool
}

func UserPage(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	s, _ := mgo.Dial("127.0.0.1")
	defer s.Close()
	upd := UserPageData{}
	cv := s.DB("webmtv").C("videos")
	cu := s.DB("webmtv").C("users")

	err := cu.Find(bson.M{"id": id}).One(&upd.Me)
	if err != nil {
		ReturnInfo(w, err.Error(), false)
		return
	}
	sid, _ := r.Cookie("WEBMTV-SESSION-ID")
	if sid.Value == upd.Me.Sessionid {
		upd.IsMyPage = true
	}
	err = cv.Find(bson.M{"ownerid": id}).Limit(10).Sort("-uploadtime").All(&upd.MyVideos)
	if err != nil {
		ReturnInfo(w, err.Error(), false)
		return
	}
	t, _ := template.ParseFiles("./html/user.html")
	t.Execute(w, upd)
}
