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
	"github.com/tsironis/syntropic-studio/views/design_views"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

/********** Handlers for Design Views **********/

type DesignService interface {
	CreateDesign(t services.Design) (services.Design, error)
	GetAllDesigns(createdBy int) ([]services.Design, error)
	GetDesignById(t services.Design) (services.Design, error)
	UpdateDesign(t services.Design) (services.Design, error)
	DeleteDesign(t services.Design) error
}

func NewDesignHandler(ps DesignService) *DesignHandler {

	return &DesignHandler{
		DesignServices: ps,
	}
}

type DesignHandler struct {
	DesignServices DesignService
}

func (th *DesignHandler) createDesignHandler(c echo.Context) error {
	// isError = false
	c.Set("ISERROR", false)
	fromProtected, ok := c.Get("FROMPROTECTED").(bool)
	if !ok {
		return errors.New("invalid type for key 'FROMPROTECTED'")
	}

	if c.Request().Method == "POST" {
		design := services.Design{
			CreatedBy: c.Get(user_id_key).(int),
			Title:     strings.Trim(c.FormValue("title"), " "),
			QGISData:  strings.Trim(c.FormValue("QGISData"), " "),
		}

		_, err := th.DesignServices.CreateDesign(design)
		if err != nil {
			return err
		}

		setFlashmessages(c, "success", "Design created successfully!!")

		return c.Redirect(http.StatusSeeOther, "/design/list")
	}

	return renderView(c, design_views.DesignIndex(
		"| Create Design",
		c.Get(username_key).(string),
		fromProtected,
		c.Get("ISERROR").(bool),
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		design_views.CreateDesign(),
	))
}

func (th *DesignHandler) designListHandler(c echo.Context) error {
	// isError = false
	c.Set("ISERROR", false)
	fromProtected, ok := c.Get("FROMPROTECTED").(bool)
	if !ok {
		return errors.New("invalid type for key 'FROMPROTECTED'")
	}

	userId := c.Get(user_id_key).(int)

	designs, err := th.DesignServices.GetAllDesigns(userId)
	if err != nil {
		return err
	}

	titlePage := fmt.Sprintf(
		"| %s's Designs",
		cases.Title(language.English).String(c.Get(username_key).(string)),
	)

	return renderView(c, design_views.DesignIndex(
		titlePage,
		c.Get(username_key).(string),
		fromProtected,
		c.Get("ISERROR").(bool),
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		design_views.DesignList(titlePage, designs),
	))
}

func (th *DesignHandler) updateDesignHandler(c echo.Context) error {
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

	t := services.Design{
		ID:        idParams,
		CreatedBy: c.Get(user_id_key).(int),
	}

	design, err := th.DesignServices.GetDesignById(t)
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

		design := services.Design{
			Title:     strings.Trim(c.FormValue("title"), " "),
			QGISData:  strings.Trim(c.FormValue("qgis_data"), " "),
			Status:    status,
			CreatedBy: c.Get(user_id_key).(int),
			ID:        idParams,
		}

		_, err := th.DesignServices.UpdateDesign(design)
		if err != nil {
			return err
		}

		setFlashmessages(c, "success", "Design successfully updated!!")

		return c.Redirect(http.StatusSeeOther, "/design/list")
	}

	return renderView(c, design_views.DesignIndex(
		fmt.Sprintf("| Edit Design #%d", design.ID),
		c.Get(username_key).(string),
		fromProtected,
		c.Get("ISERROR").(bool),
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"), // ↓ getting time zone from context ↓
		design_views.UpdateDesign(design, c.Get(tzone_key).(string)),
	))
}

func (th *DesignHandler) deleteDesignHandler(c echo.Context) error {
	idParams, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		fmt.Println(err)
		return err
	}

	t := services.Design{
		CreatedBy: c.Get(user_id_key).(int),
		ID:        idParams,
	}

	err = th.DesignServices.DeleteDesign(t)
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

	setFlashmessages(c, "success", "Design successfully deleted!!")

	return c.Redirect(http.StatusSeeOther, "/design/list")
}

func (th *DesignHandler) logoutHandler(c echo.Context) error {
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
