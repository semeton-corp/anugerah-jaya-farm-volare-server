package enum

type PaymentMethod uint8

const (
	PaymentMethodUnknown     PaymentMethod = 0
	PaymentMethodPaidOff     PaymentMethod = 1
	PaymentMethodinstallment PaymentMethod = 2
)

var (
	PaymentMethodMap = map[PaymentMethod]string{
		PaymentMethodPaidOff:     "Lunas",
		PaymentMethodinstallment: "Cicil",
	}
)

func (c PaymentMethod) String() string {
	return PaymentMethodMap[c]
}

func ValueOfPaymentMethod(value string) PaymentMethod {
	for k, v := range PaymentMethodMap {
		if v == value {
			return k
		}
	}
	return PaymentMethodUnknown
}

func (c PaymentMethod) IsValid() bool {
	_, ok := PaymentMethodMap[c]
	return ok
}
