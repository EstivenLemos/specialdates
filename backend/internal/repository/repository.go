package repository

import (
	"database/sql"

	"specialdates-backend/internal/models"
)

type DateRepository interface {
	CreateDate(d *models.DateEvent) (int, error)
	ListDates() ([]models.DateEvent, error)
	GetDate(id int) (*models.DateEvent, error)
	UpdateDate(d *models.DateEvent) error
	DeleteDate(id int) error
}

type MySQLRepo struct {
	db *sql.DB
}

func NewMySQLRepository(db *sql.DB) *MySQLRepo {
	return &MySQLRepo{db: db}
}

func (r *MySQLRepo) CreateDate(d *models.DateEvent) (int, error) {
	res, err := r.db.Exec(
		`INSERT INTO dates (title, description, event_datetime, date_type, recurring) VALUES (?,?,?,?,?)`,
		d.Title, d.Description, d.EventAt, d.DateType, boolToTinyint(d.Recurring),
	)
	if err != nil {
		return 0, err
	}
	last, err := res.LastInsertId()
	return int(last), err
}

func (r *MySQLRepo) ListDates() ([]models.DateEvent, error) {
	rows, err := r.db.Query(`SELECT id, title, description, event_datetime, date_type, recurring, created_at FROM dates ORDER BY event_datetime ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.DateEvent
	for rows.Next() {
		var d models.DateEvent
		var recurring tinyInt
		if err := rows.Scan(&d.ID, &d.Title, &d.Description, &d.EventAt, &d.DateType, &recurring, &d.CreatedAt); err != nil {
			return nil, err
		}
		d.Recurring = tinyIntToBool(recurring)
		out = append(out, d)
	}
	return out, nil
}

func (r *MySQLRepo) GetDate(id int) (*models.DateEvent, error) {
	var d models.DateEvent
	var recurring tinyInt
	err := r.db.QueryRow(`SELECT id, title, description, event_datetime, date_type, recurring, created_at FROM dates WHERE id = ?`, id).
		Scan(&d.ID, &d.Title, &d.Description, &d.EventAt, &d.DateType, &recurring, &d.CreatedAt)
	if err != nil {
		return nil, err
	}
	d.Recurring = tinyIntToBool(recurring)
	return &d, nil
}

func (r *MySQLRepo) UpdateDate(d *models.DateEvent) error {
	_, err := r.db.Exec(
		`UPDATE dates SET title=?, description=?, event_datetime=?, date_type=?, recurring=? WHERE id=?`,
		d.Title, d.Description, d.EventAt, d.DateType, boolToTinyint(d.Recurring), d.ID,
	)
	return err
}

func (r *MySQLRepo) DeleteDate(id int) error {
	_, err := r.db.Exec(`DELETE FROM dates WHERE id = ?`, id)
	return err
}

// helpers
type tinyInt int

func boolToTinyint(b bool) int {
	if b {
		return 1
	}
	return 0
}

func tinyIntToBool(t tinyInt) bool {
	return t != 0
}
