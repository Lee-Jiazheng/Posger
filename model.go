// model.go
// Model is denoted the database designation and
// Databaser's interface
// Warning: All interface function should not consider the exception condition
// the upper function will handler the exception

package Posger

import (
	"gopkg.in/mgo.v2"
)


/** Databaser is An abstract factory interface, if you want to add a new database backend, you must implement all function*/
type Databaser interface {
	iUser
	iPaper
	iQuestion
}

func GetDatabaseConncect(db string) (Databaser) {
	if db == "mongodb" { return newMongodb()}
	Logger.Println("database is not the validate database backend.")
	return nil
}

func getMongodbSession() (*mgo.Session){
	session, err := mgo.Dial(MONGO_DB_SEVER)
	if err != nil {
		Logger.Println(err)
	}
	session.SetSafe(&mgo.Safe{})
	return session
}

type MongodbSetting struct {
	database	string
	collection	string
}

// Mongodb's construction, also construct it's interface implementation struct, i.e. subclass.
func newMongodb() (*Mongodb){
	var mdb = new(Mongodb)
	// construct All databases name and collections name, need to be explicit.
	mdb.MongodbUser.database = "Mongodb"
	mdb.MongodbUser.collection = "Users"	// table name.
	mdb.MongodbPaper.database = "Mongodb"
	mdb.MongodbPaper.collection = "Papers"
	mdb.MongodbQuestion.database = "Mongodb"
	mdb.MongodbQuestion.collection = "Questions"
	return mdb
}

type mongoSetting struct {
	database	string
	// collection name, i.e. table name
	collection 	string
}
// anonymous member, we can directly point the final function, such as AddUser.
type Mongodb struct{
	MongodbUser
	MongodbPaper
	MongodbQuestion
}

/** The models definition and corresponding interface. */
type User struct {
	// user's id
	UserId		string
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
	MongodbSetting
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
	//PaperId bson.ObjectId	`bson:"_id"`
	// a unique id, supposed to be a uuid
	PaperId	string
	// paper owner's username
	Owner	string
	// paper saved path
	Path	string
	// paper upload origin name
	Name	string
	// paper's built-in images extracted path
	Images	[]string
	// Create TimeStamp, second precision
	C_Time	int32
}

type iPaper interface {
	AddPaper(paper Paper) (error)
	SelectPaper(filter map[string]interface{}) ([]Paper)
	DeletePaper(filter map[string]interface{})
}

type MongodbPaper struct {
	MongodbSetting
}

func (self MongodbPaper) AddPaper(paper Paper) (error) {
	session := getMongodbSession()
	defer session.Close()

	c := session.DB(self.database).C(self.collection)
	return c.Insert(&paper)		// It can insert not only one Object, it acquires "...&interface"
}

func (self MongodbPaper) SelectPaper(filter map[string]interface{}) ([]Paper) {
	session := getMongodbSession()
	defer session.Close()
	c := session.DB(self.database).C(self.collection)
	res := []Paper{}
	err := c.Find(&filter).All(&res)
	if err != nil {
		Logger.Println("Find Paper Error, %s\n", err)
	}
	return res
}

func (self MongodbPaper) DeletePaper(filter map[string]interface{}) {
	session := getMongodbSession()
	defer session.Close()
	c := session.DB(self.database).C(self.collection)
	err := c.Remove(filter)
	if err != nil {
		Logger.Println("Find Paper Error, %s\n", err)
	}
}

// Question's model prototype
type Question struct {
	// Also UUID generated
	QuestionId	string
	// Question Content.
	Question	string
	// Answer generated by R-net
	Answer		string
	// Corresponding the passages
	Scores		[]float32
	// Extracted Answer by passages
	Answers		[]string
	// System reference passages
	Passages	[]string
}

type iQuestion interface {
	// First Add question, shall be no answer
	AddQuestion(question Question) (error)
	// filter question
	SelectQuestion(filter map[string]interface{}) ([]Question)
	// Set Question's answer and reference passages, Warning: the question must exist in the system
	SetQuestionAnswer(question Question) (error)
}

type MongodbQuestion struct {
	MongodbSetting
}

func (self MongodbQuestion) AddQuestion(question Question) (error) {
	sess := getMongodbSession()
	defer sess.Close()

	c := sess.DB(self.database).C(self.collection)
	return c.Insert(question)
}

func (self MongodbQuestion) SelectQuestion(filter map[string]interface{}) ([]Question) {
	sess := getMongodbSession()
	defer sess.Close()
	c := sess.DB(self.database).C(self.collection)
	res := []Question{}
	if err := c.Find(&filter).All(&res); err != nil {
		Logger.Println("Find Question Error, %s, the condition is %s\n", err, filter)
	}
	return res
}

func (self MongodbQuestion) SetQuestionAnswer(question Question) (error) {
	sess := getMongodbSession()
	defer sess.Close()
	c := sess.DB(self.database).C(self.collection)
	return c.Update(map[string]interface{}{"questionid": question.QuestionId}, question)
}


/**
If
*/

type SqlUser struct {

}


