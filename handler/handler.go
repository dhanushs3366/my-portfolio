package handler

import (
	"database/sql"
	"dhanushs3366/my-portfolio/services"
	"dhanushs3366/my-portfolio/services/blog"
	"dhanushs3366/my-portfolio/services/logger"
	"dhanushs3366/my-portfolio/services/user"
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Handler struct {
	router    *echo.Echo
	userStore *user.UserStore
	logStore  *logger.LogStore
	blogStore *blog.BlogStore
	// have jwt and config future
}

func Init(db *sql.DB) *Handler {
	FE_URL := os.Getenv("FE_URL")

	h := Handler{
		router:    echo.New(),
		userStore: user.NewUserStore(db),
		logStore:  logger.NewLogStore(db),
		blogStore: blog.NewBlogStore(db),
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

	h.router.GET("/debug/cookies", func(c echo.Context) error {
		cookies := c.Cookies()
		h.router.Logger.Printf("Cookies: %+v\n", cookies)
		return c.JSON(http.StatusOK, cookies)
	})

	adminRoutes := h.router.Group("/admins")
	apiRoutes := h.router.Group("/api")

	adminRoutes.Use(services.ValidateJWT)
	apiRoutes.Use(services.ValidateLoggerToken)

	// everyone
	h.router.GET("/hello", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "HIIII")
	})
	h.router.GET("/log-details", h.GetLogDetails)
	h.router.POST("/login", h.Login)

	// api
	apiRoutes.POST("/log-details", h.PostLogDetails)

	// admin
	adminRoutes.GET("/hello", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "hello admin")
	})
	adminRoutes.POST("/blog", h.createBlog)

	adminRoutes.PATCH("/user", h.UpdatePassword)
	return &h
}

func (h *Handler) Run(port uint) {
	h.router.Start(fmt.Sprintf(":%d", port))
}
