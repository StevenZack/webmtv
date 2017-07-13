package webmtv

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"net/http"
	"strconv"
)

type UserPageData struct {
	Me          User
	MyVideos    []Video
	IsMyPage    bool
	CurrentPage int
	TotalPage   int
}

func UserPage(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	s, _ := mgo.Dial("127.0.0.1")
	defer s.Close()
	upd := UserPageData{CurrentPage: 1, TotalPage: 1}
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
	//handle page
	total, _ := cv.Find(bson.M{"ownerid": id}).Count()
	upd.TotalPage = GetTotalPage(total)
	if r.FormValue("reqPage") != "" {
		upd.CurrentPage, err = strconv.Atoi(r.FormValue("reqPage"))
		if err != nil || upd.CurrentPage < 1 || upd.CurrentPage > upd.TotalPage {
			ReturnInfo(w, "The page you request doesn't exist", false)
			return
		}
	}

	err = cv.Find(bson.M{"ownerid": id}).Limit(30).Skip((upd.CurrentPage - 1) * 30).Sort("-uploadtime").All(&upd.MyVideos)
	if err != nil {
		ReturnInfo(w, err.Error(), false)
		return
	}
	t, _ := template.ParseFiles("./html/user.html")
	t.Execute(w, upd)
}
