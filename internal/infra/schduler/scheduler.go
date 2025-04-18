package scheduler

import (
	"github.com/robfig/cron/v3"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"gorm.io/gorm"
)

type Scheduler struct {
	db   *gorm.DB
	cron *cron.Cron
}

type IScheduler interface {
	InitDailyWorkStaff()
}

func NewScheduler(db *gorm.DB) IScheduler {
	return &Scheduler{
		db:   db,
		cron: cron.New(),
	}
}

func (s *Scheduler) InitDailyWorkStaff() {
	cron.New(cron.WithSeconds())

	var dailyWorks []entity.DailyWork
	s.db.Find(&dailyWorks)

	var staffs []entity.Staff
	s.db.Preload("Account").Find(&staffs)

	var dailyWorkStaffs []entity.DailyWorkStaff
	for _, dailyWork := range dailyWorks {
		for _, staff := range staffs {
			if staff.Account.Role.Id == dailyWork.RoleId {
				dailyWorkStaff := entity.DailyWorkStaff{
					DailyWorkId: dailyWork.Id,
					StaffId:     staff.AccountId,
					IsDone:      false,
				}
				dailyWorkStaffs = append(dailyWorkStaffs, dailyWorkStaff)
			}
		}
	}

	s.db.CreateInBatches(dailyWorkStaffs, len(dailyWorkStaffs))
}
