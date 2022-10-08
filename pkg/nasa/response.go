package nasa

type Response interface {
	APODResponse | NEOResponse
}

type APODResponse struct {
	Copyright      string `json:"copyright"`
	Date           string `json:"date"`
	Explanation    string `json:"explanation"`
	HDUrl          string `json:"hdurl"`
	MediaType      string `json:"media_type"`
	ServiceVersion string `json:"service_version"`
	Title          string `json:"title"`
	Url            string `json:"url"`
}

type NEOResponse struct {
	Links struct {
		Next     string `json:"next"`
		Previous string `json:"previous"`
		Self     string `json:"self"`
	} `json:"links"`
	ElementCount int64 `json:"element_count"`
}

// another response object
