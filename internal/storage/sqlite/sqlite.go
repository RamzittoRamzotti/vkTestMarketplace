package sqlite

import (
	"database/sql"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"vkTestMarketplace/internal/models"
	"vkTestMarketplace/internal/storage"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	op := "storage.sqlite.New"
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s : %s", op, err)
	}
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY,
    login VARCHAR(64) UNIQUE,
    password_hash TEXT
	);
	CREATE TABLE IF NOT EXISTS ads (
		id INTEGER PRIMARY KEY,
		title VARCHAR(100),
		description TEXT,
		image_url TEXT,
		price INTEGER,
		author_id INTEGER REFERENCES users(id),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

    
`)
	if err != nil {
		return nil, fmt.Errorf("%s : %s", op, err)
	}
	return &Storage{db: db}, nil
}

func (s *Storage) CreateUser(user *models.User) (int64, error) {
	stmt, err := s.db.Prepare("INSERT INTO users(login, password_hash) VALUES(?, ?)")
	if err != nil {
		return 0, err
	}
	res, er := stmt.Exec(user.Login, user.PasswordHash)
	if er != nil {
		if sqliteErr, ok := er.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, storage.ErrUserExists
		}
		return 0, er
	}
	return res.LastInsertId()
}

func (s *Storage) GetUserByUsername(login string) (*models.User, error) {
	row := s.db.QueryRow("SELECT id, login, password_hash FROM users WHERE login = ?", login)
	var user models.User
	if err := row.Scan(&user.ID, &user.Login, &user.PasswordHash); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (s *Storage) GetUserByID(id int) (*models.User, error) {
	row := s.db.QueryRow("SELECT id, login, password_hash FROM users WHERE id = ?", id)
	var user models.User
	if err := row.Scan(&user.ID, &user.Login, &user.PasswordHash); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (s *Storage) CreateAd(ad *models.Ad) error {
	stmt, err := s.db.Prepare(`INSERT INTO ads (title, description, image_url, price, author_id) VALUES (?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	res, err := stmt.Exec(ad.Title, ad.Description, ad.ImageURL, ad.Price, ad.AuthorID)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	ad.ID = int(id)
	return nil
}

func (s *Storage) ListAds(filter storage.AdFilter) ([]models.Ad, error) {
	query := `SELECT ads.id, ads.title, ads.description, ads.image_url, ads.price, ads.author_id, users.login
	FROM ads
	JOIN users ON users.id = ads.author_id
	`
	args := []interface{}{}
	if filter.MinPrice > 0 {
		query += " AND ads.price >= ?"
		args = append(args, filter.MinPrice)
	}
	if filter.MaxPrice > 0 {
		query += " AND ads.price <= ?"
		args = append(args, filter.MaxPrice)
	}

	sortBy := "ads.created_at"
	if filter.SortBy == "price" {
		sortBy = "ads.price"
	}

	sortOrder := "DESC"
	if filter.SortOrder == "asc" {
		sortOrder = "ASC"
	}

	query += " ORDER BY " + sortBy + " " + sortOrder
	limit := 10
	if filter.Limit > 0 {
		limit = filter.Limit
	}
	offset := 0
	if filter.Page > 1 {
		offset = (filter.Page - 1) * limit
	}
	query += " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ads []models.Ad
	for rows.Next() {
		var ad models.Ad
		if err := rows.Scan(&ad.ID, &ad.Title, &ad.Description, &ad.ImageURL, &ad.Price, &ad.AuthorID, &ad.AuthorUsername); err != nil {
			return nil, err
		}
		ads = append(ads, ad)
	}
	return ads, nil
}
