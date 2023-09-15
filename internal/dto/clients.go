package dto

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type (
	Client struct {
		Id        int      `json:"id"`
		ManagerId int      `json:"manager_id"`
		Name      string   `json:"name"`
		Email     string   `json:"email"`
		Phone     string   `json:"phone"`
		New       bool     `json:"new"`
		InArchive bool     `json:"in_archive"`
		Result    *Result  `json:"result,omitempty"`
		Results   []Result `json:"results"`
		Date      int      `json:"date"`
	}
	Result struct {
		Date         int                    `json:"date"`
		DiagnosticID int                    `json:"diagnostic-id"`
		Data         map[string]interface{} `json:"data"`
		OpenAnswer   string                 `json:"openAnswer"`
	}
	GetClient struct {
		Id      int  `json:"id"`
		IsAdmin bool `json:"is_admin"`
	}
)

func (g GetClient) Validate() error {
	return validation.ValidateStruct(&g,
		validation.Field(&g.Id, validation.Required))
}

func (c Client) ValidateAddResult() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.ManagerId, validation.Required),
		validation.Field(&c.Name, validation.Required),
		validation.Field(&c.Email, validation.Required),
		validation.Field(&c.Phone, validation.Required),
		validation.Field(&c.Result, validation.Required),
		validation.Field(&c.Date, validation.Required))
}

func (c Client) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Id, validation.Required))
}
