package main

import (
	"database/sql"
	"errors"
	"io/ioutil"
	"log"

	_ "github.com/go-sql-driver/mysql"
	yaml "gopkg.in/yaml.v2"
)

type config struct {
	Database struct {
		Name     string
		Host     string
		User     string
		Password string
	}
}
type DB struct {
	*sql.DB
}

/*NewDB - it return db instance*/
func NewDB() *DB {
	data, err := ioutil.ReadFile("conf.yml")
	if err != nil {
		log.Fatal(err)
	}

	var config config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatal(err)
	}

	dataSourceName := config.Database.User + ":" + config.Database.Password + "@tcp(" + config.Database.Host + ":3306)/" + config.Database.Name
	db, err := sql.Open("mysql", dataSourceName)
	checkError(err)

	return &DB{db}
}

func (db *DB) store(id int, key string, encryptedData string) error {
	if id == 0 {
		return errors.New("Key must be provided to store")
	}

	if key == "" {
		return errors.New("Key must be provided to store")
	}

	if encryptedData == "" {
		return errors.New("encryptedData must be provided to store")
	}

	stmt, err := db.Prepare("INSERT INTO encrypted_data (id, `key`, `data`) values (?,?,?)")
	checkError(err)

	_, err = stmt.Exec(id, key, encryptedData)
	checkError(err)

	return nil
}

/*CheckIfIDAllReadyExists - it check if the passed id is already exists in db*/
func (db *DB) CheckIfIDAllReadyExists(id int) error {
	rows, err := db.Query("SELECT id FROM encrypted_data WHERE id = ?", id)
	checkError(err)
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			return err
		}
		return errors.New("There is another entry with the same id already exists")
	}

	return nil
}

func (db *DB) retrieveDecryptedContent(id int) (string, string, error) {
	rows, err := db.Query("SELECT `data`, `key` FROM encrypted_data WHERE id = ?", id)
	if err != nil {
		return "", "", err
	}

	var data, key string
	for rows.Next() {
		err := rows.Scan(&data, &key)
		if err != nil {
			return "", "", err
		}
	}

	return data, key, nil
}

func checkError(err error) {
	if err != nil {
		log.Println("Error while communicating with database:", err.Error())
		panic(err)
	}
}
