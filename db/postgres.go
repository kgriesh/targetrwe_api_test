package db

import (
	"database/sql"
	"fmt"
	"log"

	model "github.com/targetrw-api/model"

	_ "github.com/lib/pq"
)

const (
	HOST = "database"
	PORT = 5432
)

var ErrNoMatch = fmt.Errorf("no matching record")

type DbService interface {
	GetAllProducts() (*model.ProductList, error)
	AddProduct(product *model.Product) error
	GetProductById(itemId int) (model.Product, error)
	GetConnection() *sql.DB
}

type Database struct {
	Conn *sql.DB
}

func (db Database) GetConnection() *sql.DB {
	return db.Conn
}

func NewDbService(username, password, database string) (DbService, error) {
	db := Database{}
	ds := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		HOST, PORT, username, password, database)
	conn, err := sql.Open("postgres", ds)
	if err != nil {
		return db, err
	}
	db.Conn = conn
	err = db.Conn.Ping()
	if err != nil {
		return db, err
	}
	log.Println("Database connection established")
	return db, nil
}

func (db Database) GetAllProducts() (*model.ProductList, error) {
	list := &model.ProductList{}
	rows, err := db.Conn.Query("SELECT * FROM items ORDER BY ID DESC")
	if err != nil {
		return list, err
	}
	for rows.Next() {
		var product model.Product
		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.CreatedAt)
		if err != nil {
			return list, err
		}
		list.Products = append(list.Products, product)
	}
	return list, nil
}
func (db Database) AddProduct(product *model.Product) error {
	var id int
	var createdAt string
	query := `INSERT INTO items (name, description) VALUES ($1, $2) RETURNING id, created_at`
	err := db.Conn.QueryRow(query, product.Name, product.Description).Scan(&id, &createdAt)
	if err != nil {
		return err
	}
	product.ID = id
	product.CreatedAt = createdAt
	return nil
}

func (db Database) GetProductById(itemId int) (model.Product, error) {
	item := model.Product{}
	query := `SELECT * FROM items WHERE id = $1;`
	row := db.Conn.QueryRow(query, itemId)

	err := row.Scan(&item.ID, &item.Name, &item.Description, &item.CreatedAt)

	switch err {
	case sql.ErrNoRows:
		return item, ErrNoMatch
	default:
		return item, err
	}
}
