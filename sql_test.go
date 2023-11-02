package golang_mysql

import (
	"context"
	"database/sql"
	"fmt"
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
			name VARCHAR(255),
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
		name     string
		password string
	}{
		{GenerateCustomUUID(), "User1", "password1"},
		{GenerateCustomUUID(), "User2", "password2"},
		{GenerateCustomUUID(), "User3", "password3"},
	}

	for _, userData := range usersData {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userData.password), bcrypt.DefaultCost)
		if err != nil {
			t.Errorf("failed to hash password: %v", err)
			continue
		}

		insertQuery := `
			INSERT INTO user (id, name, password) VALUES (?, ?, ?)
		`

		_, err = db.ExecContext(ctx, insertQuery, userData.id, userData.name, hashedPassword)
		if err != nil {
			t.Errorf("failed to insert user:: %v", err)
		}
	}

	fmt.Println("data has been saved successfully")
}
