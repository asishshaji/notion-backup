package models

type GetTaskResult struct {
	ID        string `json:"id"`
	EventName string `json:"eventName"`
	Request   struct {
		SpaceID       string `json:"spaceId"`
		ExportOptions struct {
			ExportType string `json:"exportType"`
			TimeZone   string `json:"timeZone"`
			Locale     string `json:"locale"`
		} `json:"exportOptions"`
		ShouldExportComments bool `json:"shouldExportComments"`
	} `json:"request"`
	Actor struct {
		Table string `json:"table"`
		ID    string `json:"id"`
	} `json:"actor"`
	RootRequest struct {
		EventName string `json:"eventName"`
		RequestID string `json:"requestId"`
	} `json:"rootRequest"`
	Headers struct {
		IP                 string `json:"ip"`
		CityFromIP         string `json:"cityFromIp"`
		CountryCodeFromIP  string `json:"countryCodeFromIp"`
		Subdivision1FromIP string `json:"subdivision1FromIp"`
	} `json:"headers"`
	EqueuedAt int64 `json:"equeuedAt"`
	Status    struct {
		Type          string `json:"type"`
		PagesExported int    `json:"pagesExported"`
		ExportURL     string `json:"exportURL"`
	} `json:"status"`
	State string `json:"state"`
}

type GetTasksDTO struct {
	Results []GetTaskResult `json:"results"`
}
