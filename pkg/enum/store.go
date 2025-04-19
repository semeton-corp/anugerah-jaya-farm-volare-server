package enum

type PaymentMethod uint8

const (
	PaymentMethodUnknown     PaymentMethod = 0
	PaymentMethodPaidOff     PaymentMethod = 1
	PaymentMethodinstallment PaymentMethod = 2
)

var (
	PaymentMethodMap = map[PaymentMethod]string{
		PaymentMethodPaidOff:     "Penuh",
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

type PaymentStatus uint8

const (
	PaymentStatusUnknown PaymentStatus = 0
	PaymentStatusPaid    PaymentStatus = 1
	PaymentStatusUnpaid  PaymentStatus = 2
)

var (
	PaymentStatusMap = map[PaymentStatus]string{
		PaymentStatusPaid:   "Lunas",
		PaymentStatusUnpaid: "Belum Lunas",
	}
)

func (c PaymentStatus) String() string {
	return PaymentStatusMap[c]
}

func ValueOfPaymentStatus(value string) PaymentStatus {
	for k, v := range PaymentStatusMap {
		if v == value {
			return k
		}
	}
	return PaymentStatusUnknown
}

func (c PaymentStatus) IsValid() bool {
	_, ok := PaymentStatusMap[c]
	return ok
}
