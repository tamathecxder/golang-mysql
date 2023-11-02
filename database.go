package golang_mysql

import (
	"database/sql"
	"time"
)

// SetMaxIdleConns sets the maximum number of idle connections in the connection pool.

// SetMaxOpenConns sets the maximum number of open (in-use and idle) connections in the pool.

// SetConnMaxIdleTime sets the maximum amount of time a connection can remain idle before it's closed and removed from the pool.

// SetConnMaxLifetime sets the maximum lifespan of a connection. After this duration, a connection is closed and removed from the pool.

func GetConnection() *sql.DB {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/golang_test?parseTime=true")

	if err != nil {
		panic(err)
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(60 * time.Minute)

	return db
}
