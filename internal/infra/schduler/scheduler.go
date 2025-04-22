package scheduler

import (
	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"gorm.io/gorm"
)

type Scheduler struct {
	db   *gorm.DB
	cron *cron.Cron
}

type IScheduler interface {
	InitScheduler()
	Start()
	Stop()
}

func NewScheduler(db *gorm.DB) IScheduler {
	return &Scheduler{
		db:   db,
		cron: cron.New(),
	}
}

func (s *Scheduler) InitScheduler() {
	s.cron.AddFunc("0 0 * * *", func() {
		s.db.Transaction(func(tx *gorm.DB) error {
			s.db = tx
			err := s.CreateDailyWorkStaff()
			if err != nil {
				return err
			}
			return nil
		})
	})
}

func (s *Scheduler) CreateDailyWorkStaff() error {
	var dailyWorks []entity.DailyWork
	if err := s.db.Find(&dailyWorks).Error; err != nil {
		return err
	}

	var staffs []entity.Staff
	if err := s.db.Preload("Account").Find(&staffs).Error; err != nil {
		return err
	}

	var dailyWorkStaffs []entity.DailyWorkStaff
	for _, dailyWork := range dailyWorks {
		for _, staff := range staffs {
			if staff.Account.Role.Id == dailyWork.RoleId {
				dailyWorkStaff := entity.DailyWorkStaff{
					DailyWorkId: dailyWork.Id,
					// StaffId:     staff.Id,
					IsDone:    false,
					CreatedBy: uuid.Nil,
				}
				dailyWorkStaffs = append(dailyWorkStaffs, dailyWorkStaff)
			}
		}
	}

	return s.db.CreateInBatches(dailyWorkStaffs, len(dailyWorkStaffs)).Error
}

func (s *Scheduler) Start() {
	s.cron.Start()
}

func (s *Scheduler) Stop() {
	s.cron.Stop()
}
