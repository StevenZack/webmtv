package webmtv

import (
	"crypto/md5"
	"errors"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"time"
)

type LoginData struct {
	InvalidAlert string
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("./html/login.html")
		t.Execute(w, nil)
		return
	}
	s, _ := mgo.Dial("127.0.0.1")
	defer s.Close()
	c := s.DB("webmtv").C("users")
	ld := LoginData{}
	id := r.FormValue("id")
	password := r.FormValue("password")
	if len(id) > 0 && len(id) <= 20 && len(password) >= 5 && len(password) <= 20 {
		u := User{}
		err := c.Find(bson.M{"id": id}).One(&u)
		if err != nil { //register
			ct := time.Now().Unix()
			h := md5.New()
			io.WriteString(h, strconv.FormatInt(ct, 10))
			token := fmt.Sprintf("%x", h.Sum(nil))
			c.Insert(&User{ID: id, Password: password, Sessionid: token})
			http.SetCookie(w, &http.Cookie{Name: "WEBMTV-SESSION-ID", Value: token, Expires: time.Now().AddDate(1, 0, 0)})
		} else if u.Password == password { //login
			http.SetCookie(w, &http.Cookie{Name: "WEBMTV-SESSION-ID", Value: u.Sessionid, Expires: time.Now().AddDate(1, 0, 0)})
		} else {
			t, _ := template.ParseFiles("./html/login.html")
			ld.InvalidAlert = "Wrong Password"
			t.Execute(w, ld)
			return
		}
		ReturnInfo(w, "Succeed", true)
		return
	}
	t, _ := template.ParseFiles("./html/login.html")
	ld.InvalidAlert = "Invalid Input"
	t.Execute(w, ld)
	return
}
func CheckOutSessionID(sid *http.Cookie) (*User, error) {
	s, _ := mgo.Dial("127.0.0.1")
	c := s.DB("webmtv").C("users")
	result := User{}
	e := c.Find(bson.M{"sessionid": sid.Value}).One(&result)
	if e != nil {
		return nil, errors.New("session id not found")
	}
	return &result, nil
}
