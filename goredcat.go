package goredcat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	prefixURL = "/api/v1"
	loginURL  = "/login"
	//salesByURL  = "/reports/loyalty/membersalesby"
	salesByURL  = "/reports/kpi/salesby"
	locationURL = "/reports/data/locations"
	authHeader  = "X-Redcat-Authtoken"
)

var (
	defaultSendTimeout = time.Second * 30
)

// Redcat The main struct of this package
type Redcat struct {
	BaseURL     string
	LoginURL    string
	Username    string
	Password    string
	accessToken string
}

// NewClient will create a Redcat client with default values
func NewClient(baseURL string) *Redcat {
	return &Redcat{
		BaseURL:  baseURL + prefixURL,
		LoginURL: baseURL + prefixURL + loginURL,
	}
}

// AccessToken will get a new access token
func (v *Redcat) AccessToken(username string, password string) (bool, error) {
	u, _ := url.ParseRequestURI(v.LoginURL)
	urlStr := fmt.Sprintf("%v", u)

	request := LoginRequest{
		Username: username,
		Password: password,
		AuthType: "U",
	}

	fmt.Println("Connecting to redcat", urlStr, request)

	body, err := json.Marshal(request)

	if err != nil {
		fmt.Println("Request marshal err", err)
		return false, err
	}

	client := &http.Client{}

	r, err := http.NewRequest("POST", urlStr, bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Redcat POST err ", err)
		return false, err
	}

	r.Header = http.Header(make(map[string][]string))
	r.Header.Set("Content-Type", "application/json")

	res, err := client.Do(r)
	if err != nil {
		fmt.Println("Redcat client.DO err ", err)
		return false, err
	}

	if res.StatusCode == 200 {

		rawResBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println("Redcat ReadAll err ", err)
			return false, err
		}

		var resp LoginResponse

		err = json.Unmarshal(rawResBody, &resp)

		if err != nil {
			fmt.Println("Redcat MARSHAL err ", err)
			return false, err
		}

		fmt.Println("Connection Response", resp, string(rawResBody))

		v.accessToken = resp.Token

		return resp.Success, nil
	}

	fmt.Println("Redcat Status Code ", res.StatusCode)

	return false, fmt.Errorf("Failed redcat login %s", res.Status)
}

// RequestSalesReport will request a report from Redcat
func (v *Redcat) RequestSalesReport(request ReportRequest) (*ReportResult, error) {
	urlStr := v.BaseURL + salesByURL

	fmt.Println("Connecting to redcat", urlStr)

	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	fmt.Println("Sending body", string(body))

	client := &http.Client{}

	r, err := http.NewRequest("POST", urlStr, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	r.Header = http.Header(make(map[string][]string))
	r.Header.Set(authHeader, v.accessToken)
	r.Header.Set("Content-Type", "application/json")

	res, err := client.Do(r)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == 200 {

		rawResBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		//fmt.Println("Got response", string(rawResBody))

		var resp ReportResult

		err = json.Unmarshal(rawResBody, &resp)

		if err != nil {
			return nil, err
		}
		return &resp, nil
	}

	return nil, fmt.Errorf("Failed to get Redcat Sales Report %s", res.Status)
}

// RequestLocations will request a report from Redcat
func (v *Redcat) RequestLocations() (string, error) {
	urlStr := v.BaseURL + locationURL

	fmt.Println("Connecting to redcat", urlStr)

	client := &http.Client{}

	r, err := http.NewRequest("POST", urlStr, nil)
	if err != nil {
		return "", err
	}

	r.Header = http.Header(make(map[string][]string))
	r.Header.Set(authHeader, v.accessToken)
	r.Header.Set("Content-Type", "application/json")

	res, err := client.Do(r)
	if err != nil {
		return "", err
	}

	if res.StatusCode == 200 {

		rawResBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return "", err
		}
		return string(rawResBody), nil

		/*var resp ReportResult

		err = json.Unmarshal(rawResBody, &resp)

		if err != nil {
			return nil, err
		}
		return &resp, nil*/
	}

	return "", fmt.Errorf("Failed to get Redcat Sales Report %s", res.Status)
}
