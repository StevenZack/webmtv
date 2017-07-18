package webmtv

import (
	"html/template"
	"net/http"
	"os/exec"
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
	Info   string
	Jump   bool
	JmpUrl string
}

func ReturnInfo(w http.ResponseWriter, err string, jmpUrl string) {
	t, _ := template.ParseFiles("./html/info.html")
	if jmpUrl == "" {
		t.Execute(w, &InfoData{Info: err, Jump: false})
	} else {
		t.Execute(w, &InfoData{Info: err, Jump: true, JmpUrl: jmpUrl})
	}
}
func GetTotalPage(num int) int {
	a := num / 30
	if num%30 > 0 {
		return a + 1
	}
	return a
}
func RestartMongodb() {
	exec.Command("systemctl", "restart", "mongodb").Run()
}
