package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sangnt1552314/ip-tracker/internal/models"
)

func GetGeoLocationData() (*models.CloudflareGeoLocationDetail, error){
	url := "https://speed.cloudflare.com/meta"

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned non-200 status code: %d", resp.StatusCode)
	}

	var cloudflareData models.CloudflareGeoLocationDetail
	if err := json.NewDecoder(resp.Body).Decode(&cloudflareData); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &cloudflareData, nil
}