package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jailtonjunior94/go-sql/internal/db"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	ctx := context.Background()
	dbConnection, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/courses")
	if err != nil {
		panic(err)
	}
	defer dbConnection.Close()

	category := CategoryParams{
		ID:          uuid.New().String(),
		Name:        "Backend",
		Description: sql.NullString{String: "Backend Course", Valid: true},
	}

	course := CourseParams{
		ID:          uuid.New().String(),
		Name:        "Go",
		Price:       10.94,
		Description: sql.NullString{String: "Go Course", Valid: true},
	}

	courseDB := NewCourseDB(dbConnection)
	err = courseDB.CreateCourseAndCategory(ctx, category, course)
	if err != nil {
		panic(err)
	}
}

type (
	CourseDB struct {
		dbConn *sql.DB
		*db.Queries
	}

	CourseParams struct {
		ID          string
		Name        string
		Price       float64
		Description sql.NullString
	}

	CategoryParams struct {
		ID          string
		Name        string
		Description sql.NullString
	}
)

func NewCourseDB(dbConn *sql.DB) *CourseDB {
	return &CourseDB{
		dbConn:  dbConn,
		Queries: db.New(dbConn),
	}
}

func (c *CourseDB) callTx(ctx context.Context, fn func(*db.Queries) error) error {
	tx, err := c.dbConn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := db.New(tx)
	err = fn(q)
	if err != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			return fmt.Errorf("error on rollback %v, original error %w", errRollback, err)
		}
		return err
	}
	return tx.Commit()
}

func (c *CourseDB) CreateCourseAndCategory(ctx context.Context, category CategoryParams, course CourseParams) error {
	return c.callTx(ctx, func(q *db.Queries) error {
		if err := q.CreateCategory(ctx, db.CreateCategoryParams{
			ID:          category.ID,
			Name:        category.Name,
			Description: category.Description,
		}); err != nil {
			return err
		}

		if err := q.CreateCourse(ctx, db.CreateCourseParams{
			ID:          course.ID,
			Name:        course.Name,
			Description: course.Description,
			Price:       course.Price,
			CategoryID:  category.ID,
		}); err != nil {
			return err
		}
		return nil
	})
}
