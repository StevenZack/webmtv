package webmtv

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"net/http"
	"time"
)

type VideoPageData struct {
	MVideo    Video
	MComments []Comment
}

func VideoPage(w http.ResponseWriter, r *http.Request) {
	vid := r.FormValue("vid")
	s, _ := mgo.Dial("127.0.0.1")
	defer s.Close()
	cv := s.DB("webmtv").C("videos")
	cc := s.DB("webmtv").C("comments")
	vpd := VideoPageData{}

	err := cv.Find(bson.M{"vid": vid}).One(&vpd.MVideo)
	if err != nil {
		ReturnInfo(w, err.Error(), false)
		return
	}
	err = cc.Find(bson.M{"vid": vid}).Limit(30).Sort("-commenttime").All(&vpd.MComments)
	if err != nil {
		ReturnInfo(w, err.Error(), false)
		return
	}

	t, _ := template.ParseFiles("./html/video.html")
	t.Execute(w, vpd)
}
func HandleComment(w http.ResponseWriter, r *http.Request) {
	vid := r.FormValue("vid")
	str := r.FormValue("cm")
	sid, err := r.Cookie("WEBMTV-SESSION-ID")
	if err != nil {
		fmt.Fprint(w, "Plz Login first")
		return
	}
	u, err := CheckOutSessionID(sid)
	if err != nil {
		http.SetCookie(w, &http.Cookie{
			Name:    "WEBMTV-SESSION-ID",
			Value:   "",
			Expires: time.Now(),
		})
		fmt.Fprint(w, "Plz Login first")
		return
	}
	s, _ := mgo.Dial("127.0.0.1")
	defer s.Close()
	cc := s.DB("webmtv").C("comments")
	newComment := Comment{
		Data:        str,
		OwnerID:     u.ID,
		CommentTime: time.Now(),
		Vid:         vid,
	}
	err = cc.Insert(&newComment)
	if err != nil {
		ReturnInfo(w, err.Error(), false)
		return
	}
	fmt.Fprint(w, u.ID)
}
