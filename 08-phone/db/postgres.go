package db

import (
	"database/sql"
	"fmt"
	"github.com/alextsa22/gophercises/08-phone/phone"
	"regexp"
)

var (
	reg            = regexp.MustCompile("\\D")
	exampleNumbers = []string{
		"1234567890", "123 456 7891",
		"(123) 456 7892", "(123) 456-7893",
		"123-456-7894", "123-456-7890",
		"1234567892", "(123)456-7892",
	}
)

type PostgresDB struct {
	db *sql.DB
}

func NewPostgresDB(db *sql.DB, name string) (*PostgresDB, error) {
	if err := resetDB(db, name); err != nil {
		return nil, err
	}

	return &PostgresDB{db: db}, nil
}

func resetDB(db *sql.DB, name string) error {
	query := fmt.Sprintf("DROP DATABASE IF EXISTS %s", name)
	if _, err := db.Exec(query); err != nil {
		return err
	}

	return createDB(db, name)
}

func createDB(db *sql.DB, name string) error {
	query := fmt.Sprintf("CREATE DATABASE %s", name)
	if _, err := db.Exec(query); err != nil {
		return err
	}

	return nil
}

func (p *PostgresDB) InitExample() error {
	if err := p.createPhoneNumbersTable(); err != nil {
		return err
	}

	for _, num := range exampleNumbers {
		if _, err := p.InsertPhone(num); err != nil {
			return err
		}
	}

	return nil
}

func (p *PostgresDB) Normalize() error {
	phones, err := p.AllPhones()
	if err != nil {
		return err
	}

	for _, phone := range phones {
		number := reg.ReplaceAllString(phone.Number, "")
		if number != phone.Number {
			existing, err := p.FindPhone(number)
			if err != nil {
				return err
			}

			if existing != nil {
				if err := p.DeletePhone(phone.Id); err != nil {
					return err
				}
			} else {
				phone.Number = number
				if err := p.UpdatePhone(phone); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (p *PostgresDB) createPhoneNumbersTable() error {
	query := `
    CREATE TABLE IF NOT EXISTS phone_numbers (
      id SERIAL,
      value VARCHAR(255)
    )`

	_, err := p.db.Exec(query)
	return err
}

func (p *PostgresDB) InsertPhone(phone string) (int, error) {
	query := `INSERT INTO phone_numbers(value) VALUES($1) RETURNING id`

	var id int
	if err := p.db.QueryRow(query, phone).Scan(&id); err != nil {
		return -1, err
	}

	return id, nil
}

func (p *PostgresDB) AllPhones() ([]phone.Phone, error) {
	query := "SELECT id, value FROM phone_numbers"
	rows, err := p.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ps []phone.Phone
	for rows.Next() {
		var p phone.Phone
		if err := rows.Scan(&p.Id, &p.Number); err != nil {
			return nil, err
		}

		ps = append(ps, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ps, nil
}

func (p *PostgresDB) FindPhone(number string) (*phone.Phone, error) {
	query := "SELECT * FROM phone_numbers WHERE value=$1"
	row := p.db.QueryRow(query, number)

	var phone phone.Phone
	if err := row.Scan(&phone.Id, &phone.Number); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &phone, nil
}

func (p *PostgresDB) GetPhone(id int) (string, error) {
	query := "SELECT * FROM phone_numbers WHERE id=$1"
	row := p.db.QueryRow(query, id)

	var number string
	if err := row.Scan(&id, &number); err != nil {
		return "", err
	}

	return number, nil
}

func (p *PostgresDB) UpdatePhone(phone phone.Phone) error {
	query := `UPDATE phone_numbers SET value=$2 WHERE id=$1`
	_, err := p.db.Exec(query, phone.Id, phone.Number)
	return err
}

func (p *PostgresDB) DeletePhone(id int) error {
	query := `DELETE FROM phone_numbers WHERE id=$1`
	_, err := p.db.Exec(query, id)
	return err
}
