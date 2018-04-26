package Posger

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

/** Databaser is An abstract factory interface, if you want to add a new database backend, you must implement all function*/
type Databaser interface {
	iUser
	iPaper
}

func GetDatabaseConncect(db string) (Databaser) {
	if db == "mongodb" { return newMongodb()}
	Logger.Fatalln("db is not the validate database backend.")
	return nil
}

func getMongodbSession() (*mgo.Session){
	session, err := mgo.Dial(MONGO_DB_SEVER)
	if err != nil {
		Logger.Fatalln(err)
	}
	session.SetSafe(&mgo.Safe{})
	return session
}

// Mongodb's construction, also construct it's interface implementation struct, i.e. subclass.
func newMongodb() (*Mongodb){
	var mdb = new(Mongodb)
	// construct All databases name and collections name, need to be explicit.
	mdb.MongodbUser.database = "Mongodb";
	mdb.MongodbUser.collection = "Users";	// table name.
	mdb.MongodbPaper.database = "Mongodb";
	mdb.MongodbPaper.collection = "Papers"
	return mdb
}

// anonymous member, we can directly point the final function, such as AddUser.
type Mongodb struct{
	MongodbUser
	MongodbPaper
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
	SelectUser(filter map[string]interface{}) ([]User)
}

type MongodbUser struct {
	database 	string
	// collection name, i.e. table name
	collection  string
}

func (self MongodbUser) AddUser(user User) (error) {
	session := getMongodbSession()
	defer session.Close()

	c := session.DB(self.database).C(self.collection)
	err := c.Insert(&user)		// It can insert not only one Object, it acquires "...&interface"
	if err != nil {
		return err
	}
	return nil;
}

// TODO: The account has been written off.
func (self MongodbUser) DeleteUser(user User) (error) {
	return nil;
}

// Selection filter condition.
func (self MongodbUser) SelectUser(filter map[string]interface{}) ([]User) {
	session := getMongodbSession()
	defer session.Close()
	c := session.DB(self.database).C(self.collection)
	res := []User{}
	err := c.Find(&filter).All(&res)
	if err != nil {
		panic(err)
	}
	return res
}

/** Next is user's paper mongodb implementation */
type Paper struct {
	// mongodb's object id
	PaperId bson.ObjectId	`bson:"_id"`
	// paper owner's username
	Owner	string
	// paper saved path
	Path	string
}

type iPaper interface {
	AddPaper(paper Paper) (error)
	SelectPaper(filter map[string]interface{}) ([]Paper)
}

type MongodbPaper struct {
	database 	string
	// collection name, i.e. table name
	collection  string
}

func (self MongodbPaper) AddPaper(paper Paper) (error) {
	session := getMongodbSession()
	defer session.Close()

	c := session.DB(self.database).C(self.collection)
	err := c.Insert(&paper)		// It can insert not only one Object, it acquires "...&interface"
	if err != nil {
		return err;
	}
	return nil;
}

func (self MongodbPaper) SelectPaper(filter map[string]interface{}) ([]Paper) {
	session := getMongodbSession()
	defer session.Close()
	c := session.DB(self.database).C(self.collection)
	res := []Paper{}
	err := c.Find(&filter).All(&res)
	if err != nil {
		Logger.Fatalf("Find Paper Error, %s\n", err)
	}
	return res
}

/**
If
*/

type SqlUser struct {

}


