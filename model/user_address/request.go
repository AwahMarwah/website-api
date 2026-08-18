package user_address

type (
	CreateUserAddressRequest struct {
		UserID        string `json:"user_id"`
		RecipientName string `binding:"required" json:"recipient_name"`
		PhoneNumber   string `binding:"required,min=12,max=12" json:"phone_number"`
		ProvinceID    string `binding:"required" json:"province_id"`
		CityID        string `binding:"required" json:"city_id"`
		DistrictID    string `binding:"required" json:"district_id"`
		SubdistrictID string `binding:"required" json:"subdistrict_id"`
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

	ReqUpdateUserAddress struct {
		Path struct {
			ID string `uri:"id" binding:"required"`
		}

		Body struct {
			RecepientName string `binding:"required,omitempty" json:"recepient_name"`
			PhoneNumber   string `binding:"required,min=12,max=12" json:"phone_number"`
			FullAddress   string `binding:"required,omitempty" json:"full_address"`
			City          string `binding:"required,omitempty" json:"city"`
			PostalCode    string `binding:"required,min=5,max=5" json:"postal_code"`
		}
	}
)
