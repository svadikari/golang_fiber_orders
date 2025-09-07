package schemas

type Product struct {
	Name        string  `json:"name" validate:"required,min=2,max=200"`
	Description string  `json:"description" validate:"required,min=5,max=1000"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Stock       int     `json:"stock"`
	ImageURL    string  `json:"image_url"`
}
