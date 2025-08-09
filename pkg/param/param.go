package param

import (
	"time"

	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/errx"
)

type DateParam time.Time

func (cd *DateParam) UnmarshalText(text []byte) error {
	parsedTime, err := time.Parse("02-01-2006", string(text))
	if err != nil {
		return err
	}
	*cd = DateParam(parsedTime)
	return nil
}

func (cd DateParam) Value() time.Time {
	return time.Time(cd)
}

type ItemCategoryParam enum.ItemCategory

func (w *ItemCategoryParam) UnmarshalText(text []byte) error {
	parsedCategory := enum.ValueOfItemCategory(string(text))
	if !parsedCategory.IsValid() {
		return errx.BadRequest("invalid warehouse item category")
	}

	*w = ItemCategoryParam(parsedCategory)
	return nil
}

func (w ItemCategoryParam) Value() enum.ItemCategory {
	return enum.ItemCategory(w)
}

type PaymentMethodParam enum.PaymentMethod

func (p *PaymentMethodParam) UnmarshalText(text []byte) error {
	parsedMethod := enum.ValueOfPaymentMethod(string(text))
	if !parsedMethod.IsValid() {
		return errx.BadRequest("invalid payment method")
	}

	*p = PaymentMethodParam(parsedMethod)
	return nil
}

func (p PaymentMethodParam) Value() enum.PaymentMethod {
	return enum.PaymentMethod(p)
}

type MonthParam enum.Month

func (p *MonthParam) UnmarshalText(text []byte) error {
	parsedFilter := enum.ValueOfMonth(string(text))
	if !parsedFilter.IsValid() {
		return errx.BadRequest("invalid month filter")
	}

	*p = MonthParam(parsedFilter)
	return nil
}

func (p MonthParam) Value() enum.Month {
	return enum.Month(p)
}

type OverviewGraphTimeParam enum.OverviewGraphTime

func (p *OverviewGraphTimeParam) UnmarshalText(text []byte) error {
	parsedFilter := enum.ValueOfOverviewGraphFilter(string(text))
	if !parsedFilter.IsValid() {
		return errx.BadRequest("invalid overview graph filter")
	}

	*p = OverviewGraphTimeParam(parsedFilter)
	return nil
}

func (p OverviewGraphTimeParam) Value() enum.OverviewGraphTime {
	return enum.OverviewGraphTime(p)
}

type ChickenHealthItemTypeParam enum.ChickenHealthItemType

func (p *ChickenHealthItemTypeParam) UnmarshalText(text []byte) error {
	parsedFilter := enum.ValueOfChickenHealthItemType(string(text))
	if !parsedFilter.IsValid() {
		return errx.BadRequest("invalid chicken health item type")
	}

	*p = ChickenHealthItemTypeParam(parsedFilter)
	return nil
}

func (p ChickenHealthItemTypeParam) Value() enum.ChickenHealthItemType {
	return enum.ChickenHealthItemType(p)
}

type LocationWorkTypeParam enum.LocationWorkType

func (p *LocationWorkTypeParam) UnmarshalText(text []byte) error {
	parsedFilter := enum.ValueOfLocationWorkType(string(text))
	if !parsedFilter.IsValid() {
		return errx.BadRequest("invalid location work type")
	}

	*p = LocationWorkTypeParam(parsedFilter)
	return nil
}

func (p LocationWorkTypeParam) Value() enum.LocationWorkType {
	return enum.LocationWorkType(p)
}

type PresenceStatusParam enum.PresenceStatus

func (p *PresenceStatusParam) UnmarshalText(text []byte) error {
	parsedFilter := enum.ValueOfPresenceStatus(string(text))
	if !parsedFilter.IsValid() {
		return errx.BadRequest("invalid presence status")
	}

	*p = PresenceStatusParam(parsedFilter)
	return nil
}

func (p PresenceStatusParam) Value() enum.PresenceStatus {
	return enum.PresenceStatus(p)
}

type PaymentStatusParam enum.PaymentStatus

func (p *PaymentStatusParam) UnmarshalText(text []byte) error {
	parsedFilter := enum.ValueOfPaymentStatus(string(text))
	if !parsedFilter.IsValid() {
		return errx.BadRequest("invalid payment status")
	}

	*p = PaymentStatusParam(parsedFilter)
	return nil
}

func (p PaymentStatusParam) Value() enum.PaymentStatus {
	return enum.PaymentStatus(p)
}

type SupplierTypeParam enum.SupplierType

func (p *SupplierTypeParam) UnmarshalText(text []byte) error {
	parsedFilter := enum.ValueOfSupplierType(string(text))
	if !parsedFilter.IsValid() {
		return errx.BadRequest("invalid supplier type")
	}

	*p = SupplierTypeParam(parsedFilter)
	return nil
}

func (p SupplierTypeParam) Value() enum.SupplierType {
	return enum.SupplierType(p)
}
