package service

import (
	"errors"
	"time"

	"specialdates-backend/internal/models"
	"specialdates-backend/internal/repository"
)

type DateService interface {
	CreateDate(d *models.DateEvent) (int, error)
	ListDates() ([]models.DateEvent, error)
	GetDate(id int) (*models.DateEvent, error)
	UpdateDate(d *models.DateEvent) error
	DeleteDate(id int) error
}

type dateService struct {
	repo repository.DateRepository
}

func NewDateService(r repository.DateRepository) DateService {
	return &dateService{repo: r}
}

func (s *dateService) CreateDate(d *models.DateEvent) (int, error) {
	// Validaciones de negocio
	if d.Title == "" {
		return 0, errors.New("title is required")
	}
	if d.EventAt.IsZero() {
		return 0, errors.New("event date is required")
	}
	// ejemplo de regla: no permitir fechas en el pasado
	if d.EventAt.Before(time.Now().Add(-time.Minute)) {
		return 0, errors.New("event date cannot be in the past")
	}
	return s.repo.CreateDate(d)
}

func (s *dateService) ListDates() ([]models.DateEvent, error) {
	return s.repo.ListDates()
}

func (s *dateService) GetDate(id int) (*models.DateEvent, error) {
	return s.repo.GetDate(id)
}

func (s *dateService) UpdateDate(d *models.DateEvent) error {
	if d.Title == "" {
		return errors.New("title is required")
	}
	return s.repo.UpdateDate(d)
}

func (s *dateService) DeleteDate(id int) error {
	return s.repo.DeleteDate(id)
}
