package constant

const (
	EggMonitoringStatusSafety string = "Aman"
	EggMonitoringStatusCheck  string = "Periksa"
	EggMonitoringStatusUrgent string = "Urgent"

	EggUnitKg      string = "Kg"
	EggUnitButir   string = "Butir"
	EggUnitIkat    string = "Ikat"
	EggUnitKarpet  string = "Karpet"
	EggUnitPlastik string = "Plastik"

	GoodEgg    string = "Telur OK" // display in store and warehouse
	RejectEgg  string = "Telur Reject"
	CrackedEgg string = "Telur Retak"  // display in store
	BrokenEgg  string = "Telur Bonyok" // display in store

	TotalEggPerKarpet uint64 = 30 // butir
	TotalEggPerIkat   uint64 = 15 // kg
)
