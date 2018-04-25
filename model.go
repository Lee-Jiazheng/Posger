package Posger

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

/** Databaser is An abstract factory interface, if you want to add a new database backend, you must implement all function*/
type Databaser interface {
	iUser
}

func GetDatabaseConncect(db string) (Databaser) {
	if db == "mongodb" { return newMongodb()}
	panic("db is not the validate database backend.")
	return nil
}

// anonymous member, we can directly point the final function, such as AddUser.
type Mongodb struct{
	MongodbUser
}

/** The models definition and corresponding interface. */
type User struct {
	// user's id
	UserId		int
	// user'name or oauth2 website username
	Username	string
	// oauth2 authenticated website
	Source 		string
	// password, default is username
	Password	string
	// avatar path, maybe external url
	Avatar		string
	// bio
	Bio			string
	// oauth2 get user's info's url
	InfoURL		string
}

type iUser interface {
	AddUser(user User) (error)
	DeleteUser(user User) (error)
	SelectUser(username string) (User)
}

type MongodbUser struct {
	database 	string
	// collection name, i.e. table name
	collection  string
}

// Mongodb's construction, also construct it's interface implementation struct, i.e. subclass.
func newMongodb() (*Mongodb){
	var mdb = new(Mongodb)
	// construct All databases name and collections name, need to be explicit.
	mdb.MongodbUser.database = "Mongodb";
	mdb.MongodbUser.collection = "Users";	// table name.
	return mdb
}

func (self MongodbUser) AddUser(user User) (error) {
	session, err := mgo.Dial(MONGO_DB_SEVER)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetSafe(&mgo.Safe{})

	c := session.DB(self.database).C(self.collection)
	err = c.Insert(&user)		// It can insert not only one Object, it acquires "...&interface"
	if err != nil {
		return err
	}
	return nil;
}


func (self MongodbUser) DeleteUser(user User) (error) {
	return nil;
}

func (self MongodbUser) SelectUser(username string) (User) {
	session, err := mgo.Dial(MONGO_DB_SEVER)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	c := session.DB(self.database).C(self.collection)
	res := []User{}
	err = c.Find(&bson.M{"username": username}).All(&res)
	if err != nil {
		panic(err)
	}
	if len(res) == 1 {
		return res[0]
	} else {
		return res[0]
	}
}

/**
If
*/

type SqlUser struct {

}


