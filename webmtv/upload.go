package webmtv

import (
	"crypto/md5"
	"fmt"
	"gopkg.in/mgo.v2"
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
		ReturnInfo(w, err.Error(), true)
		return
	}
	u, err := CheckOutSessionID(sid)
	if err != nil {
		http.SetCookie(w, &http.Cookie{Name: "WEBMTV-SESSION-ID", Value: "", Expires: time.Now()})
		ReturnInfo(w, err.Error(), true)
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
	s, _ := mgo.Dial("127.0.0.1")
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
		ReturnInfo(w, err.Error(), false)
		return
	}
	//return upload-succeed page
	ReturnInfo(w, "succeed", true)
}
