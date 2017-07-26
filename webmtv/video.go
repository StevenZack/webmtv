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
	s, err := mgo.Dial("127.0.0.1")
	if err != nil {
		go RestartMongodb()
		ReturnInfo(w, err.Error(), "")
		return
	}
	defer s.Close()
	cv := s.DB("webmtv").C("videos")
	cc := s.DB("webmtv").C("comments")
	vpd := VideoPageData{}

	err = cv.Find(bson.M{"vid": vid}).One(&vpd.MVideo)
	if err != nil {
		ReturnInfo(w, err.Error(), "")
		return
	}
	err = cc.Find(bson.M{"vid": vid}).Limit(30).Sort("-commenttime").All(&vpd.MComments)
	if err != nil {
		ReturnInfo(w, err.Error(), "")
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
	s, err := mgo.Dial("127.0.0.1")
	if err != nil {
		go RestartMongodb()
		ReturnInfo(w, err.Error(), "")
		return
	}
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
		ReturnInfo(w, err.Error(), "")
		return
	}
	fmt.Fprint(w, u.ID)
}

type PlayListPageData struct {
	Me     PlayList
	Videos []Video
}

func PlayListPage(w http.ResponseWriter, r *http.Request) {
	vid := r.FormValue("vid")
	if vid == "" {
		ReturnInfo(w, "no such playlist", "/")
		return
	}
	s, err := mgo.Dial("127.0.0.1")
	if err != nil {
		go RestartMongodb()
		ReturnInfo(w, err.Error(), "")
		return
	}
	cv := s.DB("webmtv").C("videos")
	cpl := s.DB("webmtv").C("playlists")

	plpd := PlayListPageData{}

	err = cpl.Find(bson.M{"vid": vid}).One(&plpd.Me)
	if err != nil {
		ReturnInfo(w, err.Error(), "")
		return
	}
	err = cv.Find(bson.M{"playlistid": vid}).All(&plpd.Videos)
	if err != nil {
		ReturnInfo(w, err.Error(), "")
		return
	}
	t, err := template.ParseFiles("./html/playlist.html")
	if err != nil {
		fmt.Println(err.Error())
		ReturnInfo(w, err.Error(), "")
		return
	}
	t.Execute(w, plpd)
}
