package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync/atomic"
)

type user struct {
	Email    string
	Password string
	Tel      string
	Name     string
	Surname  string
	Admin    bool
}

type event struct{
	Evname   string
	Evstatus string
	Evdescpiption string
}

type server struct {
	userToId    map[string]uint32
	users       map[uint32]user
	nextUserId  uint32
	eventToId   map[string]uint32
	events      map[uint32]user
	nextEventId uint32
}

type regHandler struct {
	srv *server
}

func (u *regHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	log.Printf("%+v", req)

	b, _ := ioutil.ReadAll(req.Body)
	log.Printf("%+v", string(b))

	var user user
	if err := json.Unmarshal(b, &user); err != nil {
		log.Print(err)
		rw.WriteHeader(400)
		return
	}

	user.Admin = false;

	myUserId := atomic.AddUint32(&u.srv.nextUserId, 1)
	u.srv.users[myUserId] = user
	u.srv.userToId[user.Email] = myUserId

	rw.WriteHeader(200)

	ret, _ := json.Marshal(myUserId)
	rw.Write(ret)

	log.Printf("%d %+v %s", myUserId, user, ret)
}

type authHandler struct {
	srv *server
}

type authData struct {
	Email    string
	Password string
}

func (u *authHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	log.Printf("%+v", req)

	b, _ := ioutil.ReadAll(req.Body)
	log.Printf("%+v", string(b))

	var request authData
	if err := json.Unmarshal(b, &request); err != nil {
		log.Print(err)
		rw.WriteHeader(500)
		return
	}

	uid, exists := u.srv.userToId[request.Email]
	if !exists {
		rw.WriteHeader(403)
		return
	}

	user, exists := u.srv.users[uid]
	if !exists {
		rw.WriteHeader(403)
		return
	}

	if user.Password != request.Password {
		rw.WriteHeader(403)
		return
	}

	rw.WriteHeader(200)

	ret, _ := json.Marshal(uid)
	rw.Write(ret)
}

type getHandler struct {
	srv *server
}

type  regeventHandler struct{
	srv *server
}

func (u *regeventHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	log.Printf("%+v", req)

	b, _ := ioutil.ReadAll(req.Body)
	log.Printf("%+v", string(b))

	var event event
	if err := json.Unmarshal(b, &event); err != nil {
		log.Print(err)
		rw.WriteHeader(400)
		return
	}

	//user.Admin = false;

	myEventId := atomic.AddUint32(&u.srv.nextEventId, 1)
	u.srv.events[myEventId] = event
	u.srv.eventToId[event.Evdescpiption] = myEventId

	rw.WriteHeader(200)

	ret, _ := json.Marshal(myEventId)
	rw.Write(ret)

	log.Printf("%d %+v %s", myEventId, event, ret)
}


func (u *getHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	qs := req.URL.Query()
	uidStrings := qs["user_id"]

	if len(uidStrings) < 1 {
		rw.WriteHeader(400)
		return
	}

	uid, err := strconv.ParseUint(uidStrings[0], 10, 0)
	if err != nil {
		log.Printf("cannnot convert query string param to int: %s", err)
		rw.WriteHeader(400)
		return
	}

	user, exists := u.srv.users[uint32(uid)]
	if !exists {
		rw.WriteHeader(404)
		return
	}

	log.Printf("%+v", user)
	b, err := json.Marshal(user)
	if err != nil {
		rw.WriteHeader(500)
		return
	}

	rw.WriteHeader(200)
	if _, err := rw.Write(b); err != nil {
		log.Printf("cannot write response: %s", err)
	}
}

func main() {
	srv := &server{
		userToId:   make(map[string]uint32, 1000),
		users:      make(map[uint32]user, 1000),
		nextUserId: 1,
	}

	mux := http.NewServeMux()
	mux.Handle("/reg", &regHandler{srv: srv,})
	mux.Handle("/auth", &authHandler{srv: srv,})
	mux.Handle("/get", &getHandler{srv: srv,})
	mux.Handle("/reg_event",&regeventHandler{srv: srv,})

	hs := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	if err := hs.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
