package Posger

import (
	"testing"
	"fmt"
)

func TestInsertMongodb(t *testing.T) {
	db := NewDatabaseConnecter()
	if err := db.AddUser(User{Username: "ljz", Source: "google"}); err != nil {
		t.Error(err)
	}
	// TODO: Select User to verify we really insert users successfully.
}

func TestSelectMongodb(t *testing.T) {
	db := NewDatabaseConnecter()
	// When you select, the corresponding fields must be lower-case.
	res, err := db.SelectUser(map[string]string{"username": "ljz", "source": "github"});
	if err != nil {	t.Error(err); return }
	for _, user := range res {
		fmt.Print(user.Username)
	}
}
