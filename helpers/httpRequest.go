package helper

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

func GetAPI(endpoint string, params map[string]string) (string, error) {
	// Create a URL object and add parameters
	baseURL, err := url.Parse(endpoint)
	if err != nil {
		return "", err
	}

	// Add query parameters to the URL
	queryParams := url.Values{}
	for key, value := range params {
		queryParams.Add(key, value)
	}
	baseURL.RawQuery = queryParams.Encode()

	// Make the GET request
	resp, err := http.Get(baseURL.String())
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read and return the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func PostAPI(endpoint string, jsonData map[string]interface{}) (string, error) {
	// Marshal the map into JSON format
	jsonValue, err := json.Marshal(jsonData)
	if err != nil {
		return "", err
	}

	// Create a new POST request with the JSON data
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonValue))
	if err != nil {
		return "", err
	}

	// Set the appropriate headers
	req.Header.Set("Content-Type", "application/json")

	// Make the POST request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read and return the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
