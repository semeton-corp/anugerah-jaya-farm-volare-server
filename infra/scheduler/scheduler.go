package scheduler

import (
	"database/sql"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/robfig/cron/v3"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/constant"
	datatype "github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/custom/data_type"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/util"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"gorm.io/gorm"
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
			err := s.createUserSalaryPaymentPerDaily(tx)
			if err != nil {
				s.log.Error("failed to create user salary payment per daily", zap.Error(err))
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

	s.cron.AddFunc("0 0 1 * *", func() {
		s.db.Transaction(func(tx *gorm.DB) error {
			err := s.refreshChickenCageNeedFeed(tx)
			if err != nil {
				s.log.Error("failed to refresh chicken cage need feed", zap.Error(err))
				return err
			}
			return nil
		})
	})

	s.cron.AddFunc("0 18 * * *", func() {
		s.db.Transaction(func(tx *gorm.DB) error {
			err := s.createKpiChickenCage(tx)
			if err != nil {
				s.log.Error("failed to create kpi chicken performance", zap.Error(err))
				return err
			}
			return nil
		})
	})

	s.cron.AddFunc("0 18 * * *", func() {
		s.db.Transaction(func(tx *gorm.DB) error {
			err := s.createNotificationWhenKPIPerformanceUserBad(tx)
			if err != nil {
				s.log.Error("failed to create notification when kpi performance usere bad", zap.Error(err))
				return err
			}
			return nil
		})
	})

	s.cron.AddFunc("01 00 * * *", func() {
		s.db.Transaction(func(tx *gorm.DB) error {
			err := s.createNotificationTotalItemSaleShipToday(tx)
			if err != nil {
				s.log.Error("failed to create notification total item sale item need ship today", zap.Error(err))
				return err
			}
			return nil
		})
	})

	s.cron.AddFunc("01 00 * * *", func() {
		s.db.Transaction(func(tx *gorm.DB) error {
			err := s.checkChickenCageIfNeedVaccineRoutine(tx)
			if err != nil {
				s.log.Error("failed to check chicken cage if need vaccine routine", zap.Error(err))
				return err
			}
			return nil
		})
	})

	s.cron.AddFunc("01 00 * * *", func() {
		s.db.Transaction(func(tx *gorm.DB) error {
			err := s.createNotificationWarehouseItemInDangerStatus(tx)
			if err != nil {
				s.log.Error("failed to create notification warehouse item in danger status", zap.Error(err))
				return err
			}
			return nil
		})
	})

	s.cron.AddFunc("01 00 * * *", func() {
		s.db.Transaction(func(tx *gorm.DB) error {
			err := s.createNotificationStoreItemGoodEggInDanger(tx)
			if err != nil {
				s.log.Error("failed to create notification store item good egg in danger", zap.Error(err))
				return err
			}
			return nil
		})
	})

	s.cron.AddFunc("01 00 * * *", func() {
		s.db.Transaction(func(tx *gorm.DB) error {
			err := s.createNotificationWhen3DaysBeforeDeadlinePaymentDate(tx)
			if err != nil {
				s.log.Error("failed to create notification when 3 days before deadline payment date", zap.Error(err))
				return err
			}
			return nil
		})
	})

	s.cron.AddFunc("01 00 * * *", func() {
		s.db.Transaction(func(tx *gorm.DB) error {
			err := s.createNotificationAfkirChickenSaleWillTaken(tx)
			if err != nil {
				s.log.Error("failed to create notification item arrive", zap.Error(err))
				return err
			}
			return nil
		})
	})

	s.cron.AddFunc("01 00 * * *", func() {
		s.db.Transaction(func(tx *gorm.DB) error {
			err := s.createNotificationItemArrive(tx)
			if err != nil {
				s.log.Error("failed to create notification item arrive", zap.Error(err))
				return err
			}
			return nil
		})
	})

	s.cron.AddFunc("01 00 * * *", func() {
		s.db.Transaction(func(tx *gorm.DB) error {
			err := s.refreshChickenCageNeedFeed(tx)
			if err != nil {
				s.log.Error("failed to refresh chicken cage need feed", zap.Error(err))
				return err
			}
			return nil
		})
	})

	s.cron.AddFunc("01 00 1 * *", func() {
		s.db.Transaction(func(tx *gorm.DB) error {
			err := s.createCashflowHistoryMonthly(tx)
			if err != nil {
				s.log.Error("failed to create cashflow history monthly", zap.Error(err))
				return err
			}
			return nil
		})
	})
}

func (s *Scheduler) createDailyWorkUser(tx *gorm.DB) error {
	s.log.Info("creating daily work user")

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

	s.log.Info(fmt.Sprintf("daily work user created: %d", len(dailyWorkUsers)))
	return tx.CreateInBatches(dailyWorkUsers, len(dailyWorkUsers)).Error
}

func (s *Scheduler) createUserPresence(tx *gorm.DB) error {
	s.log.Info("creating user presence")

	var users []entity.User
	tx.Model(&entity.User{}).Preload("Role").Find(&users)

	totalCreatedUserPresence := 0
	for _, user := range users {
		if user.Role.Name != constant.RoleOwner {
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
	s.log.Info("checking user presence")

	var userPresences []entity.UserPresence
	if err := tx.Where("status = ? AND start_time IS NOT NULL", enum.PresenceStatusPresent).Find(&userPresences).Error; err != nil {
		return err
	}

	endTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 17, 0, 0, 0, time.Local)

	for _, userPresence := range userPresences {
		userPresence.Status = enum.PresenceStatusPresent
		userPresence.EndTime = datatype.TimeOnly{Time: &endTime}
		if err := tx.Save(&userPresence).Error; err != nil {
			return err
		}
	}

	s.log.Info(fmt.Sprintf("user presence checked: %d", len(userPresences)))
	return nil
}

func (s *Scheduler) refreshChickenCageNeedFeed(tx *gorm.DB) error {
	s.log.Info("refresh chicken cage need feed every day")

	var chickenCages []entity.ChickenCage

	subQuery := tx.Model(&entity.ChickenCage{}).
		Select("MAX(id)").
		Group("cage_id")

	query := tx.Where("chicken_cages.id IN (?)", subQuery)

	err := query.
		Order("chicken_cages.created_at DESC").
		Find(&chickenCages).Error
	if err != nil {
		return err
	}

	chickenCageIds := make([]uint64, 0)
	for _, chickenCage := range chickenCages {
		chickenCageIds = append(chickenCageIds, chickenCage.Id)
	}

	err = tx.Model(&entity.ChickenCage{}).Where("id IN ?", chickenCageIds).Updates(map[string]any{
		"is_need_feed": true,
	}).Error
	if err != nil {
		return err
	}

	s.log.Info(fmt.Sprintf("Success update %d chicken cage to get feed", len(chickenCages)))
	return nil
}

// Todo : check this!!!
func (s *Scheduler) createKpiChickenCage(tx *gorm.DB) error {
	s.log.Info("create kpi performance chicken cage")

	var chickenCages []entity.ChickenCage
	query := tx.Model(&entity.ChickenCage{})
	subQuery := tx.Model(&entity.ChickenCage{}).
		Select("MAX(id)").
		Group("cage_id")
	query = query.Where("chicken_cages.id IN (?)", subQuery)
	err := query.
		Preload("Cage.Location").
		Preload("ChickenProcurement").
		Preload("Cage.CagePlacement.User.Role").
		Order("chicken_cages.created_at DESC").
		Find(&chickenCages).Error

	if err != nil {
		return err
	}

	data := make([]entity.ChickenPerformance, 0)
	notifications := make([]entity.Notification, 0)

	today := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)
	chickenMonitorinMap := make(map[uint64]entity.ChickenMonitoring)
	eggMonitoringMap := make(map[uint64]entity.EggMonitoring)

	var eggMonitorings []entity.EggMonitoring
	err = tx.Model(&entity.EggMonitoring{}).Where("DATE(created_at) = ?", today).Find(&eggMonitorings).Error
	if err != nil {
		return err
	}
	for _, eggMonitoring := range eggMonitorings {
		eggMonitoringMap[eggMonitoring.ChickenCageId] = eggMonitoring
	}

	var chickenMonitorings []entity.ChickenMonitoring
	err = tx.Model(&entity.ChickenMonitoring{}).Where("DATE(created_at) = ?", today).Find(&chickenMonitorings).Error
	if err != nil {
		return err
	}
	for _, chickenMonitoring := range chickenMonitorings {
		chickenMonitorinMap[chickenMonitoring.ChickenCageId] = chickenMonitoring
	}

	for _, chickenCage := range chickenCages {
		avgConsumption := 0.0
		if chickenCage.TotalChicken > 0 {
			avgConsumption = chickenMonitorinMap[chickenCage.Id].TotalFeed / float64(chickenCage.TotalChicken)
		}

		totalEggCount := eggMonitoringMap[chickenCage.Id].TotalGoodEgg + eggMonitoringMap[chickenCage.Id].TotalCrackedEgg
		avgWeight := 0.0
		if totalEggCount > 0 {
			avgWeight = (eggMonitoringMap[chickenCage.Id].TotalWeightGoodEgg + eggMonitoringMap[chickenCage.Id].TotalWeightCrackedEgg) / float64(totalEggCount)
		}

		mortality := 0.0
		if chickenCage.TotalChicken > 0 {
			mortality = float64(chickenMonitorinMap[chickenCage.Id].TotalDeathChicken) / float64(chickenCage.TotalChicken)
		}

		fcr := 0.0
		if chickenCage.TotalChicken > 0 {
			fcr = float64(totalEggCount) / float64(chickenCage.TotalChicken) * 100.0
		}

		hdp := 0.0
		if totalEggCount > 0 {
			hdp = float64(chickenMonitorinMap[chickenCage.Id].TotalFeed) / float64(totalEggCount) * 100.0
		}

		var goodEgg entity.Item
		err = tx.Model(&entity.Item{}).Where("name = ? AND unit = ? AND category = ?", constant.GoodEgg, constant.UnitKg, enum.ItemCategoryEgg).First(&goodEgg).Error
		if err != nil {
			return err
		}

		var goodEggItemPrice entity.ItemPrice
		err = tx.Model(&entity.ItemPrice{}).Where("item_id = ? AND sale_unit = ?", goodEgg.Id, enum.SaleUnitKg).First(&goodEggItemPrice).Error
		if err != nil {
			return err
		}

		getTotalExpenseProductionInMonth := func(db *gorm.DB, month enum.Month, year uint64) (decimal.Decimal, error) {
			totalExpenseProduction := decimal.Zero
			startDate, endDate := util.GetStartDateAndEndDateInMonth(int(year), time.Month(month))

			var warehouseItemProcurements []entity.WarehouseItemProcurement
			if err := db.Where("DATE(deadline_payment_date) BETWEEN ? AND ?", startDate, endDate).
				Find(&warehouseItemProcurements).Error; err != nil {
				return decimal.Zero, err
			}
			for _, e := range warehouseItemProcurements {
				totalExpenseProduction = totalExpenseProduction.Add(e.TotalPrice)
			}

			var warehouseItemCornProcurements []entity.WarehouseItemCornProcurement
			if err := db.Where("DATE(deadline_payment_date) BETWEEN ? AND ?", startDate, endDate).
				Find(&warehouseItemCornProcurements).Error; err != nil {
				return decimal.Zero, err
			}
			for _, e := range warehouseItemCornProcurements {
				totalExpenseProduction = totalExpenseProduction.Add(e.TotalPrice)
			}

			var chickenProcurements []entity.ChickenProcurement
			if err := db.Where("DATE(deadline_payment_date) BETWEEN ? AND ?", startDate, endDate).
				Find(&chickenProcurements).Error; err != nil {
				return decimal.Zero, err
			}
			for _, e := range chickenProcurements {
				totalExpenseProduction = totalExpenseProduction.Add(e.TotalPrice)
			}

			var expenses []entity.Expense
			if err := db.Where("DATE(created_at) BETWEEN ? AND ?", startDate, endDate).
				Find(&expenses).Error; err != nil {
				return decimal.Zero, err
			}
			for _, e := range expenses {
				totalExpenseProduction = totalExpenseProduction.Add(e.Nominal)
			}

			var userSalaryPayments []entity.UserSalaryPayment
			if err := db.Where("DATE(created_at) BETWEEN ? AND ?", startDate, endDate).
				Find(&userSalaryPayments).Error; err != nil {
				return decimal.Zero, err
			}
			for _, e := range userSalaryPayments {
				totalExpenseProduction = totalExpenseProduction.Add(
					e.BaseSalary.Add(e.AdditionalWorkSalary).
						Add(e.BonusSalary).
						Add(e.CompentationSalary).
						Sub(e.Cashbond),
				)
			}

			return totalExpenseProduction, nil
		}

		totalExpenseProduction, err := getTotalExpenseProductionInMonth(tx, enum.Month(time.Now().Month()), uint64(time.Now().Year()))
		if err != nil {
			return err
		}
		totalDayInMonth := util.TotalDaysInMonth(today.Year(), today.Month())
		totalExpensePerDay := totalExpenseProduction.Div(decimal.NewFromUint64(totalDayInMonth))

		chickenAge := util.GetChickenAgeByChickenCage(&chickenCage)
		chickenCategory := util.GetChickenCategoryByChickenCage(&chickenCage)

		newData := entity.ChickenPerformance{
			CageName:                     chickenCage.Cage.Name,
			ChickenCategory:              chickenCategory,
			ChickenAge:                   chickenAge,
			TotalChicken:                 chickenCage.TotalChicken,
			TotalEgg:                     totalEggCount,
			AverageConsumptionPerChicken: avgConsumption,
			AverageWeightPerEgg:          avgWeight,
			MortalityRate:                mortality,
			FCR:                          fcr,
			HDP:                          hdp,
		}

		if chickenAge >= 90 {
			newData.Productivity = enum.ChickenProductivityAfkir
		} else {
			var (
				totalPrice = decimal.Zero
			)

			if eggMonitoringMap[chickenCage.Id].TotalWeightGoodEgg != 0.0 {
				totalPrice = goodEggItemPrice.Price.Mul(decimal.NewFromFloat(eggMonitoringMap[chickenCage.Id].TotalWeightGoodEgg))
			}

			if totalPrice.Sub(totalExpensePerDay).GreaterThanOrEqual(decimal.NewFromInt(constant.MinProfitForCageNotAfkir)) {
				newData.Productivity = enum.ChickenProductivityProductive
			} else {
				newData.Productivity = enum.ChickenProductivityAfkir
			}
		}

		kpiChicken := (mortality + hdp) / 2
		if kpiChicken < constant.ThresholdKpiChicken {
			for _, cagePlacement := range chickenCage.Cage.CagePlacement {
				if cagePlacement.User.Role.Name == constant.RolePekerjaKandang {
					notifications = append(notifications, entity.Notification{
						CageId:               sql.NullInt64{Int64: int64(chickenCage.CageId), Valid: true},
						UserId:               uuid.NullUUID{UUID: cagePlacement.UserId, Valid: true},
						NotificationContexts: []string{constant.ChickenKPINotificationContext},
						Description:          fmt.Sprintf(constant.KPIPerformanceChickenBadNotification, chickenCage.Cage.Name),
					})
				}
			}
		}

		data = append(data, newData)
	}

	if len(data) > 0 {
		err = tx.Model(&entity.ChickenPerformance{}).CreateInBatches(&data, len(data)).Error
		if err != nil {
			return err
		}
	}

	if len(notifications) > 0 {
		err = tx.Model(&entity.Notification{}).CreateInBatches(&notifications, len(notifications)).Error
		if err != nil {
			return err
		}
	}

	s.log.Info(fmt.Sprintf("success create %d kpi chicken cage", len(data)))
	return nil
}

func (s *Scheduler) createUserSalaryPaymentPerMonth(tx *gorm.DB) error {
	s.log.Info("create user salary payments per month")

	var users []entity.User
	if err := tx.Preload("Role").Find(&users).Error; err != nil {
		return err
	}

	data := make([]entity.UserSalaryPayment, 0)
	for _, user := range users {
		if user.Role.Name != constant.RoleOwner && user.SalaryInterval == enum.SalaryIntervalMonthly {
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

	s.log.Info(fmt.Sprintf("success create %d user salary payment", len(data)))
	return nil
}

func (s *Scheduler) createUserSalaryPaymentPerDaily(tx *gorm.DB) error {
	s.log.Info("create user salary payments per daily")

	var users []entity.User
	if err := tx.Preload("Role").Find(&users).Error; err != nil {
		return err
	}

	data := make([]entity.UserSalaryPayment, 0)
	for _, user := range users {
		if user.Role.Name != constant.RoleOwner && user.SalaryInterval == enum.SalaryIntervalDaily {
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

	s.log.Info(fmt.Sprintf("success create %d user salary payment", len(data)))
	return nil
}

func (s *Scheduler) createCashflowHistoryMonthly(tx *gorm.DB) error {
	s.log.Info("create cashflow history monthly")

	now := time.Now()
	year, month, _ := now.Date()

	lastMonth := month - 1
	lastYear := year
	if lastMonth == 0 {
		lastMonth = 12
		lastYear = year - 1
	}

	startDate, endDate := util.GetStartDateAndEndDateInMonth(lastYear, lastMonth)

	var locations []entity.Location
	if err := tx.Find(&locations).Error; err != nil {
		return err
	}

	data := make([]entity.CashflowHistory, 0)

	for _, loc := range locations {
		totalIncome := decimal.Zero
		totalExpense := decimal.Zero
		totalReceivables := decimal.Zero
		totalDebt := decimal.Zero
		totalStoreEggSale := decimal.Zero
		totalWarehouseEggSale := decimal.Zero

		var warehouseSalePayments []entity.WarehouseSalePayment
		if err := tx.Joins("JOIN warehouse_sales ws ON ws.id = warehouse_sale_payments.warehouse_sale_id").
			Joins("JOIN warehouses w ON w.id = ws.warehouse_id").
			Where("w.location_id = ? AND warehouse_sale_payments.payment_date BETWEEN ? AND ?", loc.Id, startDate, endDate).
			Find(&warehouseSalePayments).Error; err != nil {
			return err
		}
		for _, e := range warehouseSalePayments {
			totalIncome = totalIncome.Add(e.Nominal)
		}

		var storeSalePayments []entity.StoreSalePayment
		if err := tx.Joins("JOIN store_sales ss ON ss.id = store_sale_payments.store_sale_id").
			Joins("JOIN stores st ON st.id = ss.store_id").
			Where("st.location_id = ? AND store_sale_payments.payment_date BETWEEN ? AND ?", loc.Id, startDate, endDate).
			Find(&storeSalePayments).Error; err != nil {
			return err
		}
		for _, e := range storeSalePayments {
			totalIncome = totalIncome.Add(e.Nominal)
		}

		var afkirPayments []entity.AfkirChickenSalePayment
		if err := tx.Joins("JOIN afkir_chicken_sales acs ON acs.id = afkir_chicken_sale_payments.afkir_chicken_sale_id").
			Joins("JOIN chicken_cages cc ON cc.id = acs.chicken_cage_id").
			Joins("JOIN cages c ON c.id = cc.cage_id").
			Where("c.location_id = ? AND afkir_chicken_sale_payments.payment_date BETWEEN ? AND ?", loc.Id, startDate, endDate).
			Find(&afkirPayments).Error; err != nil {
			return err
		}
		for _, e := range afkirPayments {
			totalIncome = totalIncome.Add(e.Nominal)
		}

		var userCashAdvancePayment []entity.UserCashAdvancePayment
		if err := tx.Joins("JOIN user_cash_advances uca ON uca.id = user_cash_advance_payments.user_cash_advance_id").
			Joins("JOIN users u ON u.id = uca.user_id").
			Where("u.location_id = ? AND user_cash_advance_payments.payment_date BETWEEN ? AND ?", loc.Id, startDate, endDate).
			Find(&userCashAdvancePayment).Error; err != nil {
			return err
		}
		for _, e := range userCashAdvancePayment {
			totalIncome = totalIncome.Add(e.Nominal)
		}

		var chickenProcurementPayments []entity.ChickenProcurementPayment
		if err := tx.Joins("JOIN chicken_procurements cp ON cp.id = chicken_procurement_payments.chicken_procurement_id").
			Joins("JOIN cages c ON c.id = cp.cage_id").
			Where("c.location_id = ? AND chicken_procurement_payments.payment_date BETWEEN ? AND ?", loc.Id, startDate, endDate).
			Find(&chickenProcurementPayments).Error; err != nil {
			return err
		}
		for _, e := range chickenProcurementPayments {
			totalExpense = totalExpense.Add(e.Nominal)
		}

		var salaries []entity.UserSalaryPayment
		if err := tx.Joins("JOIN users u ON u.id = user_salary_payments.user_id").
			Where("u.location_id = ? AND user_salary_payments.payment_date BETWEEN ? AND ? AND user_salary_payments.is_paid = ?", loc.Id, startDate, endDate, true).
			Find(&salaries).Error; err != nil {
			return err
		}
		for _, e := range salaries {
			totalExpense = totalExpense.
				Add(e.BaseSalary).
				Add(e.AdditionalWorkSalary).
				Add(e.BonusSalary).
				Add(e.CompentationSalary).
				Sub(e.Cashbond)
		}

		var warehouseItemProcurementPayments []entity.WarehouseItemProcurementPayment
		if err := tx.
			Joins("LEFT JOIN warehouse_item_procurements ON warehouse_item_procurements.id = warehouse_item_procurement_payments.warehouse_item_procurement_id").Joins("LEFT JOIN warehouses ON warehouses.id = warehouse_item_procurements.warehouse_id").
			Where("warehouses.location_id = ? AND warehouse_item_procurement_payments.payment_date BETWEEN ? AND ?", loc.Id, startDate, endDate).
			Find(&warehouseItemProcurementPayments).Error; err != nil {
			return err
		}
		for _, e := range warehouseItemProcurementPayments {
			totalExpense = totalExpense.Add(e.Nominal)
		}

		var warehouseItemCornProcurementPayments []entity.WarehouseItemCornProcurementPayment
		if err := tx.Joins("LEFT JOIN warehouse_item_corn_procurements ON warehouse_item_corn_procurements.id = warehouse_item_corn_procurement_payments.warehouse_item_corn_procurement_id").Joins("LEFT JOIN warehouses ON warehouses.id = warehouse_item_corn_procurements.warehouse_id").
			Where("warehouses.location_id = ? AND warehouse_item_corn_procurement_payments.payment_date BETWEEN ? AND ?", loc.Id, startDate, endDate).
			Find(&warehouseItemCornProcurementPayments).Error; err != nil {
			return err
		}
		for _, e := range warehouseItemCornProcurementPayments {
			totalExpense = totalExpense.Add(e.Nominal)
		}

		var expenses []entity.Expense
		if err := tx.Where("location_id = ? AND created_at BETWEEN ? AND ?", loc.Id, startDate, endDate).Find(&expenses).Error; err != nil {
			return err
		}
		for _, e := range expenses {
			totalExpense = totalExpense.Add(e.Nominal)
		}

		var storeSales []entity.StoreSale
		if err := tx.Preload("Payments").
			Joins("JOIN stores st ON st.id = store_sales.store_id").
			Where("st.location_id = ? AND store_sales.created_at BETWEEN ? AND ?", loc.Id, startDate, endDate).
			Find(&storeSales).Error; err != nil {
			return err
		}
		for _, sale := range storeSales {
			totalStoreEggSale = totalStoreEggSale.Add(sale.TotalPrice)
		}

		var warehouseSales []entity.WarehouseSale
		if err := tx.Preload("Payments").
			Joins("JOIN warehouses w ON w.id = warehouse_sales.warehouse_id").
			Where("w.location_id = ? AND warehouse_sales.created_at BETWEEN ? AND ?", loc.Id, startDate, endDate).
			Find(&warehouseSales).Error; err != nil {
			return err
		}
		for _, sale := range warehouseSales {
			totalWarehouseEggSale = totalWarehouseEggSale.Add(sale.TotalPrice)
		}

		var warehouseSaleReceivables []entity.WarehouseSale
		if err := tx.Preload("Payments").
			Joins("JOIN warehouses w ON w.id = warehouse_sales.warehouse_id").
			Where("w.location_id = ? AND warehouse_sales.created_at BETWEEN ? AND ?", loc.Id, startDate, endDate).
			Where("warehouse_sales.payment_status IN ?", []enum.PaymentStatus{
				enum.PaymentStatusNotPaid,
				enum.PaymentStatusUnpaid,
			}).Find(&warehouseSaleReceivables).Error; err != nil {
			return err
		}
		for _, sale := range warehouseSaleReceivables {
			total := sale.TotalPrice
			for _, p := range sale.Payments {
				total = total.Sub(p.Nominal)
			}
			totalReceivables = totalReceivables.Add(total)
		}

		var afkirSales []entity.AfkirChickenSale
		if err := tx.Preload("Payments").
			Joins("JOIN chicken_cages cc ON cc.id = afkir_chicken_sales.chicken_cage_id").
			Joins("JOIN cages c ON c.id = cc.cage_id").
			Where("c.location_id = ? AND afkir_chicken_sales.created_at BETWEEN ? AND ?", loc.Id, startDate, endDate).
			Where("afkir_chicken_sales.payment_status IN ?", []enum.PaymentStatus{
				enum.PaymentStatusNotPaid,
				enum.PaymentStatusUnpaid,
			}).Find(&afkirSales).Error; err != nil {
			return err
		}
		for _, sale := range afkirSales {
			total := sale.TotalPrice
			for _, p := range sale.Payments {
				total = total.Sub(p.Nominal)
			}
			totalReceivables = totalReceivables.Add(total)
		}

		var userAdvances []entity.UserCashAdvance
		if err := tx.Preload("Payments").
			Joins("JOIN users u ON u.id = user_cash_advances.user_id").
			Where("u.location_id = ? AND user_cash_advances.created_at BETWEEN ? AND ?", loc.Id, startDate, endDate).
			Where("user_cash_advances.payment_status IN ?", []enum.PaymentStatus{
				enum.PaymentStatusNotPaid,
				enum.PaymentStatusUnpaid,
			}).Find(&userAdvances).Error; err != nil {
			return err
		}
		for _, adv := range userAdvances {
			total := adv.Nominal
			for _, p := range adv.Payments {
				total = total.Sub(p.Nominal)
			}
			totalReceivables = totalReceivables.Add(total)
		}

		var warehouseItemProcurements []entity.WarehouseItemProcurement
		if err := tx.Preload("Payments").
			Joins("JOIN warehouses w ON w.id = warehouse_item_procurements.warehouse_id").
			Where("w.location_id = ? AND warehouse_item_procurements.created_at BETWEEN ? AND ?", loc.Id, startDate, endDate).
			Where("warehouse_item_procurements.payment_status IN ?", []enum.PaymentStatus{
				enum.PaymentStatusNotPaid,
				enum.PaymentStatusUnpaid,
			}).Find(&warehouseItemProcurements).Error; err != nil {
			return err
		}
		for _, procurement := range warehouseItemProcurements {
			total := procurement.TotalPrice
			for _, p := range procurement.Payments {
				total = total.Sub(p.Nominal)
			}
			totalDebt = totalDebt.Add(total)
		}

		var warehouseItemCornProcurements []entity.WarehouseItemCornProcurement
		if err := tx.Preload("Payments").
			Joins("JOIN warehouses w ON w.id = warehouse_item_corn_procurements.warehouse_id").
			Where("w.location_id = ? AND warehouse_item_corn_procurements.created_at BETWEEN ? AND ?", loc.Id, startDate, endDate).
			Where("warehouse_item_corn_procurements.payment_status IN ?", []enum.PaymentStatus{
				enum.PaymentStatusNotPaid,
				enum.PaymentStatusUnpaid,
			}).Find(&warehouseItemCornProcurements).Error; err != nil {
			return err
		}
		for _, procurement := range warehouseItemCornProcurements {
			total := procurement.TotalPrice
			for _, p := range procurement.Payments {
				total = total.Sub(p.Nominal)
			}
			totalDebt = totalDebt.Add(total)
		}

		var chickenProcurements []entity.ChickenProcurement
		if err := tx.Preload("Payments").
			Joins("JOIN cages c ON c.id = chicken_procurements.cage_id").
			Where("c.location_id = ? AND chicken_procurements.created_at BETWEEN ? AND ?", loc.Id, startDate, endDate).
			Where("chicken_procurements.payment_status IN ?", []enum.PaymentStatus{
				enum.PaymentStatusNotPaid,
				enum.PaymentStatusUnpaid,
			}).Find(&chickenProcurements).Error; err != nil {
			return err
		}
		for _, procurement := range chickenProcurements {
			total := procurement.TotalPrice
			for _, p := range procurement.Payments {
				total = total.Sub(p.Nominal)
			}
			totalDebt = totalDebt.Add(total)
		}

		history := entity.CashflowHistory{
			LocationId:       loc.Id,
			Income:           totalIncome,
			Expense:          totalExpense,
			Receivables:      totalReceivables,
			Debt:             totalDebt,
			Cash:             totalIncome.Add(totalReceivables),
			Profit:           totalIncome.Sub(totalExpense),
			WarehouseEggSale: totalWarehouseEggSale,
			StoreEggSale:     totalStoreEggSale,
			CreatedAt:        endDate,
		}

		data = append(data, history)
	}

	if len(data) > 0 {
		if err := tx.CreateInBatches(&data, len(data)).Error; err != nil {
			return err
		}
	}

	s.log.Info(fmt.Sprintf("success create %d cashflow histories", len(data)))
	return nil
}

func (s *Scheduler) createNotificationWhen3DaysBeforeDeadlinePaymentDate(tx *gorm.DB) error {
	s.log.Info("create notification for 3 days before deadline payment date")

	now := time.Now()
	targetDate := now.AddDate(0, 0, 3).Format("2006-01-02")

	notifications := make([]entity.Notification, 0)
	var warehouseSales []entity.WarehouseSale
	if err := tx.Preload("Warehouse").Preload("Customer").Where("DATE(deadline_payment_date) = ?", targetDate).
		Find(&warehouseSales).Error; err != nil {
		return err
	}
	for _, ws := range warehouseSales {
		notifications = append(notifications, entity.Notification{
			WarehouseId:          sql.NullInt64{Int64: int64(ws.WarehouseId), Valid: true},
			NotificationContexts: pq.StringArray{constant.WarehouseSaleNotificationContext, constant.ReceivablesNotificationContext},
			Description:          fmt.Sprintf(constant.PaymentReceivablesDeadlineNotification, ws.Customer.Name),
		})
	}

	var storeSales []entity.StoreSale
	if err := tx.Preload("Store").Preload("Customer").Where("DATE(deadline_payment_date) = ?", targetDate).
		Find(&storeSales).Error; err != nil {
		return err
	}
	for _, ss := range storeSales {
		notifications = append(notifications, entity.Notification{
			StoreId:              sql.NullInt64{Int64: int64(ss.StoreId), Valid: true},
			NotificationContexts: pq.StringArray{constant.StoreSaleNotificationContext, constant.ReceivablesNotificationContext},
			Description:          fmt.Sprintf(constant.PaymentReceivablesDeadlineNotification, ss.Customer.Name),
		})
	}

	var afkirSales []entity.AfkirChickenSale
	if err := tx.Where("DATE(deadline_payment_date) = ?", targetDate).Preload("AfkirChickenCustomer").Preload("ChickenCage.Cage").
		Find(&afkirSales).Error; err != nil {
		return err
	}
	for _, as := range afkirSales {
		notifications = append(notifications, entity.Notification{
			CageId:               sql.NullInt64{Int64: int64(as.ChickenCage.CageId), Valid: true},
			NotificationContexts: pq.StringArray{constant.AfkirChickenSaleNotificationContext, constant.ReceivablesNotificationContext},
			Description:          fmt.Sprintf(constant.PaymentReceivablesDeadlineNotification, as.AfkirChickenCustomer.Name),
		})
	}

	var chickenProcurements []entity.ChickenProcurement
	if err := tx.Where("DATE(deadline_payment_date) = ?", targetDate).Preload("Supplier").Preload("Cage").
		Find(&chickenProcurements).Error; err != nil {
		return err
	}
	for _, cp := range chickenProcurements {
		notifications = append(notifications, entity.Notification{
			CageId:               sql.NullInt64{Int64: int64(cp.CageId), Valid: true},
			NotificationContexts: pq.StringArray{constant.ChickenProcurementNotificationContext, constant.DebtNotificationContext},
			Description:          fmt.Sprintf(constant.PaymentDebtDeadlineNotification, cp.Supplier.Name),
		})
	}

	var itemProcurements []entity.WarehouseItemProcurement
	if err := tx.Where("DATE(deadline_payment_date) = ?", targetDate).Preload("Warehouse").Preload("Supplier").
		Find(&itemProcurements).Error; err != nil {
		return err
	}
	for _, wp := range itemProcurements {
		notifications = append(notifications, entity.Notification{
			WarehouseId:          sql.NullInt64{Int64: int64(wp.WarehouseId), Valid: true},
			NotificationContexts: pq.StringArray{constant.WarehouseItemProcurementNotificationContext, constant.DebtNotificationContext},
			Description:          fmt.Sprintf(constant.PaymentDebtDeadlineNotification, wp.Supplier.Name),
		})
	}

	var cornProcurements []entity.WarehouseItemCornProcurement
	if err := tx.Where("DATE(deadline_payment_date) = ?", targetDate).
		Find(&cornProcurements).Error; err != nil {
		return err
	}
	for _, cp := range cornProcurements {
		notifications = append(notifications, entity.Notification{
			WarehouseId:          sql.NullInt64{Int64: int64(cp.WarehouseId), Valid: true},
			NotificationContexts: pq.StringArray{constant.WarehouseItemCornProcurementNotificationContext, constant.DebtNotificationContext},
			Description:          fmt.Sprintf(constant.PaymentDebtDeadlineNotification, cp.Supplier.Name),
		})
	}

	var userCashAdvances []entity.UserCashAdvance
	if err := tx.Where("DATE(deadline_payment_date) = ?", targetDate).Preload("User").
		Find(&userCashAdvances).Error; err != nil {
		return err
	}
	for _, cp := range userCashAdvances {
		notifications = append(notifications, entity.Notification{
			UserId:               uuid.NullUUID{UUID: cp.UserId, Valid: true},
			NotificationContexts: pq.StringArray{constant.UserCashAdvanceNotificationContext, constant.ReceivablesNotificationContext},
			Description:          fmt.Sprintf(constant.PaymentDebtDeadlineNotification, cp.User.Name),
		})
	}

	err := tx.Model(&entity.Notification{}).CreateInBatches(&notifications, 10).Error
	if err != nil {
		return err
	}

	s.log.Info(fmt.Sprintf("success create %d notification for 3 days before deadline payment date", len(notifications)))

	return nil
}

func (s *Scheduler) checkChickenCageIfNeedVaccineRoutine(tx *gorm.DB) error {
	s.log.Info("check chicken cage if need routine vaccine")

	var chickenCages []*entity.ChickenCage
	query := tx.Model(&entity.ChickenCage{})
	subQuery := tx.Model(&entity.ChickenCage{}).
		Select("MAX(id)").
		Group("cage_id")
	query = query.Where("chicken_cages.id IN (?)", subQuery)

	err := query.
		Preload("Cage.Location").
		Preload("ChickenProcurement").
		Preload("Cage.CagePlacement.User.Role").
		Order("chicken_cages.created_at DESC").
		Find(&chickenCages).Error
	if err != nil {
		return fmt.Errorf("failed to fetch chicken cages: %w", err)
	}

	chickenCageNeedVaccineIds := make([]uint64, 0)
	notifications := make([]*entity.Notification, 0)

	for _, chickenCage := range chickenCages {
		if !chickenCage.ChickenProcurementId.Valid {
			continue
		}

		chickenAge := util.GetChickenAgeByChickenCage(chickenCage)
		var chickenHealthItems []*entity.ChickenHealthItem
		err := tx.Model(&entity.ChickenHealthItem{}).
			Where("chicken_age = ? AND type = ?", int64(chickenAge), enum.ChickenHealthItemTypeVaccineRoutine).
			Find(&chickenHealthItems).Error
		if err != nil {
			return fmt.Errorf("failed to fetch chicken health items for age %d: %w", chickenAge, err)
		}

		if len(chickenHealthItems) == 0 {
			continue
		}

		if chickenAge != uint64(chickenCage.LatestChickenAgeVaccineRoutine.Int64) {
			var vaccineRoutineNames []string
			for _, chickenHealthItem := range chickenHealthItems {
				vaccineRoutineNames = append(vaccineRoutineNames, chickenHealthItem.Name)
			}

			chickenCageNeedVaccineIds = append(chickenCageNeedVaccineIds, chickenCage.Id)
			notifications = append(notifications, &entity.Notification{
				CageId:               sql.NullInt64{Int64: int64(chickenCage.CageId), Valid: true},
				NotificationContexts: pq.StringArray{constant.VaccineMonitoringNotificationContext},
				Description:          fmt.Sprintf(constant.VaccineRoutineNotification, strings.Join(vaccineRoutineNames, ","), chickenCage.Cage.Name),
			})
		}
	}

	if len(chickenCageNeedVaccineIds) > 0 {
		err := tx.Model(&entity.ChickenCage{}).
			Where("id IN ?", chickenCageNeedVaccineIds).
			Updates(map[string]any{
				"is_need_routine_vaccine": true,
			}).Error
		if err != nil {
			return fmt.Errorf("failed to update chicken cages vaccine status: %w", err)
		}
	}

	if len(notifications) > 0 {
		err := tx.Model(&entity.Notification{}).
			CreateInBatches(notifications, len(notifications)).Error
		if err != nil {
			return fmt.Errorf("failed to create vaccine notifications: %w", err)
		}
	}

	s.log.Info(fmt.Sprintf("success check chicken cage if need routine vaccine, total chicken cage need vaccine: %d", len(chickenCageNeedVaccineIds)))
	return nil
}

func (s *Scheduler) createNotificationTotalItemSaleShipToday(tx *gorm.DB) error {
	s.log.Info("create notification for total shipped today")

	today := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)

	var storeResults []struct {
		StoreId   uint64
		StoreName string
		Total     int64
	}

	err := tx.Table("store_sales AS ss").
		Select("ss.store_id, st.name AS store_name, COUNT(*) AS total").
		Joins("JOIN stores st ON ss.store_id = st.id").
		Where("DATE(ss.send_date) = ?", today).
		Group("ss.store_id, st.name").
		Scan(&storeResults).Error
	if err != nil {
		return err
	}

	var warehouseResults []struct {
		WarehouseId   uint64
		WarehouseName string
		Total         int64
	}

	err = tx.Table("warehouse_sales AS ws").
		Select("ws.warehouse_id, wh.name AS warehouse_name, COUNT(*) AS total").
		Joins("JOIN warehouses wh ON ws.warehouse_id = wh.id").
		Where("DATE(ws.send_date) = ?", today).
		Group("ws.warehouse_id, wh.name").
		Scan(&warehouseResults).Error
	if err != nil {
		return err
	}

	var notifications []entity.Notification
	for _, r := range storeResults {
		if r.Total > 0 {
			notifications = append(notifications, entity.Notification{
				StoreId:              sql.NullInt64{Int64: int64(r.StoreId), Valid: true},
				NotificationContexts: pq.StringArray{constant.StoreSaleNotificationContext},
				Description:          fmt.Sprintf(constant.ItemShipTodayStoreSaleNotification, r.Total, r.StoreName),
			})
			s.log.Info(fmt.Sprintf("store %s (id=%d) has %d sales to ship today", r.StoreName, r.StoreId, r.Total))
		}
	}

	for _, r := range warehouseResults {
		if r.Total > 0 {
			notifications = append(notifications, entity.Notification{
				WarehouseId:          sql.NullInt64{Int64: int64(r.WarehouseId), Valid: true},
				NotificationContexts: pq.StringArray{constant.WarehouseSaleNotificationContext},
				Description:          fmt.Sprintf(constant.ItemShipTodayWarehouseSaleNotification, r.Total, r.WarehouseName),
			})
			s.log.Info(fmt.Sprintf("warehouse %s (id=%d) has %d sales to ship today", r.WarehouseName, r.WarehouseId, r.Total))
		}
	}

	if len(notifications) > 0 {
		if err := tx.CreateInBatches(&notifications, 10).Error; err != nil {
			return err
		}
	}

	s.log.Info(fmt.Sprintf("create %d notification about total item ship today", len(notifications)))

	return nil
}

func (s *Scheduler) createNotificationWarehouseItemInDangerStatus(tx *gorm.DB) error {
	s.log.Error("create notification if warehouse item in danger status")

	var warehouseItems []entity.WarehouseItem
	err := tx.Model(&entity.WarehouseItem{}).Preload("Item").Preload("Warehouse.Location").Find(&warehouseItems).Error
	if err != nil {
		return err
	}

	notifications := make([]entity.Notification, 0)
	for _, warehouseItem := range warehouseItems {
		if warehouseItem.Item.Category != enum.ItemCategoryEgg && warehouseItem.Item.Category != enum.ItemCategoryCornMaterial {
			daysLeft := math.Ceil(warehouseItem.Quantity / warehouseItem.Item.DailySpending.Float64)

			if daysLeft < 3 {
				notifications = append(notifications, entity.Notification{
					WarehouseId:          sql.NullInt64{Int64: int64(warehouseItem.WarehouseId), Valid: true},
					NotificationContexts: pq.StringArray{constant.WarehouseItemNotificationContext},
					Description:          fmt.Sprintf(constant.WarehouseItemInDangerNotification, warehouseItem.Item.Name, warehouseItem.Warehouse.Name),
				})
			}
		}
	}

	if len(notifications) > 0 {
		err = tx.Model(&entity.Notification{}).CreateInBatches(&notifications, len(notifications)).Error
		if err != nil {
			return err
		}
	}

	s.log.Info(fmt.Sprintf("success create %d notification if warehouse item in danger status", len(notifications)))
	return nil
}

func (s *Scheduler) createNotificationStoreItemGoodEggInDanger(tx *gorm.DB) error {
	s.log.Info("create notification store item good egg in danger")

	var goodEggItem entity.Item
	err := tx.Model(&entity.Item{}).Where("name = ? AND unit = ? AND category = ?", constant.GoodEgg, constant.UnitKg, enum.ItemCategoryEgg).First(&goodEggItem).Error
	if err != nil {
		return err
	}

	var storeItems []entity.StoreItem
	err = tx.Model(&entity.StoreItem{}).Preload("Item").Preload("Store.Location").Where("item_id = ?", goodEggItem.Id).Find(&storeItems).Error
	if err != nil {
		return err
	}

	notifications := make([]entity.Notification, 0)
	for _, storeItem := range storeItems {
		if storeItem.Quantity/float64(constant.TotalEggPerIkat) < 20.0 {
			notifications = append(notifications, entity.Notification{
				StoreId:              sql.NullInt64{Int64: int64(storeItem.StoreId), Valid: true},
				NotificationContexts: pq.StringArray{constant.StoreItemNotificationContext},
				Description:          fmt.Sprintf(constant.StoreItemInDangerNotification, storeItem.Item.Name, storeItem.Store.Name),
			})
		}
	}

	if len(notifications) > 0 {
		err = tx.Model(&entity.Notification{}).CreateInBatches(&notifications, len(notifications)).Error
		if err != nil {
			return err
		}
	}

	s.log.Info(fmt.Sprintf("success create %d notification store item good egg in danger", len(notifications)))

	return nil
}

func (s *Scheduler) createNotificationItemArrive(tx *gorm.DB) error {
	s.log.Info("create notification item arrived")

	today := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)
	var chickenProcurements []entity.ChickenProcurement
	err := tx.Model(&entity.ChickenProcurement{}).Where("DATE(estimation_arrival_date) = ?", today).Preload("Cage").Find(&chickenProcurements).Error
	if err != nil {
		return err
	}

	var warehouseItemProcurements []entity.WarehouseItemProcurement
	err = tx.Model(&entity.WarehouseItemProcurement{}).Where("DATE(estimation_arrival_date) = ?", today).Preload("Item").Find(&warehouseItemProcurements).Error
	if err != nil {
		return err
	}

	// Todo : Implement for spesific user Id
	notifications := make([]entity.Notification, 0)
	for _, chickenProcurement := range chickenProcurements {
		notifications = append(notifications, entity.Notification{
			CageId:               sql.NullInt64{Int64: int64(chickenProcurement.CageId), Valid: true},
			Description:          fmt.Sprintf(constant.WorkChickenArriveNotification, chickenProcurement.Cage.Name),
			NotificationContexts: pq.StringArray{constant.ChickenProcurementNotificationContext, constant.WorkNotificationContext},
		})
	}

	for _, warehouseItemProcurement := range warehouseItemProcurements {
		notifications = append(notifications, entity.Notification{
			WarehouseId:          sql.NullInt64{Int64: int64(warehouseItemProcurement.WarehouseId), Valid: true},
			Description:          fmt.Sprintf(constant.WorkItemArriveNotification, warehouseItemProcurement.Item.Name),
			NotificationContexts: pq.StringArray{constant.WarehouseItemProcurementNotificationContext, constant.WorkNotificationContext},
		})
	}

	if len(notifications) > 0 {
		err = tx.Model(&entity.Notification{}).CreateInBatches(&notifications, 10).Error
		if err != nil {
			return err
		}
	}

	s.log.Info(fmt.Sprintf("success create %d notification for item arrive", len(notifications)))
	return nil
}

func (s *Scheduler) createNotificationAfkirChickenSaleWillTaken(tx *gorm.DB) error {
	s.log.Info("create notification afkir chicken sale will taken")

	today := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, 1)
	var afkirChickenSales []entity.AfkirChickenSale
	err := tx.Model(&entity.AfkirChickenSale{}).Where("DATE(taken_at) = ?", today).Preload("ChickenCage.Cage").Find(&afkirChickenSales).Error
	if err != nil {
		return err
	}

	// Todo : Implement for spesific user Id
	notifications := make([]entity.Notification, 0)
	for _, afkirChickenSale := range afkirChickenSales {
		notifications = append(notifications, entity.Notification{
			CageId:               sql.NullInt64{Int64: int64(afkirChickenSale.ChickenCage.CageId), Valid: true},
			Description:          fmt.Sprintf(constant.WorkAfkirChickenTakenTommorow, afkirChickenSale.ChickenCage.Cage.Name),
			NotificationContexts: pq.StringArray{constant.AfkirChickenSaleNotificationContext, constant.WorkNotificationContext},
		})
	}

	if len(notifications) > 0 {
		err = tx.Model(&entity.Notification{}).CreateInBatches(&notifications, 10).Error
		if err != nil {
			return err
		}
	}

	s.log.Info(fmt.Sprintf("success create %d notification afkir chicken sale will taken", len(notifications)))
	return nil
}

func (s *Scheduler) createNotificationWhenKPIPerformanceUserBad(tx *gorm.DB) error {
	s.log.Info("create notification when KPI performance user bad")

	now := time.Now()
	month := int(now.Month())
	year := now.Year()

	startDate, endDate := util.GetStartDayAndEndDayByMonthFilter(enum.Month(month), year)

	var users []entity.User
	if err := tx.Model(&entity.User{}).Find(&users).Error; err != nil {
		return err
	}

	notifications := make([]entity.Notification, 0)

	for _, user := range users {
		var dailyWorkUsers []entity.DailyWorkUser
		if err := tx.Model(&entity.DailyWorkUser{}).
			Joins("JOIN daily_works ON daily_work_users.daily_work_id = daily_works.id").
			Where("daily_work_users.user_id = ?", user.Id).
			Where("daily_work_users.created_at BETWEEN ? AND ?", startDate, endDate).
			Where("daily_works.deleted_at IS NULL").
			Find(&dailyWorkUsers).Error; err != nil {
			s.log.Error("failed to query daily work", zap.Error(err))
			continue
		}

		var additionalWorkUsers []entity.AdditionalWorkUser
		if err := tx.Model(&entity.AdditionalWorkUser{}).
			Joins("JOIN additional_works ON additional_work_users.additional_work_id = additional_works.id").
			Where("additional_work_users.user_id = ?", user.Id).
			Where("additional_work_users.created_at BETWEEN ? AND ?", startDate, endDate).
			Where("additional_works.deleted_at IS NULL").
			Find(&additionalWorkUsers).Error; err != nil {
			s.log.Error("failed to query additional work", zap.Error(err))
			continue
		}

		var userPresences []entity.UserPresence
		if err := tx.Model(&entity.UserPresence{}).
			Where("user_id = ?", user.Id).
			Where("created_at BETWEEN ? AND ?", startDate, endDate).
			Find(&userPresences).Error; err != nil {
			s.log.Error("failed to query presence", zap.Error(err))
			continue
		}

		presenceScore, workScore, _ := util.CalculateKPIScoreUserInMonthViaEntity(
			additionalWorkUsers,
			dailyWorkUsers,
			userPresences,
		)

		totalScore := (0.6 * presenceScore) + (0.4 * workScore)
		if totalScore <= constant.KPIScoreBad {
			notifications = append(notifications, entity.Notification{
				UserId:               uuid.NullUUID{UUID: user.Id, Valid: true},
				NotificationContexts: pq.StringArray{constant.UserKPINotificationContext},
				Description:          fmt.Sprintf(constant.KPIPerformanceUserBadNotification, user.Name),
			})

			s.log.Info(fmt.Sprintf("user %s (id=%s) has bad KPI score %.2f", user.Name, user.Id.String(), totalScore))
		}
	}

	if len(notifications) > 0 {
		if err := tx.CreateInBatches(&notifications, 10).Error; err != nil {
			return err
		}
	}

	s.log.Info(fmt.Sprintf("create %d notifications for bad KPI performance", len(notifications)))
	return nil
}

func (s *Scheduler) refreshNeedFeed(tx *gorm.DB) error {
	s.log.Info("refresh need feed...")

	var chickenCages []*entity.ChickenCage
	query := tx.Model(&entity.ChickenCage{})
	subQuery := tx.Model(&entity.ChickenCage{}).
		Select("MAX(id)").
		Group("cage_id")
	query = query.Where("chicken_cages.id IN (?)", subQuery)

	err := query.
		Preload("Cage.Location").
		Preload("ChickenProcurement").
		Preload("Cage.CagePlacement.User.Role").
		Order("chicken_cages.created_at DESC").
		Find(&chickenCages).Error
	if err != nil {
		return fmt.Errorf("failed to fetch chicken cages: %w", err)
	}

	for _, chickenCage := range chickenCages {
		chickenCage.IsNeedFeed = true
	}

	if err := tx.Save(&chickenCages).Error; err != nil {
		return err
	}

	s.log.Info("success refresh need feed chicken cage")
	return nil
}

func (s *Scheduler) Start() {
	s.cron.Start()
}

func (s *Scheduler) Stop() {
	s.cron.Stop()
}
