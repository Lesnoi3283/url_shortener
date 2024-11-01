package logic

// PingDB just pings a database and returns result.
func PingDB(db URLStorageInterface) error {
	return db.Ping()
}
