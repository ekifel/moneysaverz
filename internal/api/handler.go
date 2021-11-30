package api

import (
	"fmt"
	"net/http"

	v1 "github.com/ekifel/moneysaverz/internal/api/v1"
	"github.com/ekifel/moneysaverz/internal/config"
	"github.com/ekifel/moneysaverz/internal/service"
	"github.com/ekifel/moneysaverz/pkg/auth"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/swag/example/basic/docs"
)

type Handler struct {
	services     *service.Services
	tokenManager auth.TokenManager
}

func NewHandler(services *service.Services, tokenManager auth.TokenManager) *Handler {
	return &Handler{
		services:     services,
		tokenManager: tokenManager,
	}
}

func (h *Handler) Init(cfg *config.Config) *gin.Engine {
	router := gin.Default()

	router.Use(
		gin.Recovery(),
		gin.Logger(),
		corsMiddleware,
	)

	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", cfg.HTTP.Host, cfg.HTTP.Port)
	if cfg.Environment != config.EnvLocal {
		docs.SwaggerInfo.Host = cfg.HTTP.Host
	}

	// Init router
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	h.initAPI(router)

	return router
}

func (h *Handler) initAPI(router *gin.Engine) {
	handlerV1 := v1.NewHandler(h.services, h.tokenManager)
	api := router.Group("/api")
	{
		handlerV1.Init(api)
	}
}
