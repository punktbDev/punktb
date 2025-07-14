package dto

import validation "github.com/go-ozzo/ozzo-validation"

type (
	Manager struct {
		Id                   int    `json:"id"`
		Login                string `json:"login,omitempty"`
		Password             string `json:"password,omitempty"`
		Name                 string `json:"name"`
		Surname              string `json:"surname"`
		Phone                string `json:"phone"`
		IsAdmin              bool   `json:"is_admin"`
		IsActive             bool   `json:"is_active"`
		AvailableDiagnostics []int  `json:"available_diagnostics"`
		IsFullAccess         bool   `json:"is_full_access"`
	}
)

func (m Manager) ValidateUpdate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Id, validation.Required),
		validation.Field(&m.Login, validation.Required),
		validation.Field(&m.Password, validation.Required),
		validation.Field(&m.Name, validation.Required),
		validation.Field(&m.Surname, validation.Required),
		//validation.Field(&m.Phone, validation.Required),
	)
}

func (m Manager) ValidateCreate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Login, validation.Required),
		validation.Field(&m.Password, validation.Required),
		//validation.Field(&m.Name, validation.Required),
		//validation.Field(&m.Surname, validation.Required),
		//validation.Field(&m.Phone, validation.Required),
	)
}
