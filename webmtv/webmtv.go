package webmtv

import (
	"html/template"
	"net/http"
	"time"
)

type User struct {
	ID        string
	Password  string
	Sessionid string
}
type Video struct {
	Uploadtime   time.Time
	Title        string
	Vid          string
	VURL         string
	Cover        string
	OwnerID      string
	IsWebTorrent bool
}
type Comment struct {
	Data        string
	OwnerID     string
	CommentTime time.Time
	Vid         string
}
type InfoData struct {
	Info       string
	JumpToHome bool
}

func ReturnInfo(w http.ResponseWriter, err string, b bool) {
	t, _ := template.ParseFiles("./html/info.html")
	t.Execute(w, &InfoData{Info: err, JumpToHome: b})
}
