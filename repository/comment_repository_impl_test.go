package repository

import (
	"context"
	"fmt"
	"golang_mysql"
	"golang_mysql/entity"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestCommentInsert(t *testing.T) {
	commentRepository := NewCommentRepository(golang_mysql.GetConnection())

	ctx := context.Background()

	comment := entity.Comment{
		Email:   "repo@gmail.com",
		Comment: "Hello from repository",
	}

	result, err := commentRepository.Insert(ctx, comment)

	if err != nil {
		panic(err)
	}

	fmt.Println(result)
}

func TestFindById(t *testing.T) {
	commentRepository := NewCommentRepository(golang_mysql.GetConnection())

	ctx := context.Background()

	comment, err := commentRepository.FindById(ctx, 13)

	if err != nil {
		panic(err)
	}

	fmt.Println(comment)
}
