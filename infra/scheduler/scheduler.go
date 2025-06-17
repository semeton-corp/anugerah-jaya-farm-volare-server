package scheduler

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	datatype "github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/custom/data_type"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	ownerRole = "Owner"
)

type IScheduler interface {
	Start()
	InitScheduler()
	Stop()
}

type Scheduler struct {
	db   *gorm.DB
	cron *cron.Cron
	log  *zap.Logger
}

func New(db *gorm.DB, log *zap.Logger) IScheduler {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to load timezone: %v", err))
	}

	return &Scheduler{
		db:   db,
		cron: cron.New(cron.WithLocation(loc), cron.WithChain(cron.Recover(cron.DefaultLogger))),
		log:  log,
	}
}

func (s *Scheduler) InitScheduler() {
	s.cron.AddFunc("01 00 * * *", func() {
		s.db.Transaction(func(tx *gorm.DB) error {
			err := s.createDailyWorkStaff(tx)
			if err != nil {
				s.log.Error("failed to create daily work staff", zap.Error(err))
				return err
			}
			return nil
		})
	})

	s.cron.AddFunc("01 00 * * *", func() {
		s.db.Transaction(func(tx *gorm.DB) error {
			err := s.createStaffPresence(tx)
			if err != nil {
				s.log.Error("failed to create staff presence", zap.Error(err))
				return err
			}
			return nil
		})
	})

	s.cron.AddFunc("01 00 * * *", func() {
		s.db.Transaction(func(tx *gorm.DB) error {
			err := s.checkForgottenStaffPresence(tx)
			if err != nil {
				s.log.Error("failed to check staff presence", zap.Error(err))
				return err
			}
			return nil
		})
	})
}

func (s *Scheduler) createDailyWorkStaff(tx *gorm.DB) error {
	s.log.Info("Creating daily work staff...")

	var dailyWorks []entity.DailyWork
	if err := tx.Model(&entity.DailyWork{}).Find(&dailyWorks).Error; err != nil {
		return err
	}

	var users []entity.User
	if err := tx.Preload("Role").Find(&users).Error; err != nil {
		return err
	}

	var dailyWorkStaffs []entity.DailyWorkStaff
	for _, dailyWork := range dailyWorks {
		for _, user := range users {
			if user.RoleId == dailyWork.RoleId {
				dailyWorkStaff := entity.DailyWorkStaff{
					DailyWorkId: dailyWork.Id,
					UserId:      user.Id,
					IsDone:      false,
					CreatedBy:   uuid.NullUUID{UUID: uuid.Nil, Valid: false},
				}
				dailyWorkStaffs = append(dailyWorkStaffs, dailyWorkStaff)
			}
		}
	}

	s.log.Info(fmt.Sprintf("Daily work staff created: %d", len(dailyWorkStaffs)))
	return tx.CreateInBatches(dailyWorkStaffs, len(dailyWorkStaffs)).Error
}

func (s *Scheduler) createStaffPresence(tx *gorm.DB) error {
	s.log.Info("Creating staff presence...")

	var users []entity.User
	tx.Model(&entity.User{}).Preload("Role").Find(&users)

	for _, user := range users {
		if user.Role.Name != ownerRole {
			staffPresence := entity.StaffPresence{
				UserId:    user.Id,
				IsPresent: false,
				CreatedBy: uuid.NullUUID{UUID: uuid.Nil, Valid: false},
			}
			if err := tx.Create(&staffPresence).Error; err != nil {
				return err
			}
		}
	}

	s.log.Info(fmt.Sprintf("Staff presence created: %d", len(users)))
	return nil
}

func (s *Scheduler) checkForgottenStaffPresence(tx *gorm.DB) error {
	s.log.Info("Checking staff presence...")

	var staffPresences []entity.StaffPresence
	if err := tx.Where("is_present = ? AND start_time IS NOT NULL", false).Find(&staffPresences).Error; err != nil {
		return err
	}

	endTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 17, 0, 0, 0, time.UTC)

	for _, staffPresence := range staffPresences {
		staffPresence.IsPresent = true
		staffPresence.EndTime = datatype.TimeOnly{Time: endTime}
		if err := tx.Save(&staffPresence).Error; err != nil {
			return err
		}
	}

	s.log.Info(fmt.Sprintf("Staff presence checked: %d", len(staffPresences)))
	return nil
}

func (s *Scheduler) Start() {
	s.cron.Start()
}

func (s *Scheduler) Stop() {
	s.cron.Stop()
}
