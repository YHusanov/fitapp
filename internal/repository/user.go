package repository

import (
	"diplomIshi/internal/models"
	"gorm.io/gorm"
	"time"
)

type UserRepository struct {
	Db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{Db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	return r.Db.Create(user).Error
}

func (r *UserRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := r.Db.First(&user, id).Error
	return &user, err
}

func (r *UserRepository) Update(user *models.User) error {
	return r.Db.Save(user).Error
}

func (r *UserRepository) AddPoint(point *models.Point) error {
	return r.Db.Create(point).Error
}

func (r *UserRepository) AddAchievement(achievement *models.Achievement) error {
	return r.Db.Create(achievement).Error
}

func (r *UserRepository) GetPoints(userID uint) ([]models.Point, error) {
	var points []models.Point
	err := r.Db.Where("user_id = ?", userID).Find(&points).Error
	return points, err
}

func (r *UserRepository) GetAchievements(userID uint) ([]models.Achievement, error) {
	var achievements []models.Achievement
	err := r.Db.Where("user_id = ?", userID).Find(&achievements).Error
	return achievements, err
}

func (r *UserRepository) GetRecentExerciseDays(userID uint) ([]time.Time, error) {
	var progress []models.Progress
	err := r.Db.Where("user_id = ?", userID).Order("date desc").Limit(7).Find(&progress).Error
	if err != nil {
		return nil, err
	}
	var days []time.Time
	for _, p := range progress {
		days = append(days, p.Date)
	}
	return days, nil
}

func (r *UserRepository) CreateReminder(reminder *models.Reminder) error {
	return r.Db.Create(reminder).Error
}

func (r *UserRepository) GetReminders(userID uint) ([]models.Reminder, error) {
	var reminders []models.Reminder
	err := r.Db.Where("user_id = ? AND active = ?", userID, true).Find(&reminders).Error
	return reminders, err
}

func (r *UserRepository) UpdateReminder(reminder *models.Reminder) error {
	return r.Db.Save(reminder).Error
}

func (r *UserRepository) CreatePost(post *models.Post) error {
	return r.Db.Create(post).Error
}

func (r *UserRepository) GetPosts() ([]models.Post, error) {
	var posts []models.Post
	err := r.Db.Order("created_at desc").Find(&posts).Error
	return posts, err
}

func (r *UserRepository) GetPostByID(postID uint) (*models.Post, error) {
	var post models.Post
	err := r.Db.First(&post, postID).Error
	return &post, err
}

func (r *UserRepository) CreateComment(comment *models.Comment) error {
	return r.Db.Create(comment).Error
}

func (r *UserRepository) GetComments(postID uint) ([]models.Comment, error) {
	var comments []models.Comment
	err := r.Db.Where("post_id = ?", postID).Order("created_at asc").Find(&comments).Error
	return comments, err
}

func (r *UserRepository) DeleteOldPlans(userID uint) error {
	if err := r.Db.Where("user_id = ?", userID).Delete(&models.Meal{}).Error; err != nil {
		return err
	}
	return r.Db.Where("user_id = ?", userID).Delete(&models.Exercise{}).Error
}
