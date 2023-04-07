package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jailtonjunior94/go-sql/internal/db"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

func main() {
	ctx := context.Background()
	dbConnection, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/courses")
	if err != nil {
		panic(err)
	}
	defer dbConnection.Close()

	queries := db.New(dbConnection)

	if err = queries.CreateCategory(ctx, db.CreateCategoryParams{
		ID:          uuid.New().String(),
		Name:        "Backend",
		Description: sql.NullString{String: "Backend Description", Valid: true},
	}); err != nil {
		panic(err)
	}

	categories, err := queries.ListCategories(ctx)
	if err != nil {
		panic(err)
	}

	for _, category := range categories {
		fmt.Println(category.ID, category.Name, category.Description)
	}
}
