package services

import (
	"diplomIshi/internal/models"
	"diplomIshi/internal/repository"
	"github.com/robfig/cron/v3"
	"log"
)

type ReminderService struct {
	repo    *repository.UserRepository
	cron    *cron.Cron
	entries map[uint]cron.EntryID // Track cron entries by reminder ID
}

func NewReminderService(repo *repository.UserRepository) *ReminderService {
	return &ReminderService{
		repo:    repo,
		cron:    cron.New(),
		entries: make(map[uint]cron.EntryID),
	}
}

func (s *ReminderService) Start() {
	s.cron.Start()
}

func (s *ReminderService) Stop() {
	s.cron.Stop()
}

func (s *ReminderService) ScheduleReminder(reminder *models.Reminder) error {
	entryID, err := s.cron.AddFunc(reminder.Time, func() {
		// Simulate sending a push notification (replace with real implementation)
		log.Printf("Reminder for User %d: %s", reminder.UserID, reminder.Message)
		// Example: Integrate with FCM here
	})
	if err != nil {
		return err
	}
	s.entries[reminder.ID] = entryID
	return nil
}

func (s *ReminderService) LoadReminders(userID uint) error {
	reminders, err := s.repo.GetReminders(userID)
	if err != nil {
		return err
	}
	for i := range reminders {
		if err := s.ScheduleReminder(&reminders[i]); err != nil {
			log.Printf("Failed to schedule reminder %d: %v", reminders[i].ID, err)
		}
	}
	return nil
}

func (s *ReminderService) RemoveReminder(reminderID uint) {
	if entryID, exists := s.entries[reminderID]; exists {
		s.cron.Remove(entryID)
		delete(s.entries, reminderID)
	}
}
