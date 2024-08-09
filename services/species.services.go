package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/tsironis/syntropic-studio/db"
)

func NewSpeciesServices(t Species, tStore db.Store) *SpeciesServices {

	return &SpeciesServices{
		Species:      t,
		SpeciesStore: tStore,
	}
}

type Species struct {
	ID                   int       `json:"id"`
	CreatedBy            int       `json:"created_by"`
	BinomialNomenclature string    `json:"binomial_nomenclature"`
	CommonName           string    `json:"common_name"`
	Stratum              string    `json:"stratum]"`
	System               string    `json:"system"`
	Lifecycle            string    `json:"lifecycle"`
	BiologicalForm       string    `json:"biological_form"`
	USDAHardinessZone    string    `json:"usda_hardiness_zone"`
	PrimaryFunction      string    `json:"primary_function"`
	SecondaryFunction    string    `json:"secondary_function"`
	Notes                string    `json:"Notes"`
	Private              bool      `json:"private,omitempty"`
	CreatedAt            time.Time `json:"created_at,omitempty"`
}

type SpeciesServices struct {
	Species      Species
	SpeciesStore db.Store
}

func (ts *SpeciesServices) CreateSpecies(t Species) (Species, error) {
	result := ts.SpeciesStore.Db.Create(&t)
	if result.Error != nil {
		return Species{}, result.Error
	}

	fmt.Printf("Species created: %d\n", t)

	return ts.Species, nil
}

func (ts *SpeciesServices) GetAllSpeciess(createdBy int) ([]Species, error) {
	speciess := []Species{}

	result := ts.SpeciesStore.Db.Where("created_by = ?", createdBy).Order("created_at desc").Find(&speciess)
	if result.Error != nil {
		return []Species{}, result.Error
	}
	return speciess, nil
}

func (ts *SpeciesServices) GetSpeciesById(t Species) (Species, error) {
	result := ts.SpeciesStore.Db.Where(&t, "created_by", "id").Find(&ts.Species)
	if result.Error != nil {
		return Species{}, result.Error
	}

	return ts.Species, nil
}

func (ts *SpeciesServices) UpdateSpecies(t Species) (Species, error) {
	result := ts.SpeciesStore.Db.Model(&ts.Species).Where(&t, "created_by", "id").Updates(&t)
	if result.Error != nil {
		return Species{}, result.Error
	}

	return ts.Species, nil
}

func (ts *SpeciesServices) DeleteSpecies(t Species) error {
	result := ts.SpeciesStore.Db.Delete(&t)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 1 {
		return errors.New("an affected row was expected")
	}

	return nil
}
