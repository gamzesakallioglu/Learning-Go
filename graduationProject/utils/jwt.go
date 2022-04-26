package utils

import (
	"errors"

	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func GetUserFromCtx(c *gin.Context) (*jwt.DecodedToken, error) {

	user, exist := c.Get("__user__")
	if !exist {
		return nil, errors.New("customer not found")
	}

	return user.(*jwt.DecodedToken), nil
}
