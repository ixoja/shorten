package storage

import (
	"database/sql"
	"github.com/icrowley/fake"
	"github.com/ixoja/shorten/internal/model"
	"github.com/stretchr/testify/require"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func dbConnection(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		t.Error(err)
	}
	return db
}

func TestSQLite_InitDB(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db := dbConnection(t)
		s := SQLite{*db}
		err := s.InitDB()
		assert.NoError(t, err)
		db.Close()
	})
}

func TestSQLite_Save(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db := dbConnection(t)
		s := SQLite{*db}
		err := s.InitDB()
		require.NoError(t, err)
		longURL := fake.DomainName()
		stored := &model.StoredURL{LongURL: longURL}
		stored, err = s.Save(stored)
		require.NoError(t, err)
		newID := stored.ID
		assert.NotEqual(t, "", newID)
		stored, ok, err := s.Get(stored.ID)
		require.NoError(t, err)
		assert.True(t, ok)
		assert.Equal(t, longURL, stored.LongURL)
		assert.Equal(t, newID, stored.ID)
		db.Close()
	})
}

func TestSQLite_GetByURL(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db := dbConnection(t)
		s := SQLite{*db}
		err := s.InitDB()
		require.NoError(t, err)
		longURL := fake.DomainName()
		stored := &model.StoredURL{LongURL: longURL}
		stored, err = s.Save(stored)
		require.NoError(t, err)
		newID := stored.ID
		assert.NotEqual(t, "", newID)
		stored, ok, err := s.GetByURL(stored.LongURL)
		require.NoError(t, err)
		assert.True(t, ok)
		assert.Equal(t, longURL, stored.LongURL)
		assert.Equal(t, newID, stored.ID)
		db.Close()
	})
}
