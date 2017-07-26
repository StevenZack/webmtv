package webmtv

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"net/http"
	"time"
)

type UploadData struct {
	PlayLists []PlayList
}

func Upload(w http.ResponseWriter, r *http.Request) {
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
	s, err := mgo.Dial("127.0.0.1")
	if err != nil {
		go RestartMongodb()
		ReturnInfo(w, err.Error(), "")
		return
	}
	defer s.Close()

	mgoNewVideo := s.DB("webmtv").C("videos")
	cpl := s.DB("webmtv").C("playlists")

	ud := UploadData{}
	if r.Method == "GET" {
		err := cpl.Find(bson.M{"ownerid": u.ID}).All(&ud.PlayLists)
		if err != nil {
			ReturnInfo(w, err.Error(), "")
			return
		}
		t, _ := template.ParseFiles("./html/upload.html")
		t.Execute(w, &ud)
		return
	}

	mVideo := r.FormValue("video")
	mCover := r.FormValue("cover")
	title := r.FormValue("title")
	isWebTorrent := r.FormValue("videoType") == "webtorrent"
	pl := r.FormValue("playlist")

	vtoken := NewToken()

	err = cpl.Find(bson.M{"ownerid": u.ID, "title": pl}).All(&ud.PlayLists)
	if err == nil && len(ud.PlayLists) > 0 { //insert into existed playlist
		newVideo := Video{
			Uploadtime:    time.Now(),
			Title:         title,
			Vid:           vtoken,
			VURL:          mVideo,
			Cover:         mCover,
			OwnerID:       u.ID,
			IsWebTorrent:  isWebTorrent,
			PlayListID:    ud.PlayLists[0].Vid,
			PlayListTitle: pl,
		}
		err = mgoNewVideo.Insert(&newVideo)
		if err != nil {
			ReturnInfo(w, err.Error(), "")
			return
		}
		err = cpl.Update(bson.M{"ownerid": u.ID, "title": pl}, bson.M{"$set": bson.M{"listlength": len(ud.PlayLists) + 1}})
		if err != nil {
			ReturnInfo(w, err.Error(), "")
			return
		}
		//return upload-succeed page
		ReturnInfo(w, "succeed", "/")
		return
	}
	//insert into a new playlist
	pltoken := NewToken()

	newVideo := Video{
		Uploadtime:    time.Now(),
		Title:         title,
		Vid:           vtoken,
		VURL:          mVideo,
		Cover:         mCover,
		OwnerID:       u.ID,
		IsWebTorrent:  isWebTorrent,
		PlayListID:    pltoken,
		PlayListTitle: pl,
	}
	err = mgoNewVideo.Insert(&newVideo)
	if err != nil {
		ReturnInfo(w, err.Error(), "")
		return
	}
	newPlayList := PlayList{
		Title:      pl,
		Vid:        pltoken,
		Cover:      mCover,
		OwnerID:    u.ID,
		ListLength: 1,
	}
	err = cpl.Insert(&newPlayList)
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
