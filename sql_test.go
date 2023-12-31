package golang_mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func TestCreateTable(t *testing.T) {
	db := GetConnection()
	ctx := context.Background()

	defer db.Close()

	createTableQuery := `
		CREATE TABLE user (
			id CHAR(36) PRIMARY KEY,
			username VARCHAR(255),
			password VARCHAR(100)
		)
	`

	_, err := db.ExecContext(ctx, createTableQuery)

	if err != nil {
		panic(err)
	}

	fmt.Println("table has been successfully created")
}

func TestExecSql(t *testing.T) {
	db := GetConnection()
	ctx := context.Background()

	defer db.Close()

	query := "INSERT INTO customer(id, name) VALUES ('1', 'goyoonjung')"
	_, err := db.ExecContext(ctx, query)

	if err != nil {
		panic(err)
	}

	fmt.Println("Success: Customer data inserted")
}

func TestQuerySql(t *testing.T) {
	db := GetConnection()
	ctx := context.Background()

	defer db.Close()

	query := "SELECT * FROM customer"
	rows, err := db.QueryContext(ctx, query)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	fmt.Println("Success: Query executed. Customer data retrieved.")

	for rows.Next() {
		var id, name string

		err := rows.Scan(&id, &name)

		if err != nil {
			panic(err)
		}

		fmt.Println("========================")
		fmt.Println("Id:", id)
		fmt.Println("Name:", name)
	}
}

func TestComplexQuery(t *testing.T) {
	db := GetConnection()
	ctx := context.Background()

	defer db.Close()

	query := "SELECT id, name, email, balance, rating, birth_date, married, created_at FROM customer"
	rows, err := db.QueryContext(ctx, query)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	fmt.Println("Success: Query executed. Customer data retrieved.")

	for rows.Next() {
		var id, name, email string
		var balance int32
		var rating float64
		var married bool
		var birthDate, createdAt time.Time

		err := rows.Scan(&id, &name, &email, &balance, &rating, &birthDate, &married, &createdAt)

		if err != nil {
			panic(err)
		}

		fmt.Println("========================")
		fmt.Println("Id:", id)
		fmt.Println("Name:", name)
		fmt.Println("Email:", email)
		fmt.Println("Balance:", balance)
		fmt.Println("Rating:", rating)
		fmt.Println("Married:", married)
		fmt.Println("Birth Date:", birthDate)
		fmt.Println("Created At:", createdAt)
	}
}

func TestNullValue(t *testing.T) {
	db := GetConnection()
	ctx := context.Background()

	defer db.Close()

	query := "SELECT id, name, email, balance, rating, birth_date, married, created_at FROM customer"
	rows, err := db.QueryContext(ctx, query)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	fmt.Println("Success: Query executed. Customer data retrieved.")

	for rows.Next() {
		var id, name string
		var email sql.NullString
		var balance int32
		var rating float64
		var married bool
		var birthDate sql.NullTime
		var createdAt time.Time

		err := rows.Scan(&id, &name, &email, &balance, &rating, &birthDate, &married, &createdAt)

		if err != nil {
			panic(err)
		}

		fmt.Println("========================")
		fmt.Println("Id:", id)
		fmt.Println("Name:", name)

		if email.Valid {
			fmt.Println("Email:", email.String)
		} else {
			fmt.Println("Email:", nil)
		}

		fmt.Println("Balance:", balance)
		fmt.Println("Rating:", rating)
		fmt.Println("Married:", married)

		if birthDate.Valid {
			fmt.Println("Birth Date:", birthDate.Time)
		} else {
			fmt.Println("Birth Date:", nil)
		}

		fmt.Println("Created At:", createdAt)
	}
}

func TestUserInsert(t *testing.T) {
	db := GetConnection()
	ctx := context.Background()
	defer db.Close()

	usersData := []struct {
		id       string
		username string
		password string
	}{
		{GenerateCustomUUID(), "User1", "password1"},
		{GenerateCustomUUID(), "User2", "password2"},
		{GenerateCustomUUID(), "User3", "password3"},
	}

	for _, userData := range usersData {
		var hashedPassword []byte
		var err error

		if userData.username == "User1" {
			hashedPassword = []byte(userData.password)
		} else {
			hashedPassword, err = bcrypt.GenerateFromPassword([]byte(userData.password), bcrypt.DefaultCost)
			if err != nil {
				t.Errorf("failed to hash password: %v", err)
				continue
			}
		}

		insertQuery := `
			INSERT INTO user (id, username, password) VALUES (?, ?, ?)
		`

		_, err = db.ExecContext(ctx, insertQuery, userData.id, userData.username, hashedPassword)
		if err != nil {
			t.Errorf("failed to insert user: %v", err)
		}
	}

	fmt.Println("data has been saved successfully")
}

func TestSqlInjection(t *testing.T) {
	db := GetConnection()
	ctx := context.Background()

	defer db.Close()

	username := "User1"
	password := "password1'; #"

	query := "SELECT id, username, password FROM user WHERE username = '" + username + "'AND password = '" + password + "'LIMIT 1"

	rows, err := db.QueryContext(ctx, query)

	if err != nil {
		panic(err)
	}

	if rows.Next() {
		var id, resultUsername, resultPassword string

		err := rows.Scan(&id, &resultUsername, &resultPassword)

		if err != nil {
			panic(err)
		}

		fmt.Println("Login berhasil")
		fmt.Println("Id:", id)
		fmt.Println("Username:", resultUsername)
	} else {
		fmt.Println("Login gagal")
	}
}

func TestSqlInjectionSafe(t *testing.T) {
	db := GetConnection()
	ctx := context.Background()

	defer db.Close()

	username := "User1"
	password := "password1'; #"

	query := "SELECT id, username, password FROM user WHERE username = ? AND password = ? LIMIT 1"

	rows, err := db.QueryContext(ctx, query, username, password)

	if err != nil {
		panic(err)
	}

	if rows.Next() {
		var id, resultUsername, resultPassword string

		err := rows.Scan(&id, &resultUsername, &resultPassword)

		if err != nil {
			panic(err)
		}

		fmt.Println("Login berhasil")
		fmt.Println("Id:", id)
		fmt.Println("Username:", resultUsername)
	} else {
		fmt.Println("Login gagal")
	}
}

func TestExecSqlSafe(t *testing.T) {
	db := GetConnection()
	ctx := context.Background()

	defer db.Close()

	id := "1"
	username := "super_admin"
	password := "admin123"

	query := "INSERT INTO user(id, username, password) VALUES (?, ?, ?)"
	_, err := db.ExecContext(ctx, query, id, username, password)

	if err != nil {
		panic(err)
	}

	fmt.Println("Success: User data inserted")
}

func TestAutoIncrement(t *testing.T) {
	db := GetConnection()
	ctx := context.Background()

	defer db.Close()

	username := "test@gmail.com"
	password := "test comment"

	query := "INSERT INTO comments (email, comment) VALUES (?, ?)"
	result, err := db.ExecContext(ctx, query, username, password)

	if err != nil {
		panic(err)
	}

	insertId, err := result.LastInsertId()

	if err != nil {
		panic(err)
	}

	fmt.Println("Success: Comment data inserted with id", insertId)
}

func TestPrepareStatement(t *testing.T) {
	db := GetConnection()
	ctx := context.Background()

	defer db.Close()

	insertQuery := "INSERT INTO comments (email, comment) VALUES (?, ?)"
	statement, err := db.PrepareContext(ctx, insertQuery)

	if err != nil {
		panic(err)
	}

	for i := 0; i < 10; i++ {
		username := "testing.mail" + strconv.Itoa(i) + "@gmail.com"
		comment := "TESTING_COMMENT_NUMBER_" + strconv.Itoa(i)

		result, err := statement.ExecContext(ctx, username, comment)

		if err != nil {
			panic(err)
		}

		lastId, err := result.LastInsertId()

		if err != nil {
			panic(err)
		}

		fmt.Println("Success: Comment data inserted with id", lastId)
	}

	defer statement.Close()
}

func TestTransaction(t *testing.T) {
	db := GetConnection()

	defer db.Close()

	ctx := context.Background()

	tx, err := db.Begin()

	if err != nil {
		panic(err)
	}

	// do transaction
	insertQuery := "INSERT INTO comments (email, comment) VALUES (?, ?)"

	for i := 0; i < 10; i++ {
		username := "trx.testing.mail" + strconv.Itoa(i) + "@gmail.com"
		comment := "DB_TRANSACTION_COMMENT_" + strconv.Itoa(i)

		result, err := tx.ExecContext(ctx, insertQuery, username, comment)

		if err != nil {
			panic(err)
		}

		lastId, err := result.LastInsertId()

		if err != nil {
			panic(err)
		}

		fmt.Println("Success: Comment data inserted with id", lastId)
	}

	err = tx.Commit()

	if err != nil {
		panic(err)
	}
}
