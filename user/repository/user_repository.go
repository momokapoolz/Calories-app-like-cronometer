package repository

import (
	"log"

	"github.com/momokapoolz/caloriesapp/user/database"
	"github.com/momokapoolz/caloriesapp/user/models"
)

// UserRepository handles database operations related to users
type UserRepository struct{}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

// Create saves a new user to the database
func (r *UserRepository) Create(user *models.User) error {
	return database.DB.Create(user).Error
}

// FindByID retrieves a user by ID
func (r *UserRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := database.DB.First(&user, id).Error
	return &user, err
}

// FindByEmail retrieves a user by email
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	log.Printf("[UserRepository] Finding user by email: %s", email)

	// Enable SQL query logging for this operation
	tx := database.DB.Debug().Where("email = ?", email).First(&user)
	if tx.Error != nil {
		log.Printf("[UserRepository] Error finding user: %v", tx.Error)
		return nil, tx.Error
	}

	log.Printf("[UserRepository] Found user: ID=%d, Email=%s, PasswordHash=%s", user.ID, user.Email, user.PasswordHash)
	return &user, nil
}

// FindByRole retrieves users by role
func (r *UserRepository) FindByRole(role string) ([]models.User, error) {
	var users []models.User
	err := database.DB.Where("role = ?", role).Find(&users).Error
	return users, err
}

// Update updates an existing user
func (r *UserRepository) Update(user *models.User) error {
	return database.DB.Save(user).Error
}

// Delete removes a user from the database
func (r *UserRepository) Delete(id uint) error {
	return database.DB.Delete(&models.User{}, id).Error
}

// FindAll retrieves all users
func (r *UserRepository) FindAll() ([]models.User, error) {
	var users []models.User
	err := database.DB.Find(&users).Error
	return users, err
}
