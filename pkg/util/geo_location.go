package util

import (
	"math"
)

func IsWithinRadius(locationLon, locationLat, userLon, userLat, radiusKm float64) bool {
	const earthRadiusKm = 6371.0

	toRad := func(deg float64) float64 {
		return deg * (3.141592653589793 / 180.0)
	}

	dLat := toRad(userLat - locationLat)
	dLon := toRad(userLon - locationLon)

	lat1 := toRad(locationLat)
	lat2 := toRad(userLat)

	a := (sin(dLat/2) * sin(dLat/2)) +
		(sin(dLon/2)*sin(dLon/2))*cos(lat1)*cos(lat2)
	cVal := 2 * atan2(sqrt(a), sqrt(1-a))

	distance := earthRadiusKm * cVal
	
	return distance <= radiusKm
}

func sin(x float64) float64      { return float64(math.Sin(float64(x))) }
func cos(x float64) float64      { return float64(math.Cos(float64(x))) }
func atan2(y, x float64) float64 { return float64(math.Atan2(float64(y), float64(x))) }
func sqrt(x float64) float64     { return float64(math.Sqrt(float64(x))) }
