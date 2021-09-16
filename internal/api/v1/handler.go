package v1

import (
	"github.com/ekifel/moneysaverz/internal/service"
	"github.com/ekifel/moneysaverz/pkg/auth"
	"github.com/gin-gonic/gin"
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

func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		h.initUsersRoutes(v1)
	}
}

// func parseIdFromPath(c *gin.Context, param string) (primitive.ObjectID, error) {
// 	idParam := c.Param(param)
// 	if idParam == "" {
// 		return primitive.ObjectID{}, fmt.Errorf("empty id param")
// 	}

// 	id, err := primitive.ObjectIDFromHex(idParam)
// 	if err != nil {
// 		return primitive.ObjectID{}, fmt.Errorf("invalid id param")
// 	}

// 	return id, nil
// }
