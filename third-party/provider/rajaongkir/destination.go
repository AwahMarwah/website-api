package rajaongkir

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func (c *Client) SearchDestination(search string) ([]Destination, error) {
	fmt.Println("MASUK KE THIRD PARTY RAJAONGKIR")
	fmt.Println("API KEY EXIST:", c.apiKey != "")
	fmt.Println("BASE URL:", c.baseURL)
	endpoint := fmt.Sprintf(
		"%s/destination/domestic-destination?search=%s",
		c.baseURL,
		url.QueryEscape(search),
	)

	req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
	req.Header.Set("key", c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result DestinationResponse
	err = json.NewDecoder(resp.Body).Decode(&result)

	return result.Data, err
}
