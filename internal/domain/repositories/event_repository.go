package repositories

import (
	"github.com/google/uuid"
	"github.com/kitamersion/go-goservice/internal/domain/entities"
	"gorm.io/gorm"
)

type EventRepository interface {
	Create(event *entities.Event) error
	GetByID(id uuid.UUID) (*entities.Event, error)
	GetByType(eventType string, limit, offset int) ([]*entities.Event, error)
	List(limit, offset int) ([]*entities.Event, error)
	Count() (int64, error)
}

type eventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) EventRepository {
	return &eventRepository{
		db: db,
	}
}

func (r *eventRepository) Create(event *entities.Event) error {
	return r.db.Create(event).Error
}

func (r *eventRepository) GetByID(id uuid.UUID) (*entities.Event, error) {
	var event entities.Event
	err := r.db.Where("id = ?", id).First(&event).Error
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *eventRepository) GetByType(eventType string, limit, offset int) ([]*entities.Event, error) {
	var events []*entities.Event
	err := r.db.Where("type = ?", eventType).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&events).Error
	return events, err
}

func (r *eventRepository) List(limit, offset int) ([]*entities.Event, error) {
	var events []*entities.Event
	err := r.db.Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&events).Error
	return events, err
}

func (r *eventRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&entities.Event{}).Count(&count).Error
	return count, err
}
