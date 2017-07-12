package webmtv

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"net/http"
)

type UserPageData struct {
	Me       User
	MyVideos []Video
}

func UserPage(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	s, _ := mgo.Dial("127.0.0.1")
	defer s.Close()
	upd := UserPageData{}
	cv := s.DB("webtv").C("videos")
	cu := s.DB("webtv").C("users")

	err := cu.Find(bson.M{"id": id}).One(&upd.Me)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	err = cv.Find(bson.M{"ownerid": id}).Limit(10).Sort("-uploadtime").All(&upd.MyVideos)
	if err != nil {
		ReturnInfo(w, err.Error(), false)
		return
	}
	fmt.Println(len(upd.MyVideos))
	t, _ := template.ParseFiles("./html/user.html")
	t.Execute(w, upd)
}
