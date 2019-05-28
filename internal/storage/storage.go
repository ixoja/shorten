package storage

import (
	"database/sql"
	"encoding/base64"
	"github.com/pkg/errors"
	"log"
	"strconv"
	"time"

	"github.com/ixoja/shorten/internal/model"
)

type SQLite struct {
	Database *sql.DB
}

func New(db *sql.DB) *SQLite {
	return &SQLite{Database: db}
}

const (
	createTable = "CREATE TABLE IF NOT EXISTS shorten " +
		"(id INTEGER PRIMARY KEY AUTOINCREMENT, longURL TEXT, lastAccess INTEGER)"
	insert   = "INSERT INTO shorten (longURL, lastAccess) VALUES (?,?)"
	get      = "SELECT id, longURL, lastAccess FROM shorten WHERE id=?"
	getByURL = "SELECT id, longURL, lastAccess FROM shorten WHERE longURL=?"
)

func (s *SQLite) InitDB() error {
	stmt, err := s.Database.Prepare(createTable)
	if err != nil {
		return errors.Wrap(err, "failed to prepare statement to init table")
	}

	_, err = stmt.Exec()
	if err != nil {
		return errors.Wrap(err, "failed to init table")
	}
	return nil
}

func (s *SQLite) Save(stored *model.StoredURL) (*model.StoredURL, error) {
	stmt, err := s.Database.Prepare(insert)
	if err != nil {
		return nil, errors.Wrap(err, "failed to prepare statement to save url")
	}

	stored.LastAccess = time.Now()
	res, err := stmt.Exec(stored.LongURL, time.Now().Unix())
	if err != nil {
		return nil, errors.Wrap(err, "failed to save url")
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get last insert id")
	}
	stored.ID = base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(int(id))))

	return stored, nil
}
func (s *SQLite) Delete(id string) error {
	return nil
}
func (s *SQLite) Get(id string) (*model.StoredURL, bool, error) {
	stmt, err := s.Database.Prepare(get)
	if err != nil {
		return nil, false, errors.Wrap(err, "failed to prepare statement to get url")
	}

	idDec, err := base64.StdEncoding.DecodeString(id)
	if err != nil {
		return nil, false, errors.Wrap(err, "failed to decode id")
	}

	res, err := stmt.Query(string(idDec))
	if err != nil {
		return nil, false, errors.Wrap(err, "failed to get url")
	}
	defer func() {
		err := res.Close()
		if err != nil {
			log.Println("failed to close prepared statement", err.Error())
		}
	}()

	if hasNext := res.Next(); hasNext {
		stored := &model.StoredURL{}
		var id string
		var timestamp int64
		err := res.Scan(&id, &stored.LongURL, &timestamp)
		if err != nil {
			log.Println("failed to close prepared statement", err.Error())
		}
		stored.ID = base64.StdEncoding.EncodeToString([]byte(id))
		stored.LastAccess = time.Unix(timestamp, 0)
		return stored, true, nil
	}

	return nil, false, nil
}
func (s *SQLite) GetByURL(longURL string) (*model.StoredURL, bool, error) {
	stmt, err := s.Database.Prepare(getByURL)
	if err != nil {
		return nil, false, errors.Wrap(err, "failed to prepare statement to save url")
	}

	res, err := stmt.Query(longURL)
	if err != nil {
		return nil, false, errors.Wrap(err, "failed to save url")
	}
	defer func() {
		err := res.Close()
		if err != nil {
			log.Println("failed to close prepared statement", err.Error())
		}
	}()

	if hasNext := res.Next(); hasNext {
		stored := &model.StoredURL{}
		var id string
		var timestamp int64
		err := res.Scan(&id, &stored.LongURL, &timestamp)
		if err != nil {
			log.Println("failed to close prepared statement", err.Error())
		}
		stored.ID = base64.StdEncoding.EncodeToString([]byte(id))
		stored.LastAccess = time.Unix(timestamp, 0)
		return stored, true, nil
	}

	return nil, false, nil
}
func (s *SQLite) EvictOlder(timestamp time.Time) error {
	return nil
}
