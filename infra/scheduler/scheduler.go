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

func (s *Scheduler) refreshChickenCageNeedFeed(tx *gorm.DB) error {
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

func (s *Scheduler) createKpiChickenCage(tx *gorm.DB) error {
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

	today := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)
	for _, chickenCage := range chickenCages {
		var chickenMonitoring entity.ChickenMonitoring
		err := tx.
			Where("chicken_cage_id = ? AND DATE(created_at) = ?", chickenCage.Id, today).
			First(&chickenMonitoring).Error
		if err != nil {
			return err
		}

		var eggMonitoring entity.EggMonitoring
		err = tx.
			Where("chicken_cage_id = ? AND DATE(created_at) = ?", chickenCage.Id, today).
			First(&eggMonitoring).Error
		if err != nil {
			return err
		}

		avgConsumption := 0.0
		if chickenCage.TotalChicken > 0 {
			avgConsumption = chickenMonitoring.TotalFeed / float64(chickenCage.TotalChicken)
		}

		totalEggCount := eggMonitoring.TotalGoodEgg + eggMonitoring.TotalCrackedEgg
		avgWeight := 0.0
		if totalEggCount > 0 {
			avgWeight = (eggMonitoring.TotalWeightGoodEgg + eggMonitoring.TotalWeightCrackedEgg) / float64(totalEggCount)
		}

		mortality := 0.0
		if chickenCage.TotalChicken > 0 {
			mortality = float64(chickenMonitoring.TotalDeathChicken) / float64(chickenCage.TotalChicken)
		}

		fcr := 0.0
		if chickenCage.TotalChicken > 0 {
			fcr = float64(totalEggCount) / float64(chickenCage.TotalChicken) * 100.0
		}

		hdp := 0.0
		if totalEggCount > 0 {
			hdp = float64(chickenMonitoring.TotalFeed) / float64(totalEggCount) * 100.0
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

			// chicken procurements
			var chickenProcurements []entity.ChickenProcurement
			if err := db.Where("DATE(deadline_payment_date) BETWEEN ? AND ?", startDate, endDate).
				Find(&chickenProcurements).Error; err != nil {
				return decimal.Zero, err
			}
			for _, e := range chickenProcurements {
				totalExpenseProduction = totalExpenseProduction.Add(e.TotalPrice)
			}

			// expenses
			var expenses []entity.Expense
			if err := db.Where("DATE(created_at) BETWEEN ? AND ?", startDate, endDate).
				Find(&expenses).Error; err != nil {
				return decimal.Zero, err
			}
			for _, e := range expenses {
				totalExpenseProduction = totalExpenseProduction.Add(e.Nominal)
			}

			// user salary payments
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

			if eggMonitoring.TotalWeightGoodEgg != 0.0 {
				totalPrice = goodEggItemPrice.Price.Mul(decimal.NewFromFloat(eggMonitoring.TotalWeightGoodEgg))
			}

			if totalPrice.Sub(totalExpensePerDay).GreaterThanOrEqual(decimal.NewFromInt(constant.MinProfitForCageNotAfkir)) {
				newData.Productivity = enum.ChickenProductivityProductive
			} else {
				newData.Productivity = enum.ChickenProductivityAfkir
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

	s.log.Info(fmt.Sprintf("success create %d kpi chicken cage", len(data)))
	return nil
}

func (s *Scheduler) createUserSalaryPaymentPerMonth(tx *gorm.DB) error {
	s.log.Info("create user salary payments per month...")

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
	s.log.Info("create user salary payments per daily...")

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
	s.log.Info("create cashflow history monthly...")

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
		if err := tx.Where("location_id = ? AND created_at BETWEEN ? AND ?", loc.Id, startDate, endDate).
			Find(&warehouseSalePayments).Error; err != nil {
			return err
		}
		for _, e := range warehouseSalePayments {
			totalIncome = totalIncome.Add(e.Nominal)
		}

		var storeSalePayments []entity.StoreSalePayment
		if err := tx.Where("location_id = ? AND created_at BETWEEN ? AND ?", loc.Id, startDate, endDate).
			Find(&storeSalePayments).Error; err != nil {
			return err
		}
		for _, e := range storeSalePayments {
			totalIncome = totalIncome.Add(e.Nominal)
		}

		var afkirPayments []entity.AfkirChickenSalePayment
		if err := tx.Where("created_at BETWEEN ? AND ?", startDate, endDate).
			Find(&afkirPayments).Error; err != nil {
			return err
		}
		for _, e := range afkirPayments {
			totalIncome = totalIncome.Add(e.Nominal)
		}

		var chickenProcurements []entity.ChickenProcurementPayment
		if err := tx.Where("location_id = ? AND created_at BETWEEN ? AND ?", loc.Id, startDate, endDate).
			Find(&chickenProcurements).Error; err != nil {
			return err
		}
		for _, e := range chickenProcurements {
			totalExpense = totalExpense.Add(e.Nominal)
		}

		var salaries []entity.UserSalaryPayment
		if err := tx.Where("location_id = ? AND created_at BETWEEN ? AND ? AND is_paid = ?", loc.Id, startDate, endDate, true).
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

		var storeSales []entity.StoreSale
		if err := tx.Preload("Payments").
			Where("location_id = ? AND deadline_payment BETWEEN ? AND ?", loc.Id, startDate, endDate).
			Find(&storeSales).Error; err != nil {
			return err
		}
		for _, sale := range storeSales {
			total := sale.TotalPrice
			for _, p := range sale.Payments {
				total = total.Sub(p.Nominal)
			}
			totalReceivables = totalReceivables.Add(total)
			totalStoreEggSale = totalStoreEggSale.Add(sale.TotalPrice)
		}

		var warehouseSales []entity.WarehouseSale
		if err := tx.Preload("Payments").
			Where("location_id = ? AND deadline_payment BETWEEN ? AND ?", loc.Id, startDate, endDate).
			Find(&warehouseSales).Error; err != nil {
			return err
		}
		for _, sale := range warehouseSales {
			total := sale.TotalPrice
			for _, p := range sale.Payments {
				total = total.Sub(p.Nominal)
			}
			totalReceivables = totalReceivables.Add(total)
			totalWarehouseEggSale = totalWarehouseEggSale.Add(sale.TotalPrice)
		}

		var warehouseItemProcurements []entity.WarehouseItemProcurement
		if err := tx.Preload("Payments").
			Where("location_id = ? AND deadline_payment BETWEEN ? AND ?", loc.Id, startDate, endDate).
			Find(&warehouseItemProcurements).Error; err != nil {
			return err
		}
		for _, procurement := range warehouseItemProcurements {
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
			Profit:           totalIncome.Add(totalReceivables).Sub(totalExpense.Add(totalDebt)),
			WarehouseEggSale: totalWarehouseEggSale,
			StoreEggSale:     totalStoreEggSale,
			CreatedAt:        time.Now(),
		}

		data = append(data, history)
	}

	// Bulk insert
	if len(data) > 0 {
		if err := tx.CreateInBatches(&data, len(data)).Error; err != nil {
			return err
		}
	}

	s.log.Info(fmt.Sprintf("success create %d cashflow histories", len(data)))
	return nil
}

func (s *Scheduler) createNotificationWhen3DaysBeforeDeadlinePaymentDate(tx *gorm.DB) error {
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
			WarehouseId:         sql.NullInt64{Int64: int64(ws.WarehouseId), Valid: true},
			NotificationContext: pq.StringArray{constant.WarehouseSaleNotificationContext, constant.ReceivablesNotificationContext},
			Description:         fmt.Sprintf(constant.PaymentReceivablesDeadlineNotification, ws.Customer.Name),
		})
	}

	var storeSales []entity.StoreSale
	if err := tx.Preload("Store").Preload("Customer").Where("DATE(deadline_payment_date) = ?", targetDate).
		Find(&storeSales).Error; err != nil {
		return err
	}
	for _, ss := range storeSales {
		notifications = append(notifications, entity.Notification{
			StoreId:             sql.NullInt64{Int64: int64(ss.StoreId), Valid: true},
			NotificationContext: pq.StringArray{constant.StoreSaleNotificationContext, constant.ReceivablesNotificationContext},
			Description:         fmt.Sprintf(constant.PaymentReceivablesDeadlineNotification, ss.Customer.Name),
		})
	}

	var afkirSales []entity.AfkirChickenSale
	if err := tx.Where("DATE(deadline_payment_date) = ?", targetDate).Preload("AfkirChickenCustomer").Preload("ChickenCage.Cage").
		Find(&afkirSales).Error; err != nil {
		return err
	}
	for _, as := range afkirSales {
		notifications = append(notifications, entity.Notification{
			CageId:              sql.NullInt64{Int64: int64(as.ChickenCage.CageId), Valid: true},
			NotificationContext: pq.StringArray{constant.AfkirChickenSaleNotificationContext, constant.ReceivablesNotificationContext},
			Description:         fmt.Sprintf(constant.PaymentReceivablesDeadlineNotification, as.AfkirChickenCustomer.Name),
		})
	}

	var chickenProcurements []entity.ChickenProcurement
	if err := tx.Where("DATE(deadline_payment_date) = ?", targetDate).Preload("Supplier").Preload("Cage").
		Find(&chickenProcurements).Error; err != nil {
		return err
	}
	for _, cp := range chickenProcurements {
		notifications = append(notifications, entity.Notification{
			CageId:              sql.NullInt64{Int64: int64(cp.CageId), Valid: true},
			NotificationContext: pq.StringArray{constant.ChickenProcurementNotificationContext, constant.DebtNotificationContext},
			Description:         fmt.Sprintf(constant.PaymentDebtDeadlineNotification, cp.Supplier.Name),
		})
	}

	var itemProcurements []entity.WarehouseItemProcurement
	if err := tx.Where("DATE(deadline_payment_date) = ?", targetDate).Preload("Warehouse").Preload("Supplier").
		Find(&itemProcurements).Error; err != nil {
		return err
	}
	for _, wp := range itemProcurements {
		notifications = append(notifications, entity.Notification{
			WarehouseId:         sql.NullInt64{Int64: int64(wp.WarehouseId), Valid: true},
			NotificationContext: pq.StringArray{constant.WarehouseItemProcurementNotificationContext, constant.DebtNotificationContext},
			Description:         fmt.Sprintf(constant.PaymentDebtDeadlineNotification, wp.Supplier.Name),
		})
	}

	var cornProcurements []entity.WarehouseItemCornProcurement
	if err := tx.Where("DATE(deadline_payment_date) = ?", targetDate).
		Find(&cornProcurements).Error; err != nil {
		return err
	}
	for _, cp := range cornProcurements {
		notifications = append(notifications, entity.Notification{
			WarehouseId:         sql.NullInt64{Int64: int64(cp.WarehouseId), Valid: true},
			NotificationContext: pq.StringArray{constant.WarehouseItemCornProcurementNotificationContext, constant.DebtNotificationContext},
			Description:         fmt.Sprintf(constant.PaymentDebtDeadlineNotification, cp.Supplier.Name),
		})
	}

	err := tx.Model(&entity.Notification{}).CreateInBatches(&notifications, 10).Error
	if err != nil {
		return err
	}

	var userCashAdvances []entity.UserCashAdvance
	if err := tx.Where("DATE(deadline_payment_date) = ?", targetDate).Preload("User").
		Find(&userCashAdvances).Error; err != nil {
		return err
	}
	for _, cp := range userCashAdvances {
		notifications = append(notifications, entity.Notification{
			UserId:              uuid.NullUUID{UUID: cp.UserId, Valid: true},
			NotificationContext: pq.StringArray{constant.UserCashAdvanceNotificationContext, constant.ReceivablesNotificationContext},
			Description:         fmt.Sprintf(constant.PaymentDebtDeadlineNotification, cp.User.Name),
		})
	}

	err = tx.Model(&entity.Notification{}).CreateInBatches(&notifications, 10).Error
	if err != nil {
		return err
	}

	return nil
}

func (s *Scheduler) checkChickenCageIfNeedVaccineRoutine(tx *gorm.DB) error {
	s.log.Info("check chicken cage if need routine vaccine")

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

	for _, chickenCage := range chickenCages {
		chickenAge := util.GetChickenAgeByChickenCage(&chickenCage)

		var chickenHealthItems []entity.ChickenHealthItem
		err := tx.Model(&entity.ChickenHealthItem{}).Where("chicken_age = ? AND type = ?", chickenAge, enum.ChickenHealthItemTypeVaccineRoutine).Find(&chickenHealthItems).Error
		if err != nil {
			return err
		}

		if len(chickenHealthItems) == 0 {
			continue
		} else if chickenAge != uint64(chickenCage.LatestChickenAgeVaccineRoutine.Int64) && chickenCage.IsNeedRoutineVaccine {
			var vaccineRoutineNames []string
			for _, chickenHealthItem := range chickenHealthItems {
				vaccineRoutineNames = append(vaccineRoutineNames, chickenHealthItem.Name)
			}

			err = tx.Model(&entity.Notification{}).Create(&entity.Notification{
				CageId:              sql.NullInt64{Int64: int64(chickenCage.CageId), Valid: true},
				NotificationContext: pq.StringArray{constant.ChickenMonitoringNotificationContext},
				Description:         fmt.Sprintf(constant.VaccineRoutineNotification, strings.Join(vaccineRoutineNames, ","), chickenCage.Cage.Name),
			}).Error

			if err != nil {
				return err
			}
		}
	}

	s.log.Info("success check chicken cage if need routine vaccine")

	return nil
}

func (s *Scheduler) createNotificationTotalItemSaleShipToday(tx *gorm.DB) error {
	s.log.Info("create notification for total shipped today")

	today := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)

	var totalWarehouseSaleNeedShip int64
	err := tx.Model(&entity.WarehouseSale{}).Where("DATE(send_date) = ?", today).Count(&totalWarehouseSaleNeedShip).Error
	if err != nil {
		return err
	}

	var totalStoreSaleNeedShip int64
	err = tx.Model(&entity.StoreSale{}).Where("DATE(send_date) = ?", today).Count(&totalStoreSaleNeedShip).Error
	if err != nil {
		return err
	}

	if totalStoreSaleNeedShip > 0 {
		err = tx.Model(&entity.Notification{}).Create(&entity.Notification{
			NotificationContext: pq.StringArray{constant.StoreSaleNotificationContext},
			Description:         fmt.Sprintf(constant.ItemShipTodayWarehouseSaleNotification, totalStoreSaleNeedShip),
		}).Error

		if err != nil {
			return err
		}

		s.log.Info(fmt.Sprintf("total %d from store sale need ship today", totalStoreSaleNeedShip))
	}

	if totalWarehouseSaleNeedShip > 0 {
		err = tx.Model(&entity.Notification{}).Create(&entity.Notification{
			NotificationContext: pq.StringArray{constant.WarehouseSaleNotificationContext},
			Description:         fmt.Sprintf(constant.ItemShipTodayWarehouseSaleNotification, totalWarehouseSaleNeedShip),
		}).Error

		if err != nil {
			return err
		}

		s.log.Info(fmt.Sprintf("total %d from warehouse sale need ship today", totalWarehouseSaleNeedShip))
	}

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
					WarehouseId:         sql.NullInt64{Int64: int64(warehouseItem.WarehouseId), Valid: true},
					NotificationContext: pq.StringArray{constant.WarehouseItemNotificationContext},
					Description:         fmt.Sprintf(constant.WarehouseItemInDangerNotification, warehouseItem.Item.Name, warehouseItem.Warehouse.Name),
				})
			}
		}
	}

	err = tx.Model(&entity.Notification{}).CreateInBatches(&notifications, len(notifications)).Error
	if err != nil {
		return err
	}

	s.log.Error(fmt.Sprintf("success create %d notification if warehouse item in danger status", len(notifications)))
	return nil
}

func (s *Scheduler) createNotificationStoreItemGoodEggInDanger(tx *gorm.DB) error {
	s.log.Info("create notification store item good egg in danger")

	var goodEggItem entity.Item
	err := tx.Model(&entity.Item{}).Where("name = ? AND unit = ? AND type = ?", constant.GoodEgg, constant.UnitKg, enum.ItemCategoryEgg).First(&goodEggItem).Error
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
				StoreId:             sql.NullInt64{Int64: int64(storeItem.StoreId), Valid: true},
				NotificationContext: pq.StringArray{constant.StoreItemNotificationContext},
				Description:         fmt.Sprintf(constant.StoreItemInDangerNotification, storeItem.Item.Name, storeItem.Store.Name),
			})
		}
	}

	err = tx.Model(&entity.Notification{}).CreateInBatches(&notifications, len(notifications)).Error
	if err != nil {
		return err
	}

	return nil
}

func (s *Scheduler) Start() {
	s.cron.Start()
}

func (s *Scheduler) Stop() {
	s.cron.Stop()
}
