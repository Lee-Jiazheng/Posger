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
		Logger.Println("Add user failed, %s\n", user.Username)
	}
}

// Filter condition, e.g.: map[string]interface{}{"username": username, "source": "github"}
func SelectUser(filter map[string]interface{}) []User {
	return GetDatabaseConncect(DATABASE).SelectUser(filter)
}

func AddPaper(paper Paper) {
	//paper.PaperId = bson.NewObjectId()
	if paper.Owner != "anonymous" && len(SelectUser(map[string]interface{}{"userid": paper.Owner})) == 0 {
		Logger.Println("the username didn't exist, " + paper.Owner)
		return
	}
	if err := GetDatabaseConncect(DATABASE).AddPaper(paper); err != nil {
		Logger.Println("Add paper failed, %s, %s\n", paper, err)
	}
}

func SelectPaper(filter map[string]interface{}) []Paper{
	return GetDatabaseConncect(DATABASE).SelectPaper(filter)
}

func DeletePaper(filter map[string]interface{}) {
	if len(SelectPaper(filter)) == 0 {
		Logger.Println("the paper filter condition found none paper, ", filter)
		return
	} else {
		GetDatabaseConncect(DATABASE).DeletePaper(filter)
	}
}

func AddQuestion(question Question) {
	//paper.PaperId = bson.NewObjectId()
	if len(SelectQuestion(map[string]interface{}{"question": question.Question})) != 0 {
		Logger.Println("the question has existed in database, ", question)
		return
	}
	if err := GetDatabaseConncect(DATABASE).AddQuestion(question); err != nil {
		Logger.Println("Add question failed, %s, %s\n", question, err)
	}
}

func SelectQuestion(filter map[string]interface{}) []Question{
	return GetDatabaseConncect(DATABASE).SelectQuestion(filter)
}

func SetQuestionAnswer(question Question) {
	if len(SelectQuestion(map[string]interface{}{"questionid": question.QuestionId})) == 0 {
		Logger.Println("the question didn't exist, " + question.QuestionId)
		return
	}
	if err := GetDatabaseConncect(DATABASE).SetQuestionAnswer(question); err != nil {
		Logger.Println("Update question exception, ", question)
	}

}

// TODO: It mantains definitive database connection session.
type ConnPool interface {
	GetConnection()	interface{}
}

func getConnection () {
	return
}