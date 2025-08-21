package mapper

import (
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
)

func UserCashAdvancePaymentToResponse(data *entity.UserCashAdvancePayment) dto.UserCashAdvancePaymentResponse {
	return dto.UserCashAdvancePaymentResponse{
		Id:            data.Id,
		Date:          data.CreatedAt.Format("02 Jan 2006"),
		Nominal:       data.Nominal.String(),
		PaymentMethod: data.PaymentMethod.String(),
		PaymentProof:  data.PaymentProof,
	}
}
