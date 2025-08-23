package postgres

import (
	"back/back/urlShortner/internal/config"
	"back/back/urlShortner/internal/config/storage"
	

	"database/sql"
	"fmt"

	"github.com/lib/pq" // init  postgres driver
)

const codeUniqueValueError = "23505"  //code 


type Storage struct{
	db *sql.DB
}

func New(cfg config.DBConfig) (*Storage, error){
	const op = "storage.postgres.New"

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)

	if err!=nil{
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS url (
		id SERIAL PRIMARY KEY,
		alias TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL
	)`)
	if err!=nil{
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_alias ON url(alias)`) //quick search for url by index
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}


	return &Storage{db:db}, nil
}


func (s *Storage) SaveURL(urlToSave string, alias string) error{
	const op = "storage.postgres.SaveURL"

	_, err := s.db.Exec("INSERT INTO url(url, alias) VALUES($1, $2)",urlToSave,alias)
	
	if err!=nil{
		if postgresErr, ok := err.(*pq.Error); ok && postgresErr.Code == codeUniqueValueError{
			return fmt.Errorf("%s: %w", op, storage.ErrURLExists)
		}
		return fmt.Errorf("%s: %w", op, err)
	}	

	
	return nil
}

func (s *Storage) GetURL(alias string) (string, error){
	const op = "storage.postgres.GetURL"

	var urlToGet string

	err := s.db.QueryRow("SELECT url FROM url WHERE alias = $1", alias).Scan(&urlToGet)

	if err!=nil{
		if err== sql.ErrNoRows{
			return "", fmt.Errorf("%s: %w", op, storage.ErrURLNotFound)
		}
		return "",fmt.Errorf("%s: %w",op, err)
	}

	return urlToGet, nil
}

func (s *Storage) DeleteURL(alias string) error{
	const op = "storage.postgres.DeleteURL"

	res, err := s.db.Exec("DELETE FROM url WHERE alias = $1", alias)

	if err!=nil{
		return fmt.Errorf("%s: %w",op, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err!=nil{
		return fmt.Errorf("%s: %w",op, err)
	}

	if rowsAffected==0{
		return fmt.Errorf("%s: %w", op, storage.ErrURLNotFound)
	}
	return nil
}