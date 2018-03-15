package Posger

import (
	"gopkg.in/mgo.v2"
)

type msg struct {
	// Msg   string        `bson:"msg"`
	Count int `bson:"count"`
}

type Person struct {
	Name  string
	Phone string
}

const (
	MONGO_DB_SEVER = "mongodb://localhost:27017"
)

type mongoDB struct {

}

// It is a factory function to adapt variety databases.
func NewDatabaseConnecter() (databaser Databaser) {
	return mongoDB{}
}

func (self mongoDB) AddUser(user User) (error) {
	session, err := mgo.Dial(MONGO_DB_SEVER)
	if err != nil {
		return err
	}
	defer session.Close()
	session.SetSafe(&mgo.Safe{})

	c := session.DB("MongoTest11").C("user")
	err = c.Insert(&user)		// It can insert not only one Object, it acquires "...&interface"
	if err != nil {
		return err
	}
	return nil
}

func (self mongoDB) DeleteUser(user User) (error) {
	return nil
}

func (self mongoDB) SelectUser(pairs map[string]string) ([]User, error) {
	result := []User{}
	for _, _ = range pairs {
		session, err := mgo.Dial(MONGO_DB_SEVER)
		if err != nil {
			return nil, err
		}
		defer session.Close()
		session.SetSafe(&mgo.Safe{})

		c := session.DB("MongoTest11").C("user")
		err = c.Find(pairs).All(&result)
		if err != nil {
			return nil, err
		}
	}
	return result, nil

}
