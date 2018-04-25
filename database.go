package Posger

const (
	DATABASE = "mongodb"		// designated the using database
	MONGO_DB_SEVER = "mongodb://localhost:27017"	// mongodb's running port.
)

func AddUser(user User) error {
	return GetDatabaseConncect(DATABASE).AddUser(user)
}

func SelectUser(username string) User {
	return GetDatabaseConncect(DATABASE).SelectUser(username)
}


// TODO: It mantains definitive database connection session.
type ConnPool interface {
	GetConnection()	interface{}
}

func getConnection () {
	return
}