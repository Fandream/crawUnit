package model

import (
	"crawlab/constants"
	"crawlab/database"
	"github.com/apex/log"
	"github.com/globalsign/mgo/bson"
	"runtime/debug"
	"time"
)

type Action struct {
	Id     bson.ObjectId `json:"_id" bson:"_id"`
	UserId bson.ObjectId `json:"user_id" bson:"user_id"`
	Type   string        `json:"type" bson:"type"`

	CreateTs time.Time `json:"create_ts" bson:"create_ts"`
	UpdateTs time.Time `json:"update_ts" bson:"update_ts"`
}

func (a *Action) Save() error {
	s, c := database.GetCol("actions")
	defer s.Close()

	a.UpdateTs = time.Now()

	if err := c.UpdateId(a.Id, a); err != nil {
		debug.PrintStack()
		return err
	}
	return nil
}

func (a *Action) Add() error {
	s, c := database.GetCol("actions")
	defer s.Close()

	a.Id = bson.NewObjectId()
	a.UpdateTs = time.Now()
	a.CreateTs = time.Now()
	if err := c.Insert(a); err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return err
	}

	return nil
}

