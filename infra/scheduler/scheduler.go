package scheduler

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/constant"
	datatype "github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/custom/data_type"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
	"github.com/shopspring/decimal"
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
			err := s.createDailyWorkUser(tx)
			if err != nil {
				s.log.Error("failed to create daily work user", zap.Error(err))
				return err
			}
			return nil
		})
	})

	s.cron.AddFunc("01 00 * * *", func() {
		s.db.Transaction(func(tx *gorm.DB) error {
			err := s.createUserPresence(tx)
			if err != nil {
				s.log.Error("failed to create user presence", zap.Error(err))
				return err
			}
			return nil
		})
	})

	s.cron.AddFunc("01 00 * * *", func() {
		s.db.Transaction(func(tx *gorm.DB) error {
			err := s.checkForgottenUserPresence(tx)
			if err != nil {
				s.log.Error("failed to check user presence", zap.Error(err))
				return err
			}
			return nil
		})
	})

	s.cron.AddFunc("0 0 1 * *", func() {
		s.db.Transaction(func(tx *gorm.DB) error {
			err := s.createUserSalaryPaymentPerMonth(tx)
			if err != nil {
				s.log.Error("failed to create user salary payment per month", zap.Error(err))
				return err
			}
			return nil
		})
	})
}

func (s *Scheduler) createDailyWorkUser(tx *gorm.DB) error {
	s.log.Info("Creating daily work user...")

	var dailyWorks []entity.DailyWork
	if err := tx.Model(&entity.DailyWork{}).Find(&dailyWorks).Error; err != nil {
		return err
	}

	var users []entity.User
	if err := tx.Preload("Role").Find(&users).Error; err != nil {
		return err
	}

	var dailyWorkUsers []entity.DailyWorkUser
	for _, dailyWork := range dailyWorks {
		for _, user := range users {
			if user.RoleId == dailyWork.RoleId {
				dailyWorkUser := entity.DailyWorkUser{
					DailyWorkId: dailyWork.Id,
					UserId:      user.Id,
					IsDone:      false,
					CreatedBy:   uuid.NullUUID{UUID: uuid.Nil, Valid: false},
				}
				dailyWorkUsers = append(dailyWorkUsers, dailyWorkUser)
			}
		}
	}

	s.log.Info(fmt.Sprintf("Daily work user created: %d", len(dailyWorkUsers)))
	return tx.CreateInBatches(dailyWorkUsers, len(dailyWorkUsers)).Error
}

func (s *Scheduler) createUserPresence(tx *gorm.DB) error {
	s.log.Info("Creating user presence...")

	var users []entity.User
	tx.Model(&entity.User{}).Preload("Role").Find(&users)

	totalCreatedUserPresence := 0
	for _, user := range users {
		if user.Role.Name != ownerRole {
			userPresence := entity.UserPresence{
				UserId:                   user.Id,
				Status:                   enum.PresenceStatusAlpha,
				SubmissionPresenceStatus: enum.SubmissionPresenceStatusNoSubmission,
				CreatedBy:                uuid.NullUUID{UUID: uuid.Nil, Valid: false},
			}
			if err := tx.Create(&userPresence).Error; err != nil {
				return err
			}

			totalCreatedUserPresence++
		}
	}

	s.log.Info(fmt.Sprintf("User presence created: %d", totalCreatedUserPresence))
	return nil
}

func (s *Scheduler) checkForgottenUserPresence(tx *gorm.DB) error {
	s.log.Info("Checking user presence...")

	var userPresences []entity.UserPresence
	if err := tx.Where("status = ? AND start_time IS NOT NULL", enum.PresenceStatusPresent).Find(&userPresences).Error; err != nil {
		return err
	}

	endTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 17, 0, 0, 0, time.UTC)

	for _, userPresence := range userPresences {
		userPresence.Status = enum.PresenceStatusPresent
		userPresence.EndTime = datatype.TimeOnly{Time: &endTime}
		if err := tx.Save(&userPresence).Error; err != nil {
			return err
		}
	}

	s.log.Info(fmt.Sprintf("User presence checked: %d", len(userPresences)))
	return nil
}

// Todo : update every day the chicken cage is need feed

// Todo : create kpi performance every 6 pm

func (s *Scheduler) createUserSalaryPaymentPerMonth(tx *gorm.DB) error {
	s.log.Info("Create user salary payments...")

	var users []entity.User
	if err := tx.Preload("Role").Find(&users).Error; err != nil {
		return err
	}

	data := make([]entity.UserSalaryPayment, 0)
	for _, user := range users {
		if user.Role.Name != constant.RoleOwner {
			data = append(data, entity.UserSalaryPayment{
				UserId:               user.Id,
				BaseSalary:           user.Salary,
				BonusSalary:          decimal.Zero,
				CompentationSalary:   decimal.Zero,
				AdditionalWorkSalary: decimal.Zero,
				Cashbond:             decimal.Zero,
			})
		}
	}

	if err := tx.Model(&entity.UserSalaryPayment{}).CreateInBatches(&data, len(data)).Error; err != nil {
		return err
	}

	s.log.Info(fmt.Sprintf("Success create %d user salary payment", len(data)))
	return nil
}

// Todo : create cashflow history every month

// Todo : notification when h-3 deadline payment date

func (s *Scheduler) Start() {
	s.cron.Start()
}

func (s *Scheduler) Stop() {
	s.cron.Stop()
}
