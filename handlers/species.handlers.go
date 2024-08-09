package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/tsironis/syntropic-studio/services"
	"github.com/tsironis/syntropic-studio/views/species_views"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

/********** Handlers for Species Views **********/

type SpeciesService interface {
	CreateSpecies(t services.Species) (services.Species, error)
	GetAllSpeciess(createdBy int) ([]services.Species, error)
	GetSpeciesById(t services.Species) (services.Species, error)
	UpdateSpecies(t services.Species) (services.Species, error)
	DeleteSpecies(t services.Species) error
}

func NewSpeciesHandler(ts SpeciesService) *SpeciesHandler {

	return &SpeciesHandler{
		SpeciesServices: ts,
	}
}

type SpeciesHandler struct {
	SpeciesServices SpeciesService
}

func (th *SpeciesHandler) createSpeciesHandler(c echo.Context) error {
	// isError = false
	c.Set("ISERROR", false)
	fromProtected, ok := c.Get("FROMPROTECTED").(bool)
	if !ok {
		return errors.New("invalid type for key 'FROMPROTECTED'")
	}

	if c.Request().Method == "POST" {
		species := services.Species{
			CreatedBy:            c.Get(user_id_key).(int),
			BinomialNomenclature: strings.Trim(c.FormValue("binomial_nomenclature"), " "),
			CommonName:           strings.Trim(c.FormValue("common_name"), " "),
			Stratum:              strings.Trim(c.FormValue("stratum"), " "),
			System:               strings.Trim(c.FormValue("system"), " "),
			Lifecycle:            strings.Trim(c.FormValue("lifecycle"), " "),
			BiologicalForm:       strings.Trim(c.FormValue("biological_form"), " "),
			USDAHardinessZone:    strings.Trim(c.FormValue("usda_hardiness_zone"), " "),
			PrimaryFunction:      strings.Trim(c.FormValue("primary_function"), " "),
			SecondaryFunction:    strings.Trim(c.FormValue("secondary_function"), " "),
			Notes:                strings.Trim(c.FormValue("notes"), " "),
		}

		_, err := th.SpeciesServices.CreateSpecies(species)
		if err != nil {
			return err
		}

		setFlashmessages(c, "success", "Species created successfully!!")

		return c.Redirect(http.StatusSeeOther, "/species/list")
	}

	return renderView(c, species_views.SpeciesIndex(
		"| Create Species",
		c.Get(username_key).(string),
		fromProtected,
		c.Get("ISERROR").(bool),
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		species_views.CreateSpecies(),
	))
}

func (th *SpeciesHandler) speciesListHandler(c echo.Context) error {
	// isError = false
	c.Set("ISERROR", false)
	fromProtected, ok := c.Get("FROMPROTECTED").(bool)
	if !ok {
		return errors.New("invalid type for key 'FROMPROTECTED'")
	}

	userId := c.Get(user_id_key).(int)

	allSpecies, err := th.SpeciesServices.GetAllSpeciess(userId)
	if err != nil {
		return err
	}

	titlePage := fmt.Sprintf(
		"| %s's Species List",
		cases.Title(language.English).String(c.Get(username_key).(string)),
	)

	return renderView(c, species_views.SpeciesIndex(
		titlePage,
		c.Get(username_key).(string),
		fromProtected,
		c.Get("ISERROR").(bool),
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		species_views.SpeciesList(titlePage, allSpecies),
	))
}

func (th *SpeciesHandler) updateSpeciesHandler(c echo.Context) error {
	// isError = false
	c.Set("ISERROR", false)
	fromProtected, ok := c.Get("FROMPROTECTED").(bool)
	if !ok {
		return errors.New("invalid type for key 'FROMPROTECTED'")
	}

	idParams, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	t := services.Species{
		ID:        idParams,
		CreatedBy: c.Get(user_id_key).(int),
	}

	species, err := th.SpeciesServices.GetSpeciesById(t)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {

			return echo.NewHTTPError(
				echo.ErrNotFound.Code,
				fmt.Sprintf(
					"something went wrong: %s",
					err,
				))
		}

		return echo.NewHTTPError(
			echo.ErrInternalServerError.Code,
			fmt.Sprintf(
				"something went wrong: %s",
				err,
			))
	}

	if c.Request().Method == "POST" {
		var private bool
		if c.FormValue("private") == "on" {
			private = true
		} else {
			private = false
		}

		species := services.Species{
			Private:              private,
			CreatedBy:            c.Get(user_id_key).(int),
			BinomialNomenclature: c.FormValue("binomial_nomenclature"),
			CommonName:           c.FormValue("common_name"),
			Stratum:              c.FormValue("stratum"),
			System:               c.FormValue("system"),
			Lifecycle:            c.FormValue("lifecycle"),
			BiologicalForm:       c.FormValue("biological_form"),
			USDAHardinessZone:    c.FormValue("usda_hardiness_zone"),
			PrimaryFunction:      c.FormValue("primary_form"),
			SecondaryFunction:    c.FormValue("secondary_form"),
			ID:                   idParams,
		}

		_, err := th.SpeciesServices.UpdateSpecies(species)
		if err != nil {
			return err
		}

		setFlashmessages(c, "success", "Species successfully updated")

		return c.Redirect(http.StatusSeeOther, "/species/list")
	}

	return renderView(c, species_views.SpeciesIndex(
		fmt.Sprintf("| Edit Species #%d", species.ID),
		c.Get(username_key).(string),
		fromProtected,
		c.Get("ISERROR").(bool),
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"), // ↓ getting time zone from context ↓
		species_views.UpdateSpecies(species, c.Get(tzone_key).(string)),
	))
}

func (th *SpeciesHandler) deleteSpeciesHandler(c echo.Context) error {
	idParams, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		fmt.Println(err)
		return err
	}

	t := services.Species{
		CreatedBy: c.Get(user_id_key).(int),
		ID:        idParams,
	}

	err = th.SpeciesServices.DeleteSpecies(t)
	if err != nil {
		if strings.Contains(err.Error(), "an affected row was expected") {

			return echo.NewHTTPError(
				echo.ErrNotFound.Code,
				fmt.Sprintf(
					"something went wrong: %s",
					err,
				))
		}

		return echo.NewHTTPError(
			echo.ErrInternalServerError.Code,
			fmt.Sprintf(
				"something went wrong: %s",
				err,
			))
	}

	setFlashmessages(c, "success", "Species successfully deleted!!")

	return c.Redirect(http.StatusSeeOther, "/species/list")
}

func (th *SpeciesHandler) logoutHandler(c echo.Context) error {
	sess, _ := session.Get(auth_sessions_key, c)
	// Revoke users authentication
	sess.Values = map[interface{}]interface{}{
		auth_key:     false,
		user_id_key:  "",
		username_key: "",
		tzone_key:    "",
	}
	sess.Save(c.Request(), c.Response())

	setFlashmessages(c, "success", "You have successfully logged out!!")

	// fromProtected = false
	c.Set("FROMPROTECTED", false)

	return c.Redirect(http.StatusSeeOther, "/login")
}
