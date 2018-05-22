package main

import "gopkg.in/mgo.v2/bson"

type User struct {
	ID          bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Email       string        `json:"email,omitempty"`
	FavoriteAds []Ad          `json:"favoriteAds,omitempty" bson:"favoriteAds"`
}

type Ad struct {
	ID          bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Title       string        `json:"title,omitempty"`
	Description string        `json:"description,omitempty"`
	SaveType    string        `json:"saveType,omitempty"`
}
