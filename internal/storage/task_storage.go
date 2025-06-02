package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/edkuzhakhmetov/go_final_project/internal/models"
)

const (
	insTask  = `INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat);`
	updTask  = `Update scheduler set date=?, title=?, comment=?, repeat=? where id=?;`
	delTask  = `delete from scheduler where id=?;`
	updDate  = `Update scheduler set date=? where id=?;`
	getTasks = `SELECT id, date, title, comment, repeat FROM scheduler order by date limit :limit;`
	getTask  = `SELECT id, date, title, comment, repeat FROM scheduler where id=:id;`
)

func executeStatement(ctx context.Context, db *sql.DB, query string, args ...interface{}) (int64, error) {
	stmt, err := db.Prepare(query)
	if err != nil {
		return 0, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	fmt.Println("args", args)
	res, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to execute statement: %w", err)
	}

	return res.RowsAffected()
}

func (s *Storage) AddTask(ctx context.Context, task models.Task) (int, error) {

	stmt, err := s.db.Prepare(insTask)
	if err != nil {
		return 0, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
	)
	if err != nil {
		return 0, fmt.Errorf("failed to execute statement: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}
	return int(id), nil
}

func (s *Storage) UpdateTask(ctx context.Context, task models.Task) error {
	rowsAffected, err := executeStatement(ctx, s.db, updTask,
		task.Date, task.Title, task.Comment, task.Repeat, task.ID)

	if err != nil {
		return fmt.Errorf("failed to execute statement: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("task with id %s not found", task.ID)
	}
	return nil

}

func (s *Storage) UpdateDate(ctx context.Context, next string, id int) error {

	rowsAffected, err := executeStatement(ctx, s.db, updDate,
		next, id)

	if err != nil {
		return fmt.Errorf("failed to execute statement: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("task with id %d not found", id)
	}
	return nil
}

func (s *Storage) DeleteTask(ctx context.Context, id int) error {
	rowsAffected, err := executeStatement(ctx, s.db, delTask,
		id)

	if err != nil {
		return fmt.Errorf("failed to execute statement: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("task with id %d not found", id)
	}
	return nil
}

func (s *Storage) GetTasks(ctx context.Context, limit int) ([]*models.Task, error) {

	var res []*models.Task
	stmt, err := s.db.Prepare(getTasks)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(sql.Named("limit", limit))

	if err != nil && err == sql.ErrNoRows {
		return nil, err

	} else if err != nil {
		return nil, fmt.Errorf("an error occurred. %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var task models.Task
		err = rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		res = append(res, &task)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("an error occurred. %w", err)
	}
	return res, nil

}

func (s *Storage) GetTask(ctx context.Context, id int) (models.Task, error) {

	var res models.Task
	stmt, err := s.db.Prepare(getTask)
	if err != nil {
		return models.Task{}, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRow(sql.Named("id", id))
	err = row.Scan(&res.ID, &res.Date, &res.Title, &res.Comment, &res.Repeat)
	if err != nil && err == sql.ErrNoRows {
		return models.Task{}, err

	} else if err != nil {
		return models.Task{}, fmt.Errorf("an error occurred. %w", err)
	}
	return res, nil

}
