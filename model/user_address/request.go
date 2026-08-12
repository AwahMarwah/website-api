package user_address

type (
	CreateUserAddressRequest struct {
		UserID        string `json:"user_id"`
		RecipientName string `binding:"required" json:"recipient_name"`
		PhoneNumber   string `binding:"required,min=12,max=12" json:"phone_number"`
		FullAddress   string `binding:"required" json:"full_address"`
		City          string `binding:"required" json:"city"`
		PostalCode    string `binding:"required,min=5,max=5" json:"postal_code"`
		IsPrimary     bool   `binding:"required" json:"is_primary"`
	}

	ReqPath struct {
		ID string `uri:"id" binding:"required"`
	}

	UpdateUserAddress struct {
		Path struct {
			ID string `uri:"id" binding:"required"`
		}

		Body struct {
			IsPrimary bool `json:"is_primary"`
		}
	}
)
