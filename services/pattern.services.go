package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/tsironis/syntropic-studio/db"
)

func NewPatternServices(t Pattern, tStore db.Store) *PatternServices {

	return &PatternServices{
		Pattern:      t,
		PatternStore: tStore,
	}
}

type Pattern struct {
	ID          int       `json:"id"`
	CreatedBy   int       `json:"created_by"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Status      bool      `json:"status,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
}

type PatternServices struct {
	Pattern      Pattern
	PatternStore db.Store
}

func (ts *PatternServices) CreatePattern(t Pattern) (Pattern, error) {

	query := `INSERT INTO patterns (created_by, title, description)
		VALUES(?, ?, ?) RETURNING *`

	stmt, err := ts.PatternStore.Db.Prepare(query)
	if err != nil {
		return Pattern{}, err
	}

	defer stmt.Close()

	err = stmt.QueryRow(
		t.CreatedBy,
		t.Title,
		t.Description,
	).Scan(
		&ts.Pattern.ID,
		&ts.Pattern.CreatedBy,
		&ts.Pattern.Title,
		&ts.Pattern.Description,
		&ts.Pattern.Status,
		&ts.Pattern.CreatedAt,
	)
	if err != nil {
		return Pattern{}, err
	}

	/* if i, err := result.RowsAffected(); err != nil || i != 1 {
		return errors.New("error: an affected row was expected")
	} */

	return ts.Pattern, nil
}

func (ts *PatternServices) GetAllPatterns(createdBy int) ([]Pattern, error) {
	query := fmt.Sprintf("SELECT id, title, status FROM patterns WHERE created_by = %d ORDER BY created_at DESC", createdBy)

	rows, err := ts.PatternStore.Db.Query(query)
	if err != nil {
		return []Pattern{}, err
	}
	// We close the resource
	defer rows.Close()

	patterns := []Pattern{}
	for rows.Next() {
		rows.Scan(&ts.Pattern.ID, &ts.Pattern.Title, &ts.Pattern.Status)

		patterns = append(patterns, ts.Pattern)
	}

	return patterns, nil
}

func (ts *PatternServices) GetPatternById(t Pattern) (Pattern, error) {

	query := `SELECT id, title, description, status, created_at FROM patterns
		WHERE created_by = ? AND id=?`

	stmt, err := ts.PatternStore.Db.Prepare(query)
	if err != nil {
		return Pattern{}, err
	}

	defer stmt.Close()

	err = stmt.QueryRow(
		t.CreatedBy,
		t.ID,
	).Scan(
		&ts.Pattern.ID,
		&ts.Pattern.Title,
		&ts.Pattern.Description,
		&ts.Pattern.Status,
		&ts.Pattern.CreatedAt,
	)
	if err != nil {
		return Pattern{}, err
	}

	return ts.Pattern, nil
}

func (ts *PatternServices) UpdatePattern(t Pattern) (Pattern, error) {

	query := `UPDATE patterns SET title = ?,  description = ?, status = ?
		WHERE created_by = ? AND id=? RETURNING id, title, description, status`

	stmt, err := ts.PatternStore.Db.Prepare(query)
	if err != nil {
		return Pattern{}, err
	}

	defer stmt.Close()

	err = stmt.QueryRow(
		t.Title,
		t.Description,
		t.Status,
		t.CreatedBy,
		t.ID,
	).Scan(
		&ts.Pattern.ID,
		&ts.Pattern.Title,
		&ts.Pattern.Description,
		&ts.Pattern.Status,
	)
	if err != nil {
		return Pattern{}, err
	}

	return ts.Pattern, nil
}

func (ts *PatternServices) DeletePattern(t Pattern) error {

	query := `DELETE FROM patterns
		WHERE created_by = ? AND id=?`

	stmt, err := ts.PatternStore.Db.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	result, err := stmt.Exec(t.CreatedBy, t.ID)
	if err != nil {
		return err
	}

	if i, err := result.RowsAffected(); err != nil || i != 1 {
		return errors.New("an affected row was expected")
	}

	return nil
}
