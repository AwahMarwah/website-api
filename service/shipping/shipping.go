package shipping

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	merchantModel "website-api/model/merchant"
	shippingModel "website-api/model/shipping"
	userAddressModel "website-api/model/user_address"
	"website-api/repository/merchant"
	"website-api/repository/product"
	product_variant "website-api/repository/product-variant"
	userAddressRepo "website-api/repository/user_address"
	"website-api/third-party/provider/rajaongkir"

	"gorm.io/gorm"
)

type (
	IService interface {
		Quote(req *shippingModel.ShippingCostRequest, userID string) ([]shippingModel.MerchantShipping, int, error)
	}

	service struct {
		userAddressRepo    userAddressRepo.IRepo
		productVariantRepo product_variant.IRepo
		productRepo        product.IRepo
		merchantRepo       merchant.IRepo
		rajaOngkir         rajaongkir.Provider
	}
)

func NewService(
	userAddressRepo userAddressRepo.IRepo,
	productVariantRepo product_variant.IRepo,
	productRepo product.IRepo,
	merchantRepo merchant.IRepo,
	rajaOngkir rajaongkir.Provider,
) IService {
	return &service{
		userAddressRepo:    userAddressRepo,
		productVariantRepo: productVariantRepo,
		productRepo:        productRepo,
		merchantRepo:       merchantRepo,
		rajaOngkir:         rajaOngkir,
	}
}

var defaultCouriers = []string{"jne", "pos", "tiki"}

// Quote menghitung ongkir per merchant untuk daftar item checkout.
func (s *service) Quote(req *shippingModel.ShippingCostRequest, userID string) ([]shippingModel.MerchantShipping, int, error) {
	// Ambil alamat tujuan
	address, err := s.userAddressRepo.Take([]string{"id", "user_id", "destination_id"}, &userAddressModel.UserAddress{ID: req.AddressID})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, http.StatusNotFound, fmt.Errorf("address not found")
		}
		return nil, http.StatusInternalServerError, fmt.Errorf("gagal mengambil alamat: %w", err)
	}

	// Ownership check
	if address.UserID != userID {
		return nil, http.StatusForbidden, fmt.Errorf("forbidden")
	}
	if address.DestinationID == 0 {
		return nil, http.StatusBadRequest, fmt.Errorf("alamat belum memiliki destination_id")
	}

	couriers := req.Couriers
	if len(couriers) == 0 {
		couriers = defaultCouriers
	}

	merchantByID := map[string]merchantModel.Merchant{}
	weightByMerchant := map[string]int{}

	for _, item := range req.Items {
		variant, err := s.productVariantRepo.FindByID(item.VariantID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, http.StatusBadRequest, fmt.Errorf("variant with id %s not found", item.VariantID)
			}
			return nil, http.StatusInternalServerError, fmt.Errorf("gagal mengambil variant: %w", err)
		}

		productInfo, err := s.productRepo.FindByID(variant.ProductID)
		if err != nil {
			return nil, http.StatusInternalServerError, fmt.Errorf("gagal mengambil produk: %w", err)
		}
		if productInfo.MerchantId == "" {
			return nil, http.StatusBadRequest, fmt.Errorf("produk %s belum memiliki merchant", productInfo.Name)
		}

		if _, ok := merchantByID[productInfo.MerchantId]; !ok {
			m, err := s.merchantRepo.FindByID(productInfo.MerchantId)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, http.StatusBadRequest, fmt.Errorf("merchant %s not found", productInfo.MerchantId)
				}
				return nil, http.StatusInternalServerError, fmt.Errorf("gagal mengambil merchant: %w", err)
			}
			merchantByID[m.ID] = m
		}

		// variant.Weight dalam kg → gram, kalikan dengan qty
		weightGramPerUnit := int(math.Round(float64(variant.Weight) * 1000))
		weightByMerchant[productInfo.MerchantId] += weightGramPerUnit * item.Qty
	}

	result := make([]shippingModel.MerchantShipping, 0, len(weightByMerchant))
	for merchantID, weightGram := range weightByMerchant {
		merch := merchantByID[merchantID]
		weightKg := int(math.Ceil(float64(weightGram) / 1000.0))

		options := make([]shippingModel.ShippingOption, 0)
		for _, courier := range couriers {
			roOptions, err := s.rajaOngkir.CalculateCost(rajaongkir.CalculateCostRequest{
				Origin:      merch.DestinationID,
				Destination: address.DestinationID,
				Weight:      weightGram, // API ongkir menerima satuan gram
				Courier:     courier,
			})
			if err != nil {
				continue // skip courier bila gagal
			}
			for _, opt := range roOptions {
				options = append(options, shippingModel.ShippingOption{
					Courier: opt.Code,
					Service: opt.Service,
					Cost:    opt.Cost,
					Etd:     opt.Etd,
				})
			}
		}

		result = append(result, shippingModel.MerchantShipping{
			MerchantID:   merchantID,
			MerchantName: merch.Name,
			WeightGram:   weightGram,
			WeightKg:     weightKg,
			Options:      options,
		})
	}

	return result, http.StatusOK, nil
}