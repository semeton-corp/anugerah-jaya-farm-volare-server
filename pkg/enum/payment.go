package enum

type PaymentMethod uint8

const (
	PaymentMethodUnknown  PaymentMethod = 0
	PaymentMethodCash     PaymentMethod = 1
	PaymentMethodTransfer PaymentMethod = 2
)

var (
	PaymentMethodMap = map[PaymentMethod]string{
		PaymentMethodCash:     "Tunai",
		PaymentMethodTransfer: "Non Tunai",
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

type PaymentType uint8

const (
	PaymentTypeUnknown     PaymentType = 0
	PaymentTypePaidOff     PaymentType = 1
	PaymentTypeinstallment PaymentType = 2
)

var (
	PaymentTypeMap = map[PaymentType]string{
		PaymentTypePaidOff:     "Penuh",
		PaymentTypeinstallment: "Cicil",
	}
)

func (c PaymentType) String() string {
	return PaymentTypeMap[c]
}

func ValueOfPaymentType(value string) PaymentType {
	for k, v := range PaymentTypeMap {
		if v == value {
			return k
		}
	}
	return PaymentTypeUnknown
}

func (c PaymentType) IsValid() bool {
	_, ok := PaymentTypeMap[c]
	return ok
}

type SaleUnit uint8

const (
	SaleUnitUnknown SaleUnit = 0
	SaleUnitIkat    SaleUnit = 1
	SaleUnitPlastik SaleUnit = 2
	SaleUnitKg      SaleUnit = 3
)

var (
	SaleUnitMap = map[SaleUnit]string{
		SaleUnitIkat:    "Ikat",
		SaleUnitPlastik: "Plastik",
		SaleUnitKg:      "Kg",
	}
)

func (c SaleUnit) String() string {
	return SaleUnitMap[c]
}

func ValueOfSaleUnit(value string) SaleUnit {
	for k, v := range SaleUnitMap {
		if v == value {
			return k
		}
	}
	return SaleUnitUnknown
}

func (c SaleUnit) IsValid() bool {
	_, ok := SaleUnitMap[c]
	return ok
}
