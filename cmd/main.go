package main

import (
	"github.com/tsironis/syntropic-studio/db"
	"github.com/tsironis/syntropic-studio/handlers"
	"github.com/tsironis/syntropic-studio/services"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// In production, the secret key of the CookieStore
// and database name would be obtained from a .env file
const (
	SECRET_KEY string = "secret"
	DB_NAME    string = "app_data.db"
)

func main() {

	e := echo.New()

	e.Static("/static", "assets")

	e.HTTPErrorHandler = handlers.CustomHTTPErrorHandler

	// Helpers Middleware
	// e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	// Session Middleware
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(SECRET_KEY))))

	store, err := db.NewStore(DB_NAME)
	if err != nil {
		e.Logger.Fatalf("failed to create store: %s", err)
	}

	store.Db.AutoMigrate(&services.User{})
	store.Db.AutoMigrate(&services.Todo{})
	store.Db.AutoMigrate(&services.Pattern{})
	store.Db.AutoMigrate(&services.Design{})
	store.Db.AutoMigrate(&services.Species{})

	us := services.NewUserServices(services.User{}, store)
	ah := handlers.NewAuthHandler(us)

	ts := services.NewTodoServices(services.Todo{}, store)
	th := handlers.NewTaskHandler(ts)

	ps := services.NewPatternServices(services.Pattern{}, store)
	ph := handlers.NewPatternHandler(ps)

	ds := services.NewDesignServices(services.Design{}, store)
	dh := handlers.NewDesignHandler(ds)

	ss := services.NewSpeciesServices(services.Species{}, store)
	sh := handlers.NewSpeciesHandler(ss)
	// Setting Routes
	handlers.SetupRoutes(e, ah, th, ph, dh, sh)

	// Start Server
	e.Logger.Fatal(e.Start(":8082"))
}

/*
https://gist.github.com/taforyou/544c60ffd072c9573971cf447c9fea44
https://gist.github.com/mhewedy/4e45e04186ed9d4e3c8c86e6acff0b17

https://github.com/CurtisVermeeren/gorilla-sessions-tutorial
*/
