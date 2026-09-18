package server

import (
	"database/sql"
	"errors"
	"short2/internal/data"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Server struct {
	db  *sql.DB
	cfg *data.ConfigData
}

func NewServer(cfg *data.ConfigData, db *sql.DB) *Server {
	return &Server{db, cfg}
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

func (s Server) GetUserIdByTokenAndAdvance(token string) (string, error) {
	var userID string
	var expiryDt time.Time
	var expiredInd string

	err := s.db.QueryRow("select person_id, expiry_dt, expired_ind from session where token = ?", token).Scan(&userID, &expiryDt, &expiredInd)
	if err != nil {
		return "", err
	}

	if time.Now().After(expiryDt) || strings.ToUpper(expiredInd) == "Y" {
		return "", errors.New("token expired")
	}

	// advance the time
	newTime := time.Now().Add(time.Duration(s.cfg.TokenTtlSecond) * time.Second)
	_, err = s.db.Exec("update session set expiry_dt = ? where token = ?", newTime, token)
	if err != nil {
		return "", err
	}

	return userID, nil
}

func (s Server) GetUserById(userId string) (*data.Person, error) {
	rows, err := s.db.Query("select person_id, login, full_name, provided_id from person where person_id = ?", userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var person data.Person
	if rows.Next() {
		err = rows.Scan(&person.ID, &person.Login, &person.FullName, &person.ProvidedID)

		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("user with this id not found")
	}

	return &person, nil
}

func (s Server) ExpireSessionByToken(token string) error {
	_, err := s.db.Exec("update session set expiry_dt = 'Y' where token = ?", token)

	return err
}

func (s Server) GetUrlByShortID(shortID string) (string, error) {
	var url string
	err := s.db.QueryRow("select full_url from urls where shortened_code = ?", shortID).Scan(&url)

	return url, err
}

func (s Server) ShortenAndSave(userID string, longUrl string) (string, error) {
	// basically try about 10 times to find a random guid, take the first 8 characters from it and see if exists
	// if not, this is it, if it does not work, will then report an error and then the user will be able to try again

	var shortened string
	errorSituation := true

	for i := 0; i < 10; i++ {
		shortened = uuid.Must(uuid.NewRandom()).String()
		runes := []rune(shortened)
		runes = runes[:8]

		shortened = string(runes)

		rs, err := s.db.Query("select url_id from urls where shortened_code = ?", string(runes))
		if err != nil {
			return "", err
		}

		found := rs.Next()
		rs.Close()

		if !found {
			errorSituation = false
			break
		}
	}

	if errorSituation {
		return "", errors.New("all the attempts are already in the database")
	}

	urlId := uuid.Must(uuid.NewV7()).String()
	_, err := s.db.Exec("insert into urls (url_id, person_id, shortened_code, full_url) values (?, ?, ?, ?)", urlId, userID, shortened, longUrl)
	if err != nil {
		return "", err
	}

	return shortened, nil
}
