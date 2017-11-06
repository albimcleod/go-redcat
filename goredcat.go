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
	prefixURL  = "/api/v1"
	loginURL   = "/login"
	salesByURL = "reports/kpi/salesby"
	authHeader = "X-Redcat-Authtoken"
)

var (
	defaultSendTimeout = time.Second * 30
)

// Redcat The main struct of this package
type Redcat struct {
	BaseURL     string
	accessToken string
}

// NewClient will create a Redcat client with default values
func NewClient(baseURL string) *Redcat {
	return &Redcat{
		BaseURL: baseURL + prefixURL,
	}
}

// AccessToken will get a new access token
func (v *Redcat) AccessToken(username string, password string) (bool, error) {
	u, _ := url.ParseRequestURI(v.BaseURL)
	u.Path = loginURL
	urlStr := fmt.Sprintf("%v", u)

	request := LoginRequest{
		Username: username,
		Password: password,
		AuthType: "M",
	}

	body, err := json.Marshal(request)
	if err != nil {
		return false, err
	}

	client := &http.Client{}
	//	client.CheckRedirect = checkRedirectFunc

	r, err := http.NewRequest("POST", urlStr, bytes.NewBuffer(body))
	if err != nil {
		return false, err
	}

	res, err := client.Do(r)
	if err != nil {
		return false, err
	}

	if res.StatusCode == 200 {

		rawResBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return false, err
		}

		var resp LoginResponse

		err = json.Unmarshal(rawResBody, &resp)

		if err != nil {
			return false, err
		}
		return resp.Success, nil
	}

	return false, fmt.Errorf("Failed redcat login %s", res.Status)
}

// RequestSalesReport will request a report from Redcat
func (v *Redcat) RequestSalesReport(request ReportRequest) (*ReportResult, error) {
	u, _ := url.ParseRequestURI(v.BaseURL)
	u.Path = salesByURL
	urlStr := fmt.Sprintf("%v", u)

	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	r, err := http.NewRequest("POST", urlStr, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	r.Header = http.Header(make(map[string][]string))
	r.Header.Set(authHeader, v.accessToken)

	res, err := client.Do(r)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == 200 {

		rawResBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		var resp ReportResult

		err = json.Unmarshal(rawResBody, &resp)

		if err != nil {
			return nil, err
		}
		return &resp, nil
	}

	return nil, fmt.Errorf("Failed to get Redcat Sales Report %s", res.Status)
}
