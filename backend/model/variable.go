package model

import (
	"crawunit/database"
	"github.com/globalsign/mgo/bson"
)

/**
全局变量
*/

type Variable struct {
	Id     bson.ObjectId `json:"_id" bson:"_id"`
	Key    string        `json:"key" bson:"key"`
	Value  string        `json:"value" bson:"value"`
	Remark string        `json:"remark" bson:"remark"`
}


func GetVariableList() []Variable {
	s, c := database.GetCol("variable")
	defer s.Close()

	var list []Variable
	if err := c.Find(nil).All(&list); err != nil {

	}
	return list
}
