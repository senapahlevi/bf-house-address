// Define a house struct with longitude and latitude fields
type house struct {
	longitude float64
	latitude  float64
}

// Calculate the optimum route between a list of houses
func calculateOptimumRoute(houses []house) []house {
	// Implement your route optimization logic here
	// For simplicity, let's assume the optimum route is the one that minimizes the total distance between houses

	// Initialize an empty slice to store the optimum route
	optimumRoute := []house{}

	// Loop through the houses and find the nearest neighbor for each one
	for i := 0; i < len(houses); i++ {
		// Initialize the minimum distance and index variables
		minDistance := math.Inf(1)
		minIndex := -1

		// Loop through the remaining houses and calculate the distance to the current house
		for j := i + 1; j < len(houses); j++ {
			distance := calculateDistance(houses[i].longitude, houses[i].latitude, houses[j].longitude, houses[j].latitude)

			// Update the minimum distance and index if a smaller distance is found
			if distance < minDistance {
				minDistance = distance
				minIndex = j
			}
		}

		// Swap the current house with its nearest neighbor in the slice
		houses[i], houses[minIndex] = houses[minIndex], houses[i]

		// Append the current house to the optimum route slice
		optimumRoute = append(optimumRoute, houses[i])
	}

	// Return the optimum route slice
	return optimumRoute
}

func main() {
	// Define some sample houses with their longitude and latitude values
	h1 := house{longitude: 10.0, latitude: 20.0}
	h2 := house{longitude: 15.0, latitude: 25.0}
	h3 := house{longitude: 12.0, latitude: 22.0}
	h4 := house{longitude: 18.0, latitude: 28.0}

	// Create a slice of houses
	houses := []house{h1, h2, h3, h4}

	// Call the calculateOptimumRoute function and print the result
	optimumRoute := calculateOptimumRoute(houses)
	fmt.Println("The optimum route is:", optimumRoute)


	package main

import (
	"fmt"
	"github.com/golang-basic/golearn/knn"
)

func main() {
	// Define the longitude and latitude values of point A and point B
	pointALat := -6.1766841
	pointALon := 106.8306534
	pointBLat := -6.1701667
	pointBLon := 106.8241962

	// Create a slice of float64 for each point
	pointA := []float64{pointALat, pointALon}
	pointB := []float64{pointBLat, pointBLon}

	// Call the EuclideanDistance function from the knn package and print the result
	distance := knn.EuclideanDistance(pointA, pointB)
	fmt.Println("The distance between point A and point B is:", distance)
}