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

type WarehouseItemCategoryParam enum.WarehouseItemCategory

func (w *WarehouseItemCategoryParam) UnmarshalText(text []byte) error {
	parsedCategory := enum.ValueOfWarehouseItemCategory(string(text))
	if !parsedCategory.IsValid() {
		return errx.BadRequest("invalid warehouse item category")
	}

	*w = WarehouseItemCategoryParam(parsedCategory)
	return nil
}

func (w WarehouseItemCategoryParam) Value() enum.WarehouseItemCategory {
	return enum.WarehouseItemCategory(w)
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
