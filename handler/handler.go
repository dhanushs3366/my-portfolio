package handler

import (
	"database/sql"
	"dhanushs3366/my-portfolio/services"
	"dhanushs3366/my-portfolio/services/logger"
	"dhanushs3366/my-portfolio/services/user"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Handler struct {
	router    *echo.Echo
	userStore *user.UserStore
	logStore  *logger.LogStore
	// have jwt and config future
}

func Init(db *sql.DB) *Handler {
	h := Handler{
		router:    echo.New(),
		userStore: user.NewUserStore(db),
		logStore:  logger.NewLogStore(db),
	}
	h.router.Use(middleware.Logger())
	h.router.Use(middleware.Recover())
	h.router.Use(middleware.CORS())

	h.router.GET("/hello", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "HIIII")
	})
	h.router.GET("/log-details", h.GetLogDetails)

	apiRoutes := h.router.Group("/api")
	apiRoutes.POST("/log-details", h.PostLogDetails)

	// h.router.GET("/repos", GetRepos)
	// h.router.GET("/git-user", GetGitUser)

	h.router.POST("/login", h.Login)

	adminRoutes := h.router.Group("/admins")
	adminRoutes.GET("/hello", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "hello admin")
	})
	adminRoutes.PATCH("/user", h.UpdatePassword)

	adminRoutes.Use(services.ValidateJWT)
	apiRoutes.Use(services.ValidateLoggerToken)
	return &h
}

func (h *Handler) Run(port uint) {
	h.router.Start(fmt.Sprintf(":%d", port))
}
