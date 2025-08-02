package constant

const (
	EggMonitoringStatusSafety string = "Aman"
	EggMonitoringStatusCheck  string = "Periksa"
	EggMonitoringStatusUrgent string = "Urgent"

	UnitKg      string = "Kg"
	UnitButir   string = "Butir"
	UnitIkat    string = "Ikat"
	UnitKarpet  string = "Karpet"
	UnitPlastik string = "Plastik"

	GoodEgg    string = "Telur OK" // display in store and warehouse
	RejectEgg  string = "Telur Reject"
	CrackedEgg string = "Telur Retak"  // display in store
	BrokenEgg  string = "Telur Bonyok" // display in store

	Corn string = "Jagung"

	TotalEggPerKarpet uint64 = 30 // butir
	TotalEggPerIkat   uint64 = 15 // kg
)
