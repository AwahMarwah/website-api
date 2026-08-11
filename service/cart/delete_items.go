package cart

import cartModel "website-api/model/cart"

func (s *service) DeleteItems(req *cartModel.DeleteCartRequest) error {
	return s.cartRepo.DeleteItems(req.UserID, req.VariantIDs)
}
