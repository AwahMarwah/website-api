package user_address

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	cityShipping "website-api/model/city-shipping"
	masterModel "website-api/model/master"
	userAddressModel "website-api/model/user_address"
	"website-api/utils"

	"gorm.io/gorm"
)

func (s *service) Create(reqBody *userAddressModel.CreateUserAddressRequest) error {
	// check province
	province, err := s.masterRepo.TakeProvince([]string{"id", "name"}, &masterModel.Province{ID: reqBody.ProvinceID})
	if err != nil {
		return fmt.Errorf("province not found or invalid: %w", err)
	}

	// check city
	city, err := s.masterRepo.TakeCity([]string{"id", "name"}, &masterModel.City{ID: reqBody.CityID, ProvinceID: reqBody.ProvinceID})
	if err != nil {
		return fmt.Errorf("city not found or invalid for given province: %w", err)
	}

	// check district
	district, err := s.masterRepo.TakeDistrict([]string{"id", "name"}, &masterModel.District{ID: reqBody.DistrictID, CityID: reqBody.CityID})
	if err != nil {
		return fmt.Errorf("district not found or invalid for given city: %w", err)
	}

	// check subdistrict
	subdistrict, err := s.masterRepo.TakeSubdistrict([]string{"id", "name"}, &masterModel.Subdistrict{ID: reqBody.SubdistrictID, DistrictID: reqBody.DistrictID})
	if err != nil {
		return fmt.Errorf("subdistrict not found or invalid for given district: %w", err)
	}

	// checking city shipping
	cityID, err := strconv.ParseUint(reqBody.CityID, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid city_id: %w", err)
	}

	var destinationID int64

	existingShipping, err := s.masterRepo.TakeCityShippingMapping([]string{"city_id", "provider", "destination_id"}, &cityShipping.CityShippingMapping{CityID: cityID})
	if err == nil {
		// Mapping sudah ada
		destinationID = int64(existingShipping.DestinationId)

	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		// Mapping belum ada → cari RajaOngkir

		search := fmt.Sprintf(
			"%s, %s, %s, %s, %s",
			strings.ToUpper(subdistrict.Name),
			strings.ToUpper(district.Name),
			strings.ToUpper(city.Name),
			strings.ToUpper(province.Name),
			reqBody.PostalCode,
		)

		destinations, err := s.rajaOngkirProvider.SearchDestination(search)
		if err != nil {
			return fmt.Errorf(
				"failed to search RajaOngkir destination: %w",
				err,
			)
		}

		if len(destinations) == 0 {
			return fmt.Errorf(
				"destination not found for %s",
				search,
			)
		}

		for _, destination := range destinations {
			roProvince := utils.NormalizeRegionName(destination.ProvinceName)
			roCity := utils.NormalizeRegionName(destination.CityName)
			roDistrict := utils.NormalizeRegionName(destination.DistrictName)
			roSubdistrict := utils.NormalizeRegionName(destination.SubdistrictName)

			masterProvince := utils.NormalizeRegionName(province.Name)
			masterCity := utils.NormalizeRegionName(city.Name)
			masterDistrict := utils.NormalizeRegionName(district.Name)
			masterSubdistrict := utils.NormalizeRegionName(subdistrict.Name)

			fmt.Println("========== MATCH DEBUG ==========")

			fmt.Printf(
				"Province     : %q == %q => %v\n",
				roProvince,
				masterProvince,
				roProvince == masterProvince,
			)

			fmt.Printf(
				"City         : %q == %q => %v\n",
				roCity,
				masterCity,
				roCity == masterCity,
			)

			fmt.Printf(
				"District     : %q == %q => %v\n",
				roDistrict,
				masterDistrict,
				roDistrict == masterDistrict,
			)

			fmt.Printf(
				"Subdistrict  : %q == %q => %v\n",
				roSubdistrict,
				masterSubdistrict,
				roSubdistrict == masterSubdistrict,
			)

			fmt.Printf(
				"Postal Code  : %q == %q => %v\n",
				destination.ZipCode,
				reqBody.PostalCode,
				destination.ZipCode == reqBody.PostalCode,
			)

			if roProvince == masterProvince &&
				roCity == masterCity &&
				roDistrict == masterDistrict &&
				roSubdistrict == masterSubdistrict &&
				strings.TrimSpace(destination.ZipCode) ==
					strings.TrimSpace(reqBody.PostalCode) {

				fmt.Println("✅ DESTINATION MATCH:", destination.ID)

				destinationID = destination.ID
				break
			}
		}

		fmt.Println(destinationID, "ada harusnya")

		if destinationID == 0 {
			return fmt.Errorf(
				"failed matching RajaOngkir destination for %s",
				search,
			)
		}

		// Save mapping
		newShipping := cityShipping.CityShippingMapping{
			CityID:        cityID,
			Provider:      "rajaongkir",
			DestinationId: uint64(destinationID),
		}

		if err := s.masterRepo.CreateCityShippingMapping(&newShipping); err != nil {
			return fmt.Errorf(
				"failed to create shipping mapping: %w",
				err,
			)
		}

	} else {
		return fmt.Errorf(
			"failed to check shipping mapping: %w",
			err,
		)
	}

	// 7. Create user address
	userAddress := userAddressModel.UserAddress{
		UserID:        reqBody.UserID,
		RecipientName: reqBody.RecipientName,
		PhoneNumber:   reqBody.PhoneNumber,
		//FullAddress:           reqBody.FullAddress,
		ProvinceID:    reqBody.ProvinceID,
		CityID:        reqBody.CityID,
		DistrictID:    reqBody.DistrictID,
		SubdistrictID: reqBody.SubdistrictID,
		PostalCode:    reqBody.PostalCode,
		IsPrimary:     reqBody.IsPrimary,
		//ShippingProvider:      "rajaongkir",
		DestinationID: destinationID,
	}

	if err := s.userAddressRepo.Create(&userAddress); err != nil {
		return fmt.Errorf(
			"failed to create user address: %w",
			err,
		)
	}

	return nil
}
