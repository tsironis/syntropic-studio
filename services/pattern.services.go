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
	ID            int       `json:"id"`
	CreatedBy     int       `json:"created_by"`
	Title         string    `json:"title"`
	TotalDistance int       `json:"total_distance,omitempty"`
	DistanceStep  float64   `json:"step_distance,omitempty"`
	RowData       float64   `json:"row_data,omitempty"`
	Description   string    `json:"description,omitempty"`
	Private       bool      `json:"private,omitempty"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
}

type PatternServices struct {
	Pattern      Pattern
	PatternStore db.Store
}

func (ts *PatternServices) CreatePattern(t Pattern) (Pattern, error) {
	result := ts.PatternStore.Db.Create(&t)
	if result.Error != nil {
		return Pattern{}, result.Error
	}

	fmt.Printf("Pattern created: %d\n", t)

	ts.Pattern = t

	return ts.Pattern, nil
}

func (ts *PatternServices) GetAllPatterns(createdBy int) ([]Pattern, error) {
	patterns := []Pattern{}

	result := ts.PatternStore.Db.Where("created_by = ?", createdBy).Order("created_at desc").Find(&patterns)
	if result.Error != nil {
		return []Pattern{}, result.Error
	}
	return patterns, nil
}

func (ts *PatternServices) GetPatternById(t Pattern) (Pattern, error) {
	result := ts.PatternStore.Db.Where(&t, "created_by", "id").Find(&ts.Pattern)
	if result.Error != nil {
		return Pattern{}, result.Error
	}

	return ts.Pattern, nil
}

func (ts *PatternServices) UpdatePattern(t Pattern) (Pattern, error) {
	result := ts.PatternStore.Db.Model(&ts.Pattern).Where(&t, "created_by", "id").Updates(&t)
	if result.Error != nil {
		return Pattern{}, result.Error
	}

	return ts.Pattern, nil
}

func (ts *PatternServices) DeletePattern(t Pattern) error {
	result := ts.PatternStore.Db.Delete(&t)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 1 {
		return errors.New("an affected row was expected")
	}

	return nil
}
