package enum

type CustomerType uint8

const (
	CustomerTypeUnknown CustomerType = 0
	CustomerTypeOld     CustomerType = 1
	CustomerTypeNew     CustomerType = 2
)

var (
	CustomerTypeMap = map[CustomerType]string{
		CustomerTypeOld: "Pelanggan Lama",
		CustomerTypeNew: "Pelanggan Baru",
	}
)

func (c CustomerType) String() string {
	return CustomerTypeMap[c]
}

func ValueOfCustomerType(value string) CustomerType {
	for k, v := range CustomerTypeMap {
		if v == value {
			return k
		}
	}
	return CustomerTypeUnknown
}

func (c CustomerType) IsValid() bool {
	_, ok := CustomerTypeMap[c]
	return ok
}
