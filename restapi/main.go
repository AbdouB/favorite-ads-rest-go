package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"encoding/json"
)

var a App

func main() {
	a = App{}
	a.Initialize("restfavads_mongo")
	a.Run(":3000")
}

func GetUsers(w http.ResponseWriter, r *http.Request) {

	allUsers := []User{}

	userCollection := a.Session.DB("favoriteAds").C("users")
	userCollection.Find(bson.M{}).Select(bson.M{}).All(&allUsers)

	responseJSON(w, allUsers)
}

func GetAds(w http.ResponseWriter, r *http.Request) {

	allAds := []Ad{}
	adsCollection := a.Session.DB("favoriteAds").C("ads")
	adsCollection.Find(bson.M{}).Select(bson.M{}).All(&allAds)

	responseJSON(w, allAds)

}

func GetUserFavoriteAds(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	userId := params["user_id"]
	saveType := params["save_type"]

	if !bson.IsObjectIdHex(userId) {
		responseError(w, "the user id provided is not well formatted", http.StatusBadRequest)
		return
	}
	userCollection := a.Session.DB("favoriteAds").C("users")

	user_match := bson.M{"$match": bson.M{"_id": bson.ObjectIdHex(userId)}}
	favAds_unwind := bson.M{"$unwind": "$favoriteAds"}
	group_pipe := bson.M{"$group": bson.M{
		"_id":         "$_id",
		"favoriteAds": bson.M{"$push": "$favoriteAds"},
	}}
	favAds_match := bson.M{}
	if saveType == "all" {
		favAds_match = bson.M{"$match": bson.M{"favoriteAds.savetype": bson.M{"$in": []string{"automatic", "manual"}}}}
	} else {
		favAds_match = bson.M{"$match": bson.M{"favoriteAds.savetype": saveType}}
	}

	pipe := userCollection.Pipe([]bson.M{user_match, favAds_unwind, favAds_match, group_pipe})

	result := []bson.M{}
	err := pipe.All(&result)

	if err != nil {
		responseError(w, err.Error(), http.StatusInternalServerError)
	} else if len(result) > 0 {
		responseJSON(w, result[0]["favoriteAds"])
	} else {
		responseJSON(w, []Ad{})
	}

}

func AddUserFavoriteAds(w http.ResponseWriter, r *http.Request) {
	userId := r.FormValue("userId")
	adId := r.FormValue("adId")
	saveType := r.FormValue("saveType")

	if saveType != "manual" && saveType != "automatic" {
		saveType = "automatic"
	}

	if !bson.IsObjectIdHex(userId) || !bson.IsObjectIdHex(adId) {
		responseError(w, "userId or adId are not well formatted", http.StatusBadRequest)
		return
	}

	adToBeFavored := &Ad{}
	userToBeModified := &User{}

	adsCollection := a.Session.DB("favoriteAds").C("ads")
	userCollection := a.Session.DB("favoriteAds").C("users")
	adsCollection.Find(bson.M{"_id": bson.ObjectIdHex(adId)}).Select(bson.M{}).One(&adToBeFavored)
	adToBeFavored.SaveType = saveType

	change := mgo.Change{
		Update:    bson.M{"$push": bson.M{"favoriteAds": adToBeFavored}},
		ReturnNew: true,
	}
	_, err := userCollection.Find(bson.M{"_id": bson.ObjectIdHex(userId)}).Apply(change, &userToBeModified)

	if err == nil {
		responseJSON(w, userToBeModified.FavoriteAds)
	} else {
		responseError(w, err.Error(), http.StatusInternalServerError)
	}
}

func DeleteUserFavoriteAds(w http.ResponseWriter, r *http.Request) {

	userId := r.FormValue("userId")
	adId := r.FormValue("adId")

	if !bson.IsObjectIdHex(userId) || !bson.IsObjectIdHex(adId) {
		responseError(w, "userId or adId are not well formatted", http.StatusBadRequest)
		return
	}
	userToBeModified := &User{}
	userCollection := a.Session.DB("favoriteAds").C("users")

	change := mgo.Change{
		Update:    bson.M{"$pull": bson.M{"favoriteAds": bson.M{"_id": bson.ObjectIdHex(adId)}}},
		ReturnNew: true,
	}
	_, err := userCollection.Find(bson.M{"_id": bson.ObjectIdHex(userId)}).Apply(change, &userToBeModified)

	if err == nil {
		responseJSON(w, userToBeModified.FavoriteAds)
	} else {
		responseError(w, err.Error(), http.StatusInternalServerError)
	}
}

func responseError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func responseJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
