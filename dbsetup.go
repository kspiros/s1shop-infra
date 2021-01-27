package infra

//SetupDB setup database dependecies
func SetupDB(dbconn IDBConn) error {
	//create indexes
	return dbconn.CreateUniqueIndex("user", "email")

}
