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
type ItemList struct {
	Items []Product `json:"items"`
}

func (i *Product) Bind(r *http.Request) error {
	if i.Name == "" {
		return fmt.Errorf("name is a required field")
	}
	return nil
}
func (*ItemList) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
func (*Product) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
