package models

type CloudflareGeoLocationDetail struct {
	Hostname string `json:"hostname"`
	ClientIp string `json:"clientIp"`
	HttpProtocol string `json:"httpProtocol"`
	Asn int `json:"asn"`
	AsOrganization string `json:"asOrganization"`
	Colo string `json:"colo"`
	Country string `json:"country"`
	City string `json:"city"`
	Region string `json:"region"`
	PostalCode string `json:"postalCode"`
	Latitude string `json:"latitude"`
	Longitude string `json:"longitude"`
}