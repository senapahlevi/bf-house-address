package user

// // EuclideanDistance computes the Euclidean distance between two points with latitude and longitude coordinates
// func EuclideanDistance(lat1, lon1, lat2, lon2 float64) float64 {
// 	// Earth radius in kilometers
// 	const R = 6371

// 	// Convert latitude and longitude to Cartesian coordinates
// 	p1 := r3.Vector{
// 		X: R * math.Cos(lat1) * math.Cos(lon1),
// 		Y: R * math.Cos(lat1) * math.Sin(lon1),
// 		Z: R * math.Sin(lat1),
// 	}

// 	p2 := r3.Vector{
// 		X: R * math.Cos(lat2) * math.Cos(lon2),
// 		Y: R * math.Cos(lat2) * math.Sin(lon2),
// 		Z: R * math.Sin(lat2),
// 	}

// 	// Compute the Euclidean distance between the points
// 	return p1.Sub(p2).Norm()
// }

// func main() {
// 	// Example points in Jakarta and Surabaya
// 	lat1 := -6.21462 * math.Pi / 180 // convert degrees to radians
// 	lon1 := 106.84513 * math.Pi / 180
// 	lat2 := -7.25747 * math.Pi / 180
// 	lon2 := 112.75209 * math.Pi / 180

// 	// Print the distance in kilometers
// 	fmt.Println(EuclideanDistance(lat1, lon1, lat2, lon2))
// }
