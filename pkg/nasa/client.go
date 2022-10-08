package nasa

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/igilgyrg/betera-test/pkg/logging"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	url = "https://api.nasa.gov/"
)

type NASAClient interface {
	APOD(date time.Time) (*APODResponse, error)
}

type nasaClient struct {
	token  string
	client *NASAHttpClient
}

func NewClient(token string) NASAClient {
	return &nasaClient{token: token, client: NewNasaHttpClient()}
}

func (n nasaClient) APOD(date time.Time) (*APODResponse, error) {
	apodUrl := fmt.Sprintf("%splanetary/apod", url)

	req, err := http.NewRequest("GET", apodUrl, nil)
	if err != nil {
		logging.Log().Error("error of getting apod image")
		return nil, err
	}

	q := req.URL.Query()
	q.Add("api_key", n.token)
	q.Add("date", date.Format("2006-01-02"))

	req.URL.RawQuery = q.Encode()

	resp, err := n.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("NASA API client return error %s", resp.Status))
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result *APODResponse
	err = json.Unmarshal(responseBody, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
