package Posger


const (
	DATABASE = "mongodb"		// designated the using database
	MONGO_DB_SEVER = "mongodb://localhost:27017"	// mongodb's running port.
)

func AddUser(user User) {
	if len(SelectUser(map[string]interface{}{"userid": user.UserId})) != 0 {
		Logger.Printf("the userId has existed database, %s\n", user)
		return
	}
	if err := GetDatabaseConncect(DATABASE).AddUser(user); err != nil {
		Logger.Fatalf("Add user failed, %s\n", user.Username)
	}
}

// Filter condition, e.g.: map[string]interface{}{"username": username, "source": "github"}
func SelectUser(filter map[string]interface{}) []User {
	return GetDatabaseConncect(DATABASE).SelectUser(filter)
}

func AddPaper(paper Paper) {
	//paper.PaperId = bson.NewObjectId()
	if len(SelectUser(map[string]interface{}{"username": paper.Owner})) == 0 {
		Logger.Fatalln("the username didn't exist, " + paper.Owner)
		return
	}
	if err := GetDatabaseConncect(DATABASE).AddPaper(paper); err != nil {
		Logger.Fatalf("Add paper failed, %s, %s\n", paper, err)
	}
}

func SelectPaper(filter map[string]interface{}) []Paper{
	return GetDatabaseConncect(DATABASE).SelectPaper(filter)
}

func DeletePaper(filter map[string]interface{}) {
	if len(SelectPaper(filter)) == 0 {
		Logger.Fatalln("the paper filter condition found none paper, ", filter)
		return
	} else {
		GetDatabaseConncect(DATABASE).DeletePaper(filter)
	}
}


// TODO: It mantains definitive database connection session.
type ConnPool interface {
	GetConnection()	interface{}
}

func getConnection () {
	return
}