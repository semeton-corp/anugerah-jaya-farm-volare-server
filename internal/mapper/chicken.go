package mapper

import (
	"time"

	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/util"
	"github.com/shopspring/decimal"
)

func ChickenHealthItemToResponse(chickenHealthItem *entity.ChickenHealthItem) dto.ChickenHealthItemResponse {
	response := dto.ChickenHealthItemResponse{
		Id:   chickenHealthItem.Id,
		Name: chickenHealthItem.Name,
		Type: chickenHealthItem.Type.String(),
		Note: chickenHealthItem.Note,
	}

	if chickenHealthItem.ChickenAge.Valid {
		valUint64 := uint64(chickenHealthItem.ChickenAge.Int64)
		chickenCategory := util.GetChickenCategoryByAge(valUint64).String()

		response.ChickenAge = &valUint64
		response.ChickenCategory = &chickenCategory
	} else {
		response.ChickenAge = nil
		response.ChickenCategory = nil
	}

	return response
}

func ChickenHealthMonitoringToResponse(chickenHealthMonitoring *entity.ChickenHealthMonitoring) dto.ChickenHealthMonitoringResponse {
	response := dto.ChickenHealthMonitoringResponse{
		Id:             chickenHealthMonitoring.Id,
		HealthItemName: chickenHealthMonitoring.HealthItemName,
		Type:           chickenHealthMonitoring.Type.String(),
		Dose:           chickenHealthMonitoring.Dose,
		Unit:           chickenHealthMonitoring.Unit,
		CreatedAt:      chickenHealthMonitoring.CreatedAt.Format("02-01-2006"),
	}

	if chickenHealthMonitoring.Disease.Valid {
		response.Disease = chickenHealthMonitoring.Disease.String
	} else {
		response.Disease = "-"
	}

	response.ChickenAge = chickenHealthMonitoring.ChickenAge
	response.ChickenCategory = util.GetChickenCategoryByAge(chickenHealthMonitoring.ChickenAge).String()

	return response
}

func ChickenMonitoringToResponse(chickenMonitoring *entity.ChickenMonitoring) dto.ChickenMonitoringResponse {
	return ChickenMonitoringToResponseAt(chickenMonitoring, chickenMonitoring.CreatedAt)
}

func ChickenMonitoringToResponseAt(chickenMonitoring *entity.ChickenMonitoring, asOf time.Time) dto.ChickenMonitoringResponse {
	totalLiveChicken := chickenMonitoringLiveChicken(chickenMonitoring)
	chickenCage := chickenCageToResponseAt(&chickenMonitoring.ChickenCage, asOf)
	chickenCage.TotalChicken = totalLiveChicken

	return dto.ChickenMonitoringResponse{
		Id:                 chickenMonitoring.Id,
		ChickenCage:        chickenCage,
		TotalLiveChicken:   totalLiveChicken,
		TotalSickChicken:   chickenMonitoring.TotalSickChicken,
		TotalDeatchChicken: chickenMonitoring.TotalDeathChicken,
		TotalFeed:          chickenMonitoring.TotalFeed,
		Note:               chickenMonitoring.Note,
	}
}

func ChickenMonitoringToListResponse(chickenMonitoring *entity.ChickenMonitoring) dto.ChickenMonitoringListResponse {
	return ChickenMonitoringToListResponseAt(chickenMonitoring, chickenMonitoring.CreatedAt)
}

func ChickenMonitoringToListResponseAt(chickenMonitoring *entity.ChickenMonitoring, asOf time.Time) dto.ChickenMonitoringListResponse {
	totalLiveChicken := chickenMonitoringLiveChicken(chickenMonitoring)
	chickenCage := chickenCageToResponseAt(&chickenMonitoring.ChickenCage, asOf)
	chickenCage.TotalChicken = totalLiveChicken

	response := dto.ChickenMonitoringListResponse{
		Id:                chickenMonitoring.Id,
		ChickenCage:       chickenCage,
		TotalLiveChicken:  totalLiveChicken,
		TotalSickChicken:  chickenMonitoring.TotalSickChicken,
		TotalDeathChicken: chickenMonitoring.TotalDeathChicken,
		TotalFeed:         chickenMonitoring.TotalFeed,
		MortalityRate:     chickenMonitoring.MortalityRate,
	}

	return response
}

func chickenMonitoringLiveChicken(chickenMonitoring *entity.ChickenMonitoring) uint64 {
	if chickenMonitoring.TotalChicken < chickenMonitoring.TotalDeathChicken {
		return 0
	}

	return chickenMonitoring.TotalChicken - chickenMonitoring.TotalDeathChicken
}

func ChickenProcurementDraftToResponse(data *entity.ChickenProcurementDraft) dto.ChickenProcurementDraftResponse {
	return dto.ChickenProcurementDraftResponse{
		Id:         data.Id,
		Cage:       CageToResponse(&data.Cage),
		Supplier:   SupplierToResponse(&data.Supplier),
		Quantity:   data.Quantity,
		TotalPrice: data.TotalPrice.String(),
		InputDate:  data.CreatedAt.Format("02-01-2006"),
	}
}

func AfkirChickenCustomerToListResponse(data *entity.AfkirChickenCustomer) dto.AfkirChickenCustomerListResponse {
	return dto.AfkirChickenCustomerListResponse{
		Id:          data.Id,
		Name:        data.Name,
		PhoneNumber: data.PhoneNumber,
		Address:     data.Address,
		LatestPrice: data.LatestPrice.String(),
	}
}

func AfkirChickenCustomerToResponse(data *entity.AfkirChickenCustomer) dto.AfkirChickenCustomerResponse {
	response := dto.AfkirChickenCustomerResponse{Id: data.Id,
		Name:        data.Name,
		PhoneNumber: data.PhoneNumber,
		Address:     data.Address,
		LatestPrice: data.LatestPrice.String(),
	}

	afkirChickenSales := make([]dto.AfkirChickenSaleListResponse, 0)
	for _, e := range data.AfkirChickenSales {
		afkirChickenSales = append(afkirChickenSales, AfkirChickenSaleToListResponse(&e))
	}

	response.AfkirChickenSales = afkirChickenSales
	return response
}

func AfkirChickenSaleToListResponse(data *entity.AfkirChickenSale) dto.AfkirChickenSaleListResponse {
	response := dto.AfkirChickenSaleListResponse{
		Id:                   data.Id,
		SellDate:             data.CreatedAt.Format("02-01-2006"),
		AfkirChickenCustomer: AfkirChickenCustomerToListResponse(&data.AfkirChickenCustomer),
		ChickenAge:           data.ChickenAge,
		TotalSellChicken:     data.TotalSellChicken,
		PricePerChicken:      data.PricePerChicken.String(),
		TotalPrice:           data.TotalPrice.String(),
		PaymentStatus:        data.PaymentStatus.String(),
		TakenAt:              data.TakenAt.Format("02-01-2006"),
		IsTaken:              data.IsTaken,
	}

	if data.DeadlinePaymentDate.Valid {
		response.DeadlinePaymentDate = data.DeadlinePaymentDate.Time.Format("02-01-2006")
		if time.Now().After(data.DeadlinePaymentDate.Time) {
			response.IsMoreThanDeadlinePaymentDate = true
		} else {
			response.IsMoreThanDeadlinePaymentDate = false
		}
	} else {
		response.DeadlinePaymentDate = "-"
		response.IsMoreThanDeadlinePaymentDate = false
	}

	if data.PaidDate.Valid {
		response.PaidDate = data.PaidDate.Time.Format("02-01-2006")
	} else {
		response.PaidDate = "-"
	}

	return response
}

func AfkirChickenSaleToResponse(data *entity.AfkirChickenSale) dto.AfkirChickenSaleResponse {
	response := dto.AfkirChickenSaleResponse{
		Id:                   data.Id,
		SellDate:             data.CreatedAt.Format("02-01-2006"),
		AfkirChickenCustomer: AfkirChickenCustomerToListResponse(&data.AfkirChickenCustomer),
		ChickenCage:          ChickenCageToResponse(&data.ChickenCage),
		ChickenAge:           data.ChickenAge,
		TotalSellChicken:     data.TotalSellChicken,
		PricePerChicken:      data.PricePerChicken.String(),
		TotalPrice:           data.PricePerChicken.Mul(decimal.NewFromUint64(data.TotalSellChicken)).String(),
		PaymentStatus:        data.PaymentStatus.String(),
		PaymentType:          data.PaymentType.String(),
		TakenAt:              data.TakenAt.Format("02-01-2006"),
		IsTaken:              data.IsTaken,
	}

	if data.DeadlinePaymentDate.Valid {
		response.DeadlinePaymentDate = data.DeadlinePaymentDate.Time.Format("02-01-2006")
		if time.Now().After(data.DeadlinePaymentDate.Time) {
			response.IsMoreThanDeadlinePaymentDate = true
		} else {
			response.IsMoreThanDeadlinePaymentDate = false
		}
	} else {
		response.DeadlinePaymentDate = "-"
		response.IsMoreThanDeadlinePaymentDate = false
	}

	if data.PaidDate.Valid {
		response.PaidDate = data.PaidDate.Time.Format("02-01-2006")
	} else {
		response.PaidDate = "-"
	}

	return response
}

func AfkirChickenSaleDraftToResponse(data *entity.AfkirChickenSaleDraft) dto.AfkirChickenSaleDraftResponse {
	return dto.AfkirChickenSaleDraftResponse{
		Id:                   data.Id,
		InputDate:            data.CreatedAt.Format("02-01-2006"),
		ChickenCage:          ChickenCageToResponse(&data.ChickenCage),
		AfkirChickenCustomer: AfkirChickenCustomerToListResponse(&data.AfkirChickenCustomer),
		TotalSellChicken:     data.TotalSellChicken,
		PricePerChicken:      data.PricePerChicken.String(),
		TotalPrice:           data.TotalPrice.String(),
	}
}

// Note : without payment and remaining payment
func ChickenProcurementToResponse(data *entity.ChickenProcurement) dto.ChickenProcurementResponse {
	response := dto.ChickenProcurementResponse{
		Id:                    data.Id,
		OrderDate:             data.CreatedAt.Format("02-01-2006"),
		Quantity:              data.Quantity,
		Cage:                  CageToResponse(&data.Cage),
		Supplier:              SupplierToListResponse(&data.Supplier),
		EstimationArrivalDate: data.EstimationArrivalDate.Format("02-01-2006"),
		PaymentStatus:         data.PaymentStatus.String(),
		TotalPrice:            data.TotalPrice.String(),
		PaymentType:           data.PaymentType.String(),
		IsArrived:             data.IsArrived,
		ProcurementStatus:     data.Status.String(),
		Note:                  data.Note,
	}

	if data.DeadlinePaymentDate.Valid {
		response.DeadlinePaymentDate = data.DeadlinePaymentDate.Time.Format("02-01-2006")
		if time.Now().After(data.DeadlinePaymentDate.Time) {
			response.IsMoreThanDeadlinePaymentDate = true
		} else {
			response.IsMoreThanDeadlinePaymentDate = false
		}
	} else {
		response.DeadlinePaymentDate = "-"
		response.IsMoreThanDeadlinePaymentDate = false
	}

	if data.PaidDate.Valid {
		response.PaidDate = data.PaidDate.Time.Format("02-01-2006")
	} else {
		response.PaidDate = "-"
	}

	if data.ReceiveQuantity.Valid {
		val := uint64(data.ReceiveQuantity.Int64)
		response.ReceiveQuantity = &val
	}

	return response
}

func ChickenProcurementToListResponse(data *entity.ChickenProcurement) dto.ChickenProcurementListResponse {
	response := dto.ChickenProcurementListResponse{
		Id:                    data.Id,
		OrderDate:             data.CreatedAt.Format("02-01-2006"),
		Cage:                  CageToResponse(&data.Cage),
		Quantity:              data.Quantity,
		Supplier:              SupplierToListResponse(&data.Supplier),
		EstimationArrivalDate: data.EstimationArrivalDate.Format("02-01-2006"),
		PaymentStatus:         data.PaymentStatus.String(),
		IsArrived:             data.IsArrived,
		PaymentType:           data.PaymentType.String(),
		TotalPrice:            data.TotalPrice.String(),
		ProcurementStatus:     data.Status.String(),
	}

	if data.DeadlinePaymentDate.Valid {
		response.DeadlinePaymentDate = data.DeadlinePaymentDate.Time.Format("02-01-2006")
		if time.Now().After(data.DeadlinePaymentDate.Time) {
			response.IsMoreThanDeadlinePaymentDate = true
		} else {
			response.IsMoreThanDeadlinePaymentDate = false
		}
	} else {
		response.DeadlinePaymentDate = "-"
		response.IsMoreThanDeadlinePaymentDate = false
	}

	if data.PaidDate.Valid {
		response.PaidDate = data.PaidDate.Time.Format("02-01-2006")
	} else {
		response.PaidDate = "-"
	}

	if data.ReceiveQuantity.Valid {
		val := uint64(data.ReceiveQuantity.Int64)
		response.ReceiveQuantity = &val
	}

	return response
}

// Note : Without remaining
func ChickenProcurementPaymentToResponse(data *entity.ChickenProcurementPayment) dto.ChickenProcurementPaymentResponse {
	return dto.ChickenProcurementPaymentResponse{
		Id:            data.Id,
		Date:          data.PaymentDate.Format("02-01-2006"),
		Nominal:       data.Nominal.String(),
		PaymentMethod: data.PaymentMethod.String(),
		PaymentProof:  data.PaymentProof,
	}
}

// Note : without remaining
func AfkirChickenSalePaymentToResponse(data *entity.AfkirChickenSalePayment) dto.AfkirChickenSalePaymentResponse {
	return dto.AfkirChickenSalePaymentResponse{
		Id:            data.Id,
		Date:          data.PaymentDate.Format("02-01-2006"),
		Nominal:       data.Nominal.String(),
		PaymentMethod: data.PaymentMethod.String(),
		PaymentProof:  data.PaymentProof,
	}
}

func ChickenPerformanceToResponse(data *entity.ChickenPerformance) dto.ChickenPerformanceResponse {
	return dto.ChickenPerformanceResponse{
		CageName:                     data.CageName,
		ChickenCategory:              data.ChickenCategory.String(),
		ChickenAge:                   data.ChickenAge,
		TotalChicken:                 data.TotalChicken,
		TotalEgg:                     data.TotalEgg,
		AverageConsumptionPerChicken: data.AverageConsumptionPerChicken,
		AverageWeightPerEgg:          data.AverageWeightPerEgg,
		FCR:                          data.FCR,
		HDP:                          data.HDP * 100.0,
		MortalityRate:                data.MortalityRate * 100.0,
		Productivity:                 data.Productivity.String(),
	}
}
