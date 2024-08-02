package ftc_client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var usApiPath = "/api/v1/usersource"

// GetUserSources - Returns all user sources
func (c *Client) GetUserSources() (*UserSources, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", c.HostURL, usApiPath), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	usersources := UserSources{}
	err = json.Unmarshal(body, &usersources.UserSources)
	if err != nil {
		return nil, err
	}

	return &usersources, nil
}

// GetApplication - Returns specific application
func (c *Client) GetUserSource(UserSourceId string) (*UserSource, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s/%s", c.HostURL, usApiPath, UserSourceId), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	usersource := UserSource{}

	err = json.Unmarshal(body, &usersource)
	if err != nil {
		return nil, err
	}

	return &usersource, nil
}

// CreateApplication - Create a new application
func (c *Client) CreateUserSource(usData interface{}) (*UserSource, error) {
	rb, err := json.Marshal(usData)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", c.HostURL, usApiPath), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	usersource := UserSource{}
	err = json.Unmarshal(body, &usersource)
	if err != nil {
		return nil, err
	}

	return &usersource, nil
}

// UpdateApplication - Updates an application
func (c *Client) UpdateUserSource(UserSourceId string, userSourceData interface{}) (*UserSource, error) {
	rb, err := json.Marshal(userSourceData)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s%s/%s", c.HostURL, usApiPath, UserSourceId), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	usersource := UserSource{}
	err = json.Unmarshal(body, &usersource)
	if err != nil {
		return nil, err
	}

	return &usersource, nil
}

// DeleteApplication - Deletes an application
func (c *Client) DeleteUserSource(UserSourceId string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s%s/%s", c.HostURL, usApiPath, UserSourceId), nil)
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
func (c *Client) UpdateUserSourceDomains(UserSourceId string, domainList map[string][]string) (*[]UserSourceDomainMapping, error) {
	// when array is empty, json marshal will convert it to null, need to use make
	if len(domainList["domain_ids"]) == 0 {
		domainList["domain_ids"] = make([]string, 0)
	}
	rb, err := json.Marshal(domainList)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s%s/%s/domain", c.HostURL, usApiPath, UserSourceId), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var usersourcedomainmapping []UserSourceDomainMapping

	err = json.Unmarshal(body, &usersourcedomainmapping)
	if err != nil {
		return nil, err
	}

	return &usersourcedomainmapping, nil
}

func (c *Client) GetDomain(DomainId string) (*Domain, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s/domain/%s", c.HostURL, usApiPath, DomainId), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	domain := Domain{}

	err = json.Unmarshal(body, &domain)
	if err != nil {
		return nil, err
	}

	return &domain, nil
}

func (c *Client) CreateDomain(domainData interface{}) (*Domain, error) {
	rb, err := json.Marshal(domainData)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s/domain", c.HostURL, usApiPath), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	domain := Domain{}
	err = json.Unmarshal(body, &domain)
	if err != nil {
		return nil, err
	}

	return &domain, nil
}

func (c *Client) DeleteDomain(DomainId string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s%s/domain/%s", c.HostURL, usApiPath, DomainId), nil)
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

func (c *Client) UpdateDomain(DomainId string, domainData interface{}) (*Domain, error) {
	reqBody, err := json.Marshal(domainData)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s%s/domain/%s", c.HostURL, usApiPath, DomainId), strings.NewReader(string(reqBody)))
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	domain := Domain{}
	err = json.Unmarshal(body, &domain)
	if err != nil {
		return nil, err
	}
	return &domain, nil
}
