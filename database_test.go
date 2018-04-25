package Posger

import (
	"testing"
)

func TestInsertMongodb(t *testing.T) {
	if err := AddUser(User{Username: "ljz+", Source: "google+"}); err != nil {
		t.Error(err)
	}
	// TODO: Select User to verify we really insert users successfully.
}

func TestSelectMongodb(t *testing.T) {
	user := SelectUser("gajanlee");
	if &user == nil {
		t.Error(user)
	}
	// TODO: Select User to verify we really insert users successfully.

}