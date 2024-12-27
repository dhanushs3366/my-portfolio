package handler

import (
	"dhanushs3366/my-portfolio/services"
	"dhanushs3366/my-portfolio/services/db"
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Handler struct {
	router *echo.Echo
	store  *db.Store
	// have jwt and config future
}

func Init(store *db.Store) *Handler {
	FE_URL := os.Getenv("FE_URL")

	h := Handler{
		router: echo.New(),
		store:  store,
	}
	h.router.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"time":"${time_rfc3339}", "method":"${method}", "uri":"${uri}", "status":${status}, "latency":"${latency_human}", "bytes_in":${bytes_in}, "bytes_out":${bytes_out}}` + "\n",
	}))
	h.router.Use(middleware.Recover())
	h.router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{FE_URL},
		AllowMethods:     []string{echo.GET, echo.POST, echo.DELETE, echo.PUT, echo.PATCH},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowCredentials: true,
	}))

	adminRoutes := h.router.Group("/admins")
	apiRoutes := h.router.Group("/api")

	adminRoutes.Use(services.ValidateJWT)
	apiRoutes.Use(services.ValidateLoggerToken)

	// everyone
	h.router.GET("/hello", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "HIIII")
	})

	h.router.GET("/log-details", h.getLogDetails)
	h.router.GET("/blogs", h.getBlogs)
	h.router.GET("/blogs/:ID", h.getBlog)
	h.router.GET("/repos", h.fetchGitRepos)

	h.router.POST("/login", h.login)

	// api
	apiRoutes.POST("/log-details", h.postLogDetails)

	// admin
	adminRoutes.GET("/hello", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "hello admin")
	})

	// admins/blogs
	adminRoutes.POST("/blog", h.createBlog)
	adminRoutes.PATCH("/blog/:id", h.editBlog)
	adminRoutes.DELETE("/blog/:id", h.deleteBlog)

	// admin user
	adminRoutes.POST("/user", h.createAdmin)
	adminRoutes.PATCH("/user", h.updatePassword)

	return &h
}

func (h *Handler) Run(port uint) {
	h.router.Start(fmt.Sprintf(":%d", port))
}
