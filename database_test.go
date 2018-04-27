package Posger

import (
	"testing"
	"github.com/satori/go.uuid"
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
	AddPaper(Paper{PaperId: uuid.Must(uuid.NewV4()).String(), Owner:"gajanlee", Path:"2"})
}

func TestSelectPaper(t *testing.T) {
	t.Log(SelectPaper(map[string]interface{}{"owner": "g"}))
}

func TestDeletePaper(t *testing.T) {
	DeletePaper(map[string]interface{}{"paperid": "a57c96ba-24cd-4833-9a8b-6815a2ab757e"})
}