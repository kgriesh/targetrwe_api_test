package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	database "github.com/targetrw-api/db"
	model "github.com/targetrw-api/model"
)

type MockDbService struct {
	mock.Mock
}

func (m *MockDbService) AddProduct(p *model.Product) error {
	args := m.Called(p)
	return args.Error(0)
}

func (m *MockDbService) GetAllProducts() (*model.ProductList, error) {
	args := m.Called()
	return args.Get(0).(*model.ProductList), args.Error(0)
}

func (m *MockDbService) GetProductById(id int) (model.Product, error) {
	args := m.Called(id)
	return args.Get(0).(model.Product), args.Error(1)
}

func (m *MockDbService) GetConnection() *sql.DB {
	args := m.Called()
	return args.Get(0).(*sql.DB)
}

func Test_GetProduct(t *testing.T) {
	testCases := []struct {
		desc    string
		prodId  interface{}
		dbError error
	}{
		{
			desc:    "happy path",
			prodId:  1,
			dbError: nil,
		},
		{
			desc:    "invalid product id",
			prodId:  []byte("1"),
			dbError: errors.New("invalid productId"),
		},
		{
			desc:    "db call results in error",
			prodId:  1,
			dbError: errors.New("boom"),
		},
		{
			desc:    "db call results in not found error",
			prodId:  1,
			dbError: database.ErrNoMatch,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {

			mockDb := &MockDbService{}
			mockDb.On("GetProductById", 1).Return(model.Product{}, tc.dbError)

			r := chi.NewRouter()
			r.Get("/products/{itemId}", GetProduct(mockDb))
			ts := httptest.NewServer(r)
			defer ts.Close()

			req, _ := http.NewRequest("GET", fmt.Sprintf("%s/products/%v", ts.URL, tc.prodId), nil)
			res, err := ts.Client().Do(req)

			assert.NoError(t, err)

			if tc.dbError == nil {
				assert.Equal(t, http.StatusOK, res.StatusCode)
			} else if tc.dbError == database.ErrNoMatch || tc.dbError.Error() == "invalid productId" {
				assert.Equal(t, http.StatusBadRequest, res.StatusCode)
			} else {
				assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
			}

		})
	}
}
