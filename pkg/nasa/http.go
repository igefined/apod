package nasa

import (
	"net/http"
	"time"
)

const httpClientTmt = time.Second * 5

type NASAHttpClient struct {
	http.Client
}

func NewNasaHttpClient() *NASAHttpClient {
	return &NASAHttpClient{http.Client{Timeout: httpClientTmt}}
}

// TODO refactoring
// NASAResponse GET I have tried to use generics, seems bad...
//type NASAResponse[T Response] struct {
//	response T
//}
//
//func GET[T NASAResponse](url string, params map[string]string) (*NASAResponse[T], error) {
//	req, err := http.NewRequest("GET", url, nil)
//	if err != nil {
//		logging.Log().Error("error of getting apod image")
//		return &NASAResponse[T]{}, err
//	}
//
//	q := req.URL.Query()
//	for k, v := range params {
//		q.Add(k, v)
//	}
//
//	req.URL.RawQuery = q.Encode()
//
//	client := http.Client{}
//
//	resp, err := client.Do(req)
//	if err != nil {
//		return &NASAResponse[T]{}, err
//	}
//
//	if resp.StatusCode != http.StatusOK {
//		return nil, errors.New(fmt.Sprintf("NASA API client return error %s", resp.Status))
//	}
//
//	responseBody, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		return &NASAResponse[T]{}, err
//	}
//
//	var result T
//	err = json.Unmarshal(responseBody, &result)
//	if err != nil {
//		return &NASAResponse[T]{}, err
//	}
//
//	return &NASAResponse[T]{response: result}, nil
//}
