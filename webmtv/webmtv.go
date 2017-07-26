package webmtv

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os/exec"
	"strconv"
	"time"
)

type User struct {
	ID        string
	Password  string
	Sessionid string
}
type Video struct {
	Uploadtime    time.Time
	Title         string
	Vid           string
	VURL          string
	Cover         string
	OwnerID       string
	IsWebTorrent  bool
	PlayListTitle string
	PlayListID    string
}
type PlayList struct {
	Title      string
	Vid        string
	Cover      string
	OwnerID    string
	ListLength int
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
type FollowTO struct {
	FromID string
	ToID   string
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
func NewToken() string {
	ct := time.Now().Unix()
	h5 := md5.New()
	io.WriteString(h5, strconv.FormatInt(ct, 10))
	token := fmt.Sprintf("%x", h5.Sum(nil))
	return token
}
