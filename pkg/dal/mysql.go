package dal

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var (
	dbdriver = "mysql"
	dbname   = "messageDB"
)

// MySQLDriver ... mysql db driver
type MySQLDriver struct {
	driver *sql.DB
}

// TODO: Do away with mysql raw queries

//GetMySQLDriver ... get a mysql DB driver
func GetMySQLDriver() *MySQLDriver {
	if os.Getenv("DB_DRIVER") != "" {
		dbdriver = os.Getenv("DB_DRIVER")
	}
	dbhost := os.Getenv("DB_HOST")
	dbport := os.Getenv("DB_PORT")
	dbuser := os.Getenv("DB_USER")
	dbpasswd := os.Getenv("DB_PASSWD")
	if os.Getenv("DB_NAME") != "" {
		dbname = os.Getenv("DB_NAME")
	}
	dbaccessurl := dbuser + ":" + dbpasswd + "@" + "tcp(" + dbhost + ":" + dbport + ")" + "/" + dbname

	var db MySQLDriver
	var err error
	db.driver, err = sql.Open(dbdriver, dbaccessurl)
	if err != nil {
		log.Println("Error connecting to database: ", err)
	}
	return &db
}

// AddMessage ... add new message
func (db *MySQLDriver) AddMessage(msg *Message) error {
	var statement string
	statement = fmt.Sprintf("INSERT INTO messageDB.messages (msg) VALUES ('%s');", msg.Message)
	if msg.ID > 0 {
		statement = fmt.Sprintf("INSERT INTO messageDB.messages (id, msg) VALUES ('%d, %s');", msg.ID, msg.Message)
	}
	_, err := db.driver.Exec(statement)

	if err != nil {
		return err
	}

	err = db.driver.QueryRow("SELECT LAST_INSERT_ID()").Scan(&msg.ID)

	if err != nil {
		return err
	}

	return nil
}

// GetMessage ... get message by ID
func (db *MySQLDriver) GetMessage(msg *Message) error {
	statement := fmt.Sprintf("SELECT msg FROM messageDB.messages WHERE id=%d", msg.ID)
	return db.driver.QueryRow(statement).Scan(&msg.Message)
}

// UpdateMessage ... update message by ID
func (db *MySQLDriver) UpdateMessage(msg *Message) error {
	statement := fmt.Sprintf("UPDATE messageDB.messages SET msg='%s' WHERE id=%d", msg.Message, msg.ID)
	_, err := db.driver.Exec(statement)
	return err
}

// DeleteMessage ... delete message by ID
func (db *MySQLDriver) DeleteMessage(msg *Message) error {
	statement := fmt.Sprintf("DELETE FROM messageDB.messages WHERE id=%d", msg.ID)
	_, err := db.driver.Exec(statement)
	return err
}

//GetAll ... get all messages
func (db *MySQLDriver) GetAll() ([]Message, error) {
	statement := fmt.Sprintf("SELECT id, msg FROM messageDB.messages")
	rows, err := db.driver.Query(statement)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	messages := []Message{}

	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.ID, &msg.Message); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, nil
}
