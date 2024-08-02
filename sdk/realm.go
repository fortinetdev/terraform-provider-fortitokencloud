package ftc_client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var realmApiPath = "/api/v1/realm"

// GetRealmByName - Returns realm with name
func (c *Client) GetRealmByName(RealmName string) (*Realm, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s?name=%s", c.HostURL, realmApiPath, RealmName), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	realm := []Realm{}

	err = json.Unmarshal(body, &realm)
	if err != nil {
		return nil, err
	}

	if len(realm) != 1 {
		return nil, fmt.Errorf(fmt.Sprintf("expected a single body, got %d", len(body)))
	}

	return &realm[0], nil
}

// CreateRealm - Create a new realm
func (c *Client) CreatRealm(realmData interface{}) (*Realm, error) {
	rb, err := json.Marshal(realmData)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", c.HostURL, realmApiPath), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	realm := Realm{}
	err = json.Unmarshal(body, &realm)
	if err != nil {
		return nil, err
	}

	return &realm, nil
}

// UpdateRealm - Updates an realm
func (c *Client) UpdateRealm(realmId string, realmData interface{}) (*Realm, error) {
	rb, err := json.Marshal(realmData)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s%s/%s", c.HostURL, appApiPath, realmId), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	realm := Realm{}
	err = json.Unmarshal(body, &realm)
	if err != nil {
		return nil, err
	}

	return &realm, nil
}

// DeleteRealm - Deletes a realm
func (c *Client) DeleteRealm(realmId string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s%s%s", c.HostURL, realmApiPath, realmId), nil)
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
