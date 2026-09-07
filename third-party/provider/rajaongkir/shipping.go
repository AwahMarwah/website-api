package rajaongkir

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// CalculateCost menghitung ongkos kirim via endpoint cost.
// Mendukung dua format response:
//   - Standar RajaOngkir: {"rajaongkir":{"results":[...]}}
//   - Komerce-style:      {"meta":..., "data":[...]}
func (c *Client) CalculateCost(req CalculateCostRequest) ([]ShippingOption, error) {
	form := url.Values{}
	form.Set("origin", strconv.FormatInt(req.Origin, 10))
	form.Set("destination", strconv.FormatInt(req.Destination, 10))
	form.Set("weight", strconv.Itoa(req.Weight))
	form.Set("courier", req.Courier)

	httpReq, err := http.NewRequest(http.MethodPost, c.costURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("key", c.apiKey)
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("rajaongkir cost returned status %d: %s", resp.StatusCode, string(body))
	}

	// Coba decode format standar RajaOngkir dulu
	var standard RajaOngkirCostResponse
	if err := json.Unmarshal(body, &standard); err == nil && len(standard.RajaOngkir.Results) > 0 {
		return flattenStandardCost(standard.RajaOngkir.Results), nil
	}

	// Fallback format Komerce-style: {"data":[...]}
	var komerce CalculateCostResponse
	if err := json.Unmarshal(body, &komerce); err == nil && len(komerce.Data) > 0 {
		return komerce.Data, nil
	}

	// Tidak terdecode: tampilkan raw body untuk debug
	return nil, fmt.Errorf("unable to parse rajaongkir cost response: %s", string(body))
}

func flattenStandardCost(groups []RawCostGroup) []ShippingOption {
	options := make([]ShippingOption, 0)
	for _, group := range groups {
		for _, svc := range group.Costs {
			cost := int64(0)
			etd := ""
			if len(svc.Cost) > 0 {
				cost = svc.Cost[0].Value
				etd = svc.Cost[0].Etd
			}
			options = append(options, ShippingOption{
				Name:        group.Name,
				Code:        group.Code,
				Service:     svc.Service,
				Description: svc.Description,
				Cost:        cost,
				Etd:         etd,
			})
		}
	}
	return options
}