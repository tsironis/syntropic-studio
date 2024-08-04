package handlers

import (
	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, ah *AuthHandler, th *TaskHandler, ph *PatternHandler) {
	e.GET("/", ah.flagsMiddleware(ah.homeHandler))
	e.GET("/login", ah.flagsMiddleware(ah.loginHandler))
	e.POST("/login", ah.flagsMiddleware(ah.loginHandler))
	e.GET("/register", ah.flagsMiddleware(ah.registerHandler))
	e.POST("/register", ah.flagsMiddleware(ah.registerHandler))

	todoGroup := e.Group("/todo", ah.authMiddleware)
	/* ↓ Protected Routes ↓ */
	todoGroup.GET("/list", th.todoListHandler)
	todoGroup.GET("/create", th.createTodoHandler)
	todoGroup.POST("/create", th.createTodoHandler)
	todoGroup.GET("/edit/:id", th.updateTodoHandler)
	todoGroup.POST("/edit/:id", th.updateTodoHandler)
	todoGroup.DELETE("/delete/:id", th.deleteTodoHandler)
	todoGroup.POST("/logout", th.logoutHandler)

	patternGroup := e.Group("/pattern", ah.authMiddleware)
	/* ↓ Protected Routes ↓ */
	patternGroup.GET("/list", ph.patternListHandler)
	patternGroup.GET("/create", ph.createPatternHandler)
	patternGroup.POST("/create", ph.createPatternHandler)
	patternGroup.GET("/edit/:id", ph.updatePatternHandler)
	patternGroup.POST("/edit/:id", ph.updatePatternHandler)
	patternGroup.DELETE("/delete/:id", ph.deletePatternHandler)
	patternGroup.POST("/logout", ph.logoutHandler)

	/* ↓ Fallback Page ↓ */
	e.GET("/*", RouteNotFoundHandler)
}
