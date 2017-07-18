package webmtv

import (
	"crypto/md5"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"time"
)

type UploadData struct {
}

func Upload(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("./html/upload.html")
		t.Execute(w, nil)
		return
	}
	sid, err := r.Cookie("WEBMTV-SESSION-ID")
	if err != nil {
		http.SetCookie(w, &http.Cookie{Name: "WEBMTV-SESSION-ID", Value: "", Expires: time.Now()})
		ReturnInfo(w, err.Error(), "/login")
		return
	}
	u, err := CheckOutSessionID(sid)
	if err != nil {
		http.SetCookie(w, &http.Cookie{Name: "WEBMTV-SESSION-ID", Value: "", Expires: time.Now()})
		ReturnInfo(w, err.Error(), "/login")
		return
	}

	r.ParseMultipartForm(1 << 20)
	mVideo := r.FormValue("video")
	mCover := r.FormValue("cover")
	title := r.FormValue("title")
	isWebTorrent := r.FormValue("videoType") == "webtorrent"

	ct := time.Now().Unix()
	h5 := md5.New()
	io.WriteString(h5, strconv.FormatInt(ct, 10))
	token := fmt.Sprintf("%x", h5.Sum(nil))

	//store info in mongodb
	s, err := mgo.Dial("127.0.0.1")
	if err != nil {
		go RestartMongodb()
		ReturnInfo(w, err.Error(), "")
		return
	}
	defer s.Close()
	mgoNewVideo := s.DB("webmtv").C("videos")
	newVideo := Video{
		Uploadtime:   time.Now(),
		Title:        title,
		Vid:          token,
		VURL:         mVideo,
		Cover:        mCover,
		OwnerID:      u.ID,
		IsWebTorrent: isWebTorrent,
	}
	err = mgoNewVideo.Insert(&newVideo)
	if err != nil {
		ReturnInfo(w, err.Error(), "")
		return
	}
	//return upload-succeed page
	ReturnInfo(w, "succeed", "/")
}
func EditVideo(w http.ResponseWriter, r *http.Request) {
	vid := r.FormValue("vid")
	sid, err := r.Cookie("WEBMTV-SESSION-ID")
	if err != nil {
		http.SetCookie(w, &http.Cookie{Name: "WEBMTV-SESSION-ID", Value: "", Expires: time.Now()})
		ReturnInfo(w, "please log in first", "/login")
		return
	}
	u, err := CheckOutSessionID(sid)
	if err != nil {
		http.SetCookie(w, &http.Cookie{Name: "WEBMTV-SESSION-ID", Value: "", Expires: time.Now()})
		ReturnInfo(w, "please log in first", "/login")
		return
	}
	fundVideo := Video{}
	s, err := mgo.Dial("127.0.0.1")
	if err != nil {
		go RestartMongodb()
		ReturnInfo(w, err.Error(), "")
		return
	}
	defer s.Close()
	cv := s.DB("webmtv").C("videos")
	err = cv.Find(bson.M{"vid": vid}).One(&fundVideo)
	if err != nil {
		ReturnInfo(w, "No such video", "")
		return
	}
	if fundVideo.OwnerID != u.ID {
		ReturnInfo(w, "You don't have permission to edit this video", "")
		return
	}
	if r.Method == "GET" { // GET
		t, _ := template.ParseFiles("./html/editvideo.html")
		t.Execute(w, fundVideo)
		return
	}
	//POST
	err = cv.Update(bson.M{"vid": vid}, bson.M{"$set": bson.M{
		"title":        r.FormValue("title"),
		"vurl":         r.FormValue("video"),
		"cover":        r.FormValue("cover"),
		"iswebtorrent": r.FormValue("videoType") == "webtorrent",
	}})
	if err != nil {
		ReturnInfo(w, "Update video info failed:"+err.Error(), "")
		return
	}
	ReturnInfo(w, "Succeed", "/u?id="+u.ID)
}
func DeleteVideo(w http.ResponseWriter, r *http.Request) {
	vid := r.FormValue("vid")
	sid, err := r.Cookie("WEBMTV-SESSION-ID")
	if err != nil {
		http.SetCookie(w, &http.Cookie{Name: "WEBMTV-SESSION-ID", Value: "", Expires: time.Now()})
		ReturnInfo(w, "please log in first", "/login")
		return
	}
	u, err := CheckOutSessionID(sid)
	if err != nil {
		http.SetCookie(w, &http.Cookie{Name: "WEBMTV-SESSION-ID", Value: "", Expires: time.Now()})
		ReturnInfo(w, "please log in first", "/login")
		return
	}
	fundVideo := Video{}
	s, err := mgo.Dial("127.0.0.1")
	if err != nil {
		go RestartMongodb()
		ReturnInfo(w, err.Error(), "")
		return
	}
	defer s.Close()
	cv := s.DB("webmtv").C("videos")
	err = cv.Find(bson.M{"vid": vid}).One(&fundVideo)
	if err != nil {
		ReturnInfo(w, "No such video", "")
		return
	}
	if fundVideo.OwnerID != u.ID {
		ReturnInfo(w, "You don't have permission to edit this video", "")
		return
	}
	err = cv.Remove(bson.M{"vid": vid})
	if err != nil {
		ReturnInfo(w, "delete failed:"+err.Error(), "")
		return
	}
	ReturnInfo(w, "Succeed", "/u?id="+u.ID)
}
