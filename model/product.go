package models

import (
	"fmt"
	"net/http"
)

type Product struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
}
type ProductList struct {
	Products []Product `json:"products"`
}

func (p *Product) Bind(r *http.Request) error {
	if p.Name == "" {
		return fmt.Errorf("name is a required field")
	}
	return nil
}
func (*ProductList) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
func (*Product) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
