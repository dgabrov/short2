package server

import (
	"database/sql"
	"short2/internal/data"
	"time"

	"github.com/google/uuid"
)

type Server struct {
	db  *sql.DB
	cfg *data.ConfigData
}

func (s Server) ExistsProvidedUserId(id string) (bool, error) {
	rows, err := s.db.Query("select person_id from person where provided_id = ?", id)
	if err != nil {
		return false, err
	}

	res := rows.Next()

	defer rows.Close()

	return res, nil
}

func (s Server) CreateSessionData(userId string, expiryDate time.Time, token string) error {
	_, err := s.db.Exec("insert into session (session_id, person_id, expiry_dt, token) values (?,?,?,?)",
		uuid.Must(uuid.NewV7()).String(), userId, expiryDate, token)

	return err
}

func (s Server) GetUserIdByProvidedId(providedId string) (string, error) {
	res := ""

	rows, err := s.db.Query("select person_id from person where provided_id = ?", providedId)
	if err != nil {
		return res, err
	}

	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&res)

		if err != nil {
			return res, err
		}
	}

	return res, nil
}

func (s Server) CreateUser(userId string, login string, name string, providedId string) error {
	_, err := s.db.Exec("insert into person (person_id, login, full_name, provided_id) values (?,?,?,?)",
		userId, login, name, providedId)

	return err
}

func NewServer(cfg *data.ConfigData, db *sql.DB) *Server {
	return &Server{db, cfg}
}
