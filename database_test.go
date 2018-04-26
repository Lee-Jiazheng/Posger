package Posger

import (
	"testing"
	"gopkg.in/mgo.v2/bson"
)

func TestInsertMongodb(t *testing.T) {
	AddUser(User{Username: "ljz+", Source: "google+"});	// internal panic
	// TODO: Select User to verify we really insert users successfully.
}

func TestSelectMongodb(t *testing.T) {
	user := SelectUser(map[string]interface{}{"username": "gajanlee"});
	if len(user) == 0 {
		t.Error(user)
	}
	// TODO: Select User to verify we really insert users successfully.

}

func TestInsertPaper(t *testing.T) {
	AddPaper(Paper{PaperId: bson.NewObjectId(), Owner:"g", Path:"2"})
}

func TestSelectPaper(t *testing.T) {
	t.Log(SelectPaper(map[string]interface{}{"owner": "g"}))
}