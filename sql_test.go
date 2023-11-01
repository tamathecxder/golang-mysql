package golang_mysql

import (
	"context"
	"fmt"
	"testing"
)

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

		fmt.Println("Id:", id)
		fmt.Println("Name:", name)
	}
}
