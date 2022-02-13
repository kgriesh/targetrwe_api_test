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

// ErrNoMatch - row that doesn't exist
var ErrNoMatch = fmt.Errorf("no matching record")

type Database struct {
	Conn *sql.DB
}

func Initialize(username, password, database string) (Database, error) {
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

func (db Database) GetAllItems() (*model.ItemList, error) {
	list := &model.ItemList{}
	rows, err := db.Conn.Query("SELECT * FROM items ORDER BY ID DESC")
	test, _ := db.Conn.Query("\\d")
	log.Printf("doing query test %v", test)
	if err != nil {
		log.Println("boom *************", err)
		return list, err
	}
	for rows.Next() {
		var item model.Product
		err := rows.Scan(&item.ID, &item.Name, &item.Description, &item.CreatedAt)
		if err != nil {
			return list, err
		}
		list.Items = append(list.Items, item)
	}
	return list, nil
}
func (db Database) AddItem(item *model.Product) error {
	var id int
	var createdAt string
	query := `INSERT INTO items (name, description) VALUES ($1, $2) RETURNING id, created_at`
	err := db.Conn.QueryRow(query, item.Name, item.Description).Scan(&id, &createdAt)
	if err != nil {
		return err
	}
	item.ID = id
	item.CreatedAt = createdAt
	return nil
}
func (db Database) GetItemById(itemId int) (model.Product, error) {
	item := model.Product{}
	query := `SELECT * FROM items WHERE id = $1;`
	row := db.Conn.QueryRow(query, itemId)
	switch err := row.Scan(&item.ID, &item.Name, &item.Description, &item.CreatedAt); err {
	case sql.ErrNoRows:
		return item, ErrNoMatch
	default:
		return item, err
	}
}
func (db Database) DeleteItem(itemId int) error {
	query := `DELETE FROM items WHERE id = $1;`
	_, err := db.Conn.Exec(query, itemId)
	switch err {
	case sql.ErrNoRows:
		return ErrNoMatch
	default:
		return err
	}
}
func (db Database) UpdateItem(itemId int, itemData model.Product) (model.Product, error) {
	item := model.Product{}
	query := `UPDATE items SET name=$1, description=$2 WHERE id=$3 RETURNING id, name, description, created_at;`
	err := db.Conn.QueryRow(query, itemData.Name, itemData.Description, itemId).Scan(&item.ID, &item.Name, &item.Description, &item.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return item, ErrNoMatch
		}
		return item, err
	}
	return item, nil
}
