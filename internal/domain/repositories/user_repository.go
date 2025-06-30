package repositories

import (
	"github.com/google/uuid"
	"github.com/kitamersion/go-goservice/internal/domain/entities"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *entities.User) error
	GetByID(id uuid.UUID) (*entities.User, error)
	GetByEmail(email string) (*entities.User, error)
	Update(user *entities.User) error
	Delete(id uuid.UUID) error
	List(limit, offset int) ([]*entities.User, error)
	Count() (int64, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(user *entities.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) GetByID(id uuid.UUID) (*entities.User, error) {
	var user entities.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(email string) (*entities.User, error) {
	var user entities.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *entities.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entities.User{}, "id = ?", id).Error
}

func (r *userRepository) List(limit, offset int) ([]*entities.User, error) {
	var users []*entities.User
	err := r.db.Limit(limit).Offset(offset).Find(&users).Error
	return users, err
}

func (r *userRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&entities.User{}).Count(&count).Error
	return count, err
}
