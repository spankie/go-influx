package shops

import (
	"errors"
	"log"
	"time"
)

// Product structure of my cool shop product
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Price       float32 `json:"price"`
	Image       string  `json:"image"`
	Description string  `json:"description"`
	Views       string
}

// ProductMeasurement schema of product measurement
type ProductMeasurement struct {
	ProductID   int `json:"id"`
	ProductName string
	time        time.Time
}

var allProducts = []Product{
	Product{ID: 0, Name: "Watch", Price: 56.2, Image: "/assets/img/watch.jpeg", Description: "Nice Watch"},
	Product{ID: 1, Name: "Camera", Price: 34.2, Image: "/assets/img/camera.jpeg", Description: "Nice Camera"},
	Product{ID: 2, Name: "Glass", Price: 24.2, Image: "/assets/img/glass.jpeg", Description: "Nice Glass"},
	Product{ID: 3, Name: "Toy", Price: 56.2, Image: "/assets/img/toy.jpeg", Description: "Nice Toy"},
}

// GetAll fetches all products in my shop
func (product Product) GetAll() []Product {
	return allProducts
}

// Get fetches a particular product in my shop identified By its ID
func (product Product) Get(id int) (Product, error) {
	if id < 0 && len(allProducts) >= id {
		log.Println(id)
		return Product{}, errors.New("No product found")
	}
	return allProducts[id], nil
}
