package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/tsironis/syntropic-studio/db"
)

func NewDesignServices(t Design, tStore db.Store) *DesignServices {

	return &DesignServices{
		Design:      t,
		DesignStore: tStore,
	}
}

type Design struct {
	ID        int       `json:"id"`
	CreatedBy int       `json:"created_by"`
	Title     string    `json:"title"`
	QGISData  string    `json:"qgis_data"`
	Status    bool      `json:"status,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type DesignServices struct {
	Design      Design
	DesignStore db.Store
}

func (ts *DesignServices) CreateDesign(t Design) (Design, error) {
	result := ts.DesignStore.Db.Create(&t)
	if result.Error != nil {
		return Design{}, result.Error
	}

	fmt.Printf("Design created: %d\n", t)

	return ts.Design, nil
}

func (ts *DesignServices) GetAllDesigns(createdBy int) ([]Design, error) {
	designs := []Design{}

	result := ts.DesignStore.Db.Where("created_by = ?", createdBy).Order("created_at desc").Find(&designs)
	if result.Error != nil {
		return []Design{}, result.Error
	}
	return designs, nil
}

func (ts *DesignServices) GetDesignById(t Design) (Design, error) {
	result := ts.DesignStore.Db.Where(&t, "created_by", "id").Find(&ts.Design)
	if result.Error != nil {
		return Design{}, result.Error
	}

	return ts.Design, nil
}

func (ts *DesignServices) UpdateDesign(t Design) (Design, error) {
	result := ts.DesignStore.Db.Model(&ts.Design).Where(&t, "created_by", "id").Updates(&t)
	if result.Error != nil {
		return Design{}, result.Error
	}

	return ts.Design, nil
}

func (ts *DesignServices) DeleteDesign(t Design) error {
	result := ts.DesignStore.Db.Delete(&t)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 1 {
		return errors.New("an affected row was expected")
	}

	return nil
}
