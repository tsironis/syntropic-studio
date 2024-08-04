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
	"github.com/tsironis/syntropic-studio/views/pattern_views"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

/********** Handlers for Pattern Views **********/

type PatternService interface {
	CreatePattern(t services.Pattern) (services.Pattern, error)
	GetAllPatterns(createdBy int) ([]services.Pattern, error)
	GetPatternById(t services.Pattern) (services.Pattern, error)
	UpdatePattern(t services.Pattern) (services.Pattern, error)
	DeletePattern(t services.Pattern) error
}

func NewPatternHandler(ps PatternService) *PatternHandler {

	return &PatternHandler{
		PatternServices: ps,
	}
}

type PatternHandler struct {
	PatternServices PatternService
}

func (th *PatternHandler) createPatternHandler(c echo.Context) error {
	// isError = false
	c.Set("ISERROR", false)
	fromProtected, ok := c.Get("FROMPROTECTED").(bool)
	if !ok {
		return errors.New("invalid type for key 'FROMPROTECTED'")
	}

	if c.Request().Method == "POST" {
		pattern := services.Pattern{
			CreatedBy:   c.Get(user_id_key).(int),
			Title:       strings.Trim(c.FormValue("title"), " "),
			Description: strings.Trim(c.FormValue("description"), " "),
		}

		_, err := th.PatternServices.CreatePattern(pattern)
		if err != nil {
			return err
		}

		setFlashmessages(c, "success", "Pattern created successfully!!")

		return c.Redirect(http.StatusSeeOther, "/pattern/list")
	}

	return renderView(c, pattern_views.PatternIndex(
		"| Create Pattern",
		c.Get(username_key).(string),
		fromProtected,
		c.Get("ISERROR").(bool),
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		pattern_views.CreatePattern(),
	))
}

func (th *PatternHandler) patternListHandler(c echo.Context) error {
	// isError = false
	c.Set("ISERROR", false)
	fromProtected, ok := c.Get("FROMPROTECTED").(bool)
	if !ok {
		return errors.New("invalid type for key 'FROMPROTECTED'")
	}

	userId := c.Get(user_id_key).(int)

	patterns, err := th.PatternServices.GetAllPatterns(userId)
	if err != nil {
		return err
	}

	titlePage := fmt.Sprintf(
		"| %s's Pattern List",
		cases.Title(language.English).String(c.Get(username_key).(string)),
	)

	return renderView(c, pattern_views.PatternIndex(
		titlePage,
		c.Get(username_key).(string),
		fromProtected,
		c.Get("ISERROR").(bool),
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		pattern_views.PatternList(titlePage, patterns),
	))
}

func (th *PatternHandler) updatePatternHandler(c echo.Context) error {
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

	t := services.Pattern{
		ID:        idParams,
		CreatedBy: c.Get(user_id_key).(int),
	}

	pattern, err := th.PatternServices.GetPatternById(t)
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
		var status bool
		if c.FormValue("status") == "on" {
			status = true
		} else {
			status = false
		}

		pattern := services.Pattern{
			Title:       strings.Trim(c.FormValue("title"), " "),
			Description: strings.Trim(c.FormValue("description"), " "),
			Status:      status,
			CreatedBy:   c.Get(user_id_key).(int),
			ID:          idParams,
		}

		_, err := th.PatternServices.UpdatePattern(pattern)
		if err != nil {
			return err
		}

		setFlashmessages(c, "success", "Pattern successfully updated!!")

		return c.Redirect(http.StatusSeeOther, "/pattern/list")
	}

	return renderView(c, pattern_views.PatternIndex(
		fmt.Sprintf("| Edit Pattern #%d", pattern.ID),
		c.Get(username_key).(string),
		fromProtected,
		c.Get("ISERROR").(bool),
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"), // ↓ getting time zone from context ↓
		pattern_views.UpdatePattern(pattern, c.Get(tzone_key).(string)),
	))
}

func (th *PatternHandler) deletePatternHandler(c echo.Context) error {
	idParams, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		fmt.Println(err)
		return err
	}

	t := services.Pattern{
		CreatedBy: c.Get(user_id_key).(int),
		ID:        idParams,
	}

	err = th.PatternServices.DeletePattern(t)
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

	setFlashmessages(c, "success", "Pattern successfully deleted!!")

	return c.Redirect(http.StatusSeeOther, "/pattern/list")
}

func (th *PatternHandler) logoutHandler(c echo.Context) error {
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
