package review

import "time"

type (
	ReviewResponse struct {
		ID        string    `json:"id"`
		ProductID string    `json:"product_id"`
		UserName  string    `json:"user_name"`
		Rating    int       `json:"rating"`
		Comment   string    `json:"comment"`
		CreatedAt time.Time `json:"created_at"`
	}

	SwaggerReviewList struct {
		Data    []ReviewResponse `json:"data"`
		Message string           `json:"message"`
		Page    struct {
			Current int `json:"current"`
			Size    int `json:"size"`
			Total   int `json:"total"`
		} `json:"page"`
	}
)