package pgsql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Emibotz/megu/internal/todo"
	"github.com/jackc/pgx/v5/pgxpool"
)

type todoStore struct {
	db *pgxpool.Pool
}

func (d *db) TODOStore() (*todoStore, error) {
	return &todoStore{db: d.pool}, nil
}

func (s *todoStore) Create(ctx context.Context, td *todo.TODO) error {
	query := `
INSERT INTO
	todos
	(title, content, completed, created_at, deleted_at)
VALUES
	($1, $2, $3, $4, $5)
;
	`

	_, err := s.db.Exec(ctx, query, td.Title, td.Content, td.Completed, td.CreatedAt, td.DeletedAt)
	if err != nil {
		return err
	}

	return nil
}

func (s *todoStore) List(ctx context.Context, limit int, offset int) ([]todo.ListResp, error) {
	sql := `
SELECT
	ROW_NUMBER() OVER (ORDER BY created_at) AS num,
	id, title, content, completed, created_at, deleted_at
FROM
	todos
LIMIT
	$1
OFFSET
	$2
;
	`

	rows, err := s.db.Query(ctx, sql, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []todo.ListResp = nil
	for rows.Next() {
		var resp todo.ListResp

		if err := rows.Scan(
			&resp.Num,
			&resp.TD.ID, &resp.TD.Title, &resp.TD.Content, &resp.TD.Completed, &resp.TD.CreatedAt, &resp.TD.DeletedAt,
		); err != nil {
			return nil, err
		}

		result = append(result, resp)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// 获取指定序号的待办，若没有指定记录则返回 (nil, nil)
func (s *todoStore) GetByNumber(ctx context.Context, num int) (*todo.TODO, error) {
	query := `
SELECT
	id, title, content, completed, created_at, deleted_at
FROM (
	SELECT
		ROW_NUMBER() OVER (ORDER BY created_at) AS num,
		id, title, content, completed, created_at, deleted_at
	FROM
		todos
) AS td
WHERE
	td.num = $1
;
	`

	row := s.db.QueryRow(ctx, query, num)

	var td todo.TODO
	if err := row.Scan(
		&td.ID, &td.Title, &td.Content, &td.Completed, &td.CreatedAt, &td.DeletedAt,
	); err != nil {

		// 若没有指定记录则返回 (nil, nil)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &td, nil
}

func (s *todoStore) Update(ctx context.Context, td *todo.TODO) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		tx.Rollback(ctx)
	}()

	sql := `
UPDATE
	todos
SET
	title = $2, content = $3, completed = $4, created_at = $5, deleted_at = $6
WHERE
	id = $1
;
	`

	cmdTag, err := s.db.Exec(ctx, sql, td.ID, td.Title, td.Content, td.Completed, td.CreatedAt, td.DeletedAt)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() != 1 {
		return fmt.Errorf("expected 1 row affected, got %d", cmdTag.RowsAffected())
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (s *todoStore) Delete(ctx context.Context, td *todo.TODO) error {
	sql := `
DELETE FROM
	todos
WHERE
	id = $1
;
	`

	cmdTag, err := s.db.Exec(ctx, sql, td.ID)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() != 1 {
		return fmt.Errorf("how could this happen? expected 1 row affected, got %d", cmdTag.RowsAffected())
	}

	return nil
}
