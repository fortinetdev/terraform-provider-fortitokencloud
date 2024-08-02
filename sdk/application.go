package ftc_client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var appApiPath = "/api/v1/application"

// GetApplications - Returns all applications
func (c *Client) GetApplications() (*Applications, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", c.HostURL, appApiPath), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	apps := Applications{}
	err = json.Unmarshal(body, &apps.Apps)
	if err != nil {
		return nil, err
	}

	return &apps, nil
}

// GetApplication - Returns specific application
func (c *Client) GetApplication(AppId string) (*Application, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s/%s", c.HostURL, appApiPath, AppId), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	app := Application{}

	err = json.Unmarshal(body, &app)
	if err != nil {
		return nil, err
	}

	return &app, nil
}

// CreateApplication - Create a new application
func (c *Client) CreateApplication(appData interface{}) (*Application, error) {
	rb, err := json.Marshal(appData)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", c.HostURL, appApiPath), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	app := Application{}
	err = json.Unmarshal(body, &app)
	if err != nil {
		return nil, err
	}

	return &app, nil
}

// UpdateApplication - Updates an application
func (c *Client) UpdateApplication(appId string, appData interface{}) (*Application, error) {
	rb, err := json.Marshal(appData)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s%s/%s", c.HostURL, appApiPath, appId), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	app := Application{}
	err = json.Unmarshal(body, &app)
	if err != nil {
		return nil, err
	}

	return &app, nil
}

// DeleteApplication - Deletes an application
func (c *Client) DeleteApplication(appId string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s%s/%s", c.HostURL, appApiPath, appId), nil)
	if err != nil {
		return err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return err
	}

	if string(body) != "" {
		return errors.New(string(body))
	}

	return nil
}

// UpdateApplicationUserSource - Updates an application's user sources
func (c *Client) UpdateApplicationUserSource(appId string, userSourceList map[string][]string) (*[]AppUserMapping, error) {
	// when array is empty, json marshal will convert it to null, need to use make
	if len(userSourceList["user_source_ids"]) == 0 {
		userSourceList["user_source_ids"] = make([]string, 0)
	}
	rb, err := json.Marshal(userSourceList)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s%s/%s/user_source", c.HostURL, appApiPath, appId), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var appusermapping []AppUserMapping

	err = json.Unmarshal(body, &appusermapping)
	if err != nil {
		return nil, err
	}

	return &appusermapping, nil
}
