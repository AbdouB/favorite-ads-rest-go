package main

import (
	"gopkg.in/mgo.v2"
	"github.com/gorilla/mux"
	"log"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

type App struct {
	Session     *mgo.Session
	Router *mux.Router
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/users", GetUsers).Methods("GET")
	a.Router.HandleFunc("/ads", GetAds).Methods("GET")
	a.Router.HandleFunc("/user/{user_id}/favoriteads/{save_type}", GetUserFavoriteAds).Methods("GET")
	a.Router.HandleFunc("/user/favoriteads", AddUserFavoriteAds).Methods("POST")
	a.Router.HandleFunc("/user/favoriteads", DeleteUserFavoriteAds).Methods("DELETE")
}

func (a *App) Initialize(dbServer string) {

	a.Session = GetMongoSession(dbServer)

	log.Println("Setting things up...clearing left overs...")
	a.Session.DB("favoriteAds").C("users").RemoveAll(bson.M{})
	a.Session.DB("favoriteAds").C("ads").RemoveAll(bson.M{})

	userErr := a.Session.DB("favoriteAds").C("users").Insert(
		bson.M{"email": "user1@email.com"},
		bson.M{"email": "user2@email.com"},
		bson.M{"email": "user3@email.com"},
	)

	adErr := a.Session.DB("favoriteAds").C("ads").Insert(
		bson.M{"title": "Awesome title", "description": "awesome description"},
		bson.M{"title": "Amazing title", "description": "amazing description"},
		bson.M{"title": "Great title", "description": "great description"},

	)

	if userErr != nil || adErr != nil {
		log.Fatal(userErr, adErr);
	}
	log.Println("Mock data added successfully!")

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func GetMongoSession(host string) (*mgo.Session) {

	session, err := mgo.Dial(host)

	if err != nil {
		log.Fatal(err)
	}

	return session
}
