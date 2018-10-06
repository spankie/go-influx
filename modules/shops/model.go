package shops

type Shop struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Price       float32 `json:"price"`
	Image       string  `json:"image"`
	Description string  `json:"description"`
}

func (shop Shop) GetAll() []Shop {
	return []Shop{
		Shop{ID: 0, Name: "Watch", Price: 56.2, Image: "/assets/img/watch.jpeg", Description: "Nice Watch"},
		Shop{ID: 1, Name: "Camera", Price: 34.2, Image: "/assets/img/camera.jpeg", Description: "Nice Camera"},
		Shop{ID: 2, Name: "Glass", Price: 24.2, Image: "/assets/img/glass.jpeg", Description: "Nice Glass"},
		Shop{ID: 3, Name: "Toy", Price: 56.2, Image: "/assets/img/toy.jpeg", Description: "Nice Toy"},
	}
}
