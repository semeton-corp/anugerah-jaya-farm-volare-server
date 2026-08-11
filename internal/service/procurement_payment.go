package service

import (
	"database/sql"
	"time"

	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/errx"
	"github.com/shopspring/decimal"
)

func validateInitialProcurementPayment(paymentType enum.PaymentType, totalPayment, totalPrice decimal.Decimal, paymentCount int) error {
	if totalPayment.GreaterThan(totalPrice) {
		return errx.BadRequest("total payment more than total price")
	}

	switch paymentType {
	case enum.PaymentTypePaidOff:
		if !totalPayment.Equal(totalPrice) {
			return errx.BadRequest("payment must equal total price for paid off")
		}
	case enum.PaymentTypeinstallment:
		if paymentCount == 0 || !totalPayment.GreaterThan(decimal.Zero) {
			return errx.BadRequest("payments are required for installment")
		}
	}

	return nil
}

func procurementPaymentStatus(totalPayment, totalPrice decimal.Decimal, allowOverpayment bool) (enum.PaymentStatus, sql.NullTime, error) {
	if totalPayment.GreaterThan(totalPrice) {
		if !allowOverpayment {
			return enum.PaymentStatusUnknown, sql.NullTime{}, errx.BadRequest("total payment more than total price")
		}

		return enum.PaymentStatusPaid, sql.NullTime{Time: time.Now(), Valid: true}, nil
	}

	if totalPayment.Equal(totalPrice) {
		return enum.PaymentStatusPaid, sql.NullTime{Time: time.Now(), Valid: true}, nil
	}

	if totalPayment.Equal(decimal.Zero) {
		return enum.PaymentStatusNotPaid, sql.NullTime{Valid: false}, nil
	}

	return enum.PaymentStatusUnpaid, sql.NullTime{Valid: false}, nil
}
