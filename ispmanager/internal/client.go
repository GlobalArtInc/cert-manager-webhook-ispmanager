package internal

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/google/go-querystring/query"
	"net/http"
	"strings"
)

type CreateTxtRecord struct {
	AuthInfo string `url:"authinfo"`
	Out      string `url:"out"`
	Func     string `url:"func"`
	Plid     string `url:"plid"`
	Rtype    string `url:"rtype"`
	Name     string `url:"name"`
	Sok      string `url:"sok"`
	Value    string `url:"value"`
}

type DeleteTxtRecord struct {
	AuthInfo string `url:"authinfo"`
	Out      string `url:"out"`
	Func     string `url:"func"`
	Plid     string `url:"plid"`
	Elid     string `url:"elid"`
	Elname   string `url:"elname"`
}

type APIError struct {
	Doc struct {
		Error struct {
			Type   string `json:"$type"`
			Object string `json:"$object"`
			Lang   string `json:"$lang"`
			Detail struct {
				Text string `json:"$"`
			} `json:"detail"`
			Message struct {
				Text string `json:"$"`
			} `json:"msg"`
		} `json:"error"`
	} `json:"doc"`
}

type Client struct {
	PanelUrl   string
	Username   string
	Password   string
	HttpClient *http.Client
}

func (c *Client) newRequest(method string, body interface{}) (*http.Request, error) {
	values, _ := query.Values(body)
	req, err := http.NewRequest(method, c.PanelUrl, strings.NewReader(values.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create new http request with error: %v", err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return req, nil
}

func checkResponse(res *http.Response) error {
	if res.Body == nil {
		return fmt.Errorf("request failed with status code %v and empty body", res.StatusCode)
	}

	//fmt.Printf("%s", res.Body)
	decoder := json.NewDecoder(res.Body)
	//
	var apiError APIError
	err := decoder.Decode(&apiError)
	if err != nil {
		return fmt.Errorf("failed to decode: %s", err)
	}
	if apiError != (APIError{}) {
		//fmt.Printf("dev: %s", apiError.Doc.Error.Message.Text)
		return fmt.Errorf("ISPManager Error: %s", apiError.Doc.Error.Message.Text)
	}

	return nil
}

func (c *Client) do(req *http.Request, to interface{}) (*http.Response, error) {
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed with error: %v", err)
	}
	err = checkResponse(resp)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (c Client) CreateTXT(body CreateTxtRecord) (*CreateTxtRecord, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	req, err := c.newRequest(http.MethodPost, body)
	if err != nil {
		return nil, err
	}
	record := &CreateTxtRecord{}
	_, err = c.do(req, &record)
	if err != nil {
		return nil, err
	}

	return record, nil
}

func (c Client) DeleteTXT(body DeleteTxtRecord) (*DeleteTxtRecord, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	req, err := c.newRequest(http.MethodPost, body)
	if err != nil {
		return nil, err
	}
	record := &DeleteTxtRecord{}
	_, err = c.do(req, &record)
	if err != nil {
		return nil, err
	}

	return record, nil
}
