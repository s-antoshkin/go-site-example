package main

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
)

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv("DATABASE_URL"); exists {
		return value
	}
	return defaultVal
}

func create_table() {
	conn, err := pgx.Connect(context.Background(), getEnv("DATABASE_URL", ""))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close(context.Background())

	const query = `
		CREATE TABLE IF NOT EXISTS rsvp_test (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		phone TEXT NOT NULL,
		will_attend BOOLEAN NOT NULL)`
	_, err = conn.Exec(context.Background(), query)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func drop_table() {
	conn, err := pgx.Connect(context.Background(), getEnv("DATABASE_URL", ""))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close(context.Background())

	_, err = conn.Exec(context.Background(), "DROP TABLE IF EXISTS rsvp_test")
	if err != nil {
		fmt.Println(err)
		return
	}
}

func insert_record(query string) {
	conn, err := pgx.Connect(context.Background(), getEnv("DATABASE_URL", ""))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close(context.Background())

	_, err = conn.Exec(context.Background(), query)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func TestCount(t *testing.T) {
	var count int
	create_table()

	conn, err := pgx.Connect(context.Background(), getEnv("DATABASE_URL", ""))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close(context.Background())

	insert_record("INSERT INTO rsvp_test (name, email, phone, will_attend) VALUES ('Jhon', 'mail@email.com', '1234567', true)")
	insert_record("INSERT INTO rsvp_test (name, email, phone, will_attend) VALUES ('Will', 'mail2@email.com', '7654321', false)")
	insert_record("INSERT INTO rsvp_test (name, email, phone, will_attend) VALUES ('Mihail', 'mail3@email.com', '987654', true)")

	row := conn.QueryRow(context.Background(), "SELECT COUNT(*) FROM rsvp_test")
	err = row.Scan(&count)
	if err != nil {
		fmt.Println(err)
		return
	}

	if count != 3 {
		t.Errorf("Select query returned %d", count)
	}

	drop_table()

}

func TestQueryDB(t *testing.T) {
	create_table()

	conn, err := pgx.Connect(context.Background(), getEnv("DATABASE_URL", ""))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close(context.Background())

	query := "INSERT INTO rsvp_test (name, email, phone, will_attend) VALUES ('Jhon', 'mail@email.com', '1234567', true)"
	insert_record(query)

	rows, err := conn.Query(context.Background(), `SELECT * FROM rsvp_test WHERE name=$1`, `Jhon`)
	if err != nil {
		fmt.Println(err)
		return
	}

	var col1 int
	var col2, col3, col4 string
	var col5 bool
	for rows.Next() {
		rows.Scan(&col1, &col2, &col3, &col4, &col5)
	}

	if col2 != "Jhon" {
		t.Errorf("name returned %s\n", col2)
	}

	if col3 != "mail@email.com" {
		t.Errorf("email returned %s\n", col3)
	}

	if col4 != "1234567" {
		t.Errorf("phone returned %s\n", col3)
	}

	if col5 != true {
		t.Errorf("will_attend returned %s\n", col3)
	}

	drop_table()
}
