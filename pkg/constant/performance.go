package constant

const (
	KPIStatusGood = "Baik"
	KPIStatusMid  = "Cukup"
	KPIStatusBad  = "Buruk"

	KPIScoreGood = 90
	KPIScoreMid  = 74
	KPIScoreBad  = 73

	ThresholdKpiChicken = 0.75

	// EpeiTarget is the European Production Efficiency Index target a flock is benchmarked
	// against. KPI score = EPEI / EpeiTarget (1.0 = on target).
	EpeiTarget = 250.0
)
