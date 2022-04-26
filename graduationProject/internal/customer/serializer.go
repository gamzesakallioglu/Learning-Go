package customer

import (
	"time"

	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/internal/api"
	"github.com/gamze.sakallioglu/learningGo/bitirme-projesi-gamzesakallioglu/internal/models"
)

func CustomerToResponse(c *models.Customer) *api.Customer {

	var isPassive bool
	if c.PassiveDate.After(time.Now()) {
		isPassive = true
	} else {
		isPassive = false
	}
	return &api.Customer{
		ID:        c.ID.String(),
		Name:      c.Name,
		Address:   c.Address,
		Email:     c.Email,
		Phone:     c.Phone,
		IsPassive: isPassive,
	}

}
