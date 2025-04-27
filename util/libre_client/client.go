package libreclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
)

var ErrIncorrectUsernameOrPassword = errors.New("incorrect username/password")

// host set for rus region. Default is "https://api.libreview.io"
const host = "https://api.libreview.ru"

// LibreLinkUpManager is a client for the LibreLinkUp API.
//
// SDK client written based on this dump: https://gist.github.com/khskekec/6c13ba01b10d3018d816706a32ae8ab2
type LibreLinkUpManager struct{}

func NewLibreLinkUpManager() *LibreLinkUpManager {
	return &LibreLinkUpManager{}
}

func setDefaultHeaders(r *http.Request) {
	r.Header.Set("cache-control", "no-cache")
	r.Header.Set("connection", "Keep-Alive")
	r.Header.Set("content-type", "application/json")
	r.Header.Set("product", "llu.android")
	r.Header.Set("version", "4.7")
}

func (lm LibreLinkUpManager) makeRequest(method, path string, body io.Reader, token *string) (*http.Response, error) {
	url := host + path
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	setDefaultHeaders(req)
	if token != nil {
		req.Header.Set("authorization", "Bearer "+*token)
	}
	return http.DefaultClient.Do(req)
}

func decodeResponse[Response any](resp *http.Response) (*Response, error) {
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("status code is not 200")
	}
	type response struct {
		Status int      `json:"status"`
		Data   Response `json:"data"`
	}

	var b bytes.Buffer
	resp.Body = io.NopCloser(io.TeeReader(resp.Body, &b))

	var responseData response
	err := json.NewDecoder(resp.Body).Decode(&responseData)
	if err != nil {
		return nil, err
	}
	if responseData.Status != 0 {
		type errorResponseData struct {
			Error struct {
				Message string `json:"message"`
			} `json:"error"`
		}
		var errorData errorResponseData
		err = json.Unmarshal(b.Bytes(), &errorData)
		if err != nil {
			return nil, err
		}
		if errorData.Error.Message == "incorrect username/password" {
			return nil, ErrIncorrectUsernameOrPassword
		} 
		return nil, errors.New(errorData.Error.Message)
	}
	return &responseData.Data, nil
}

func (lm LibreLinkUpManager) Login(email, password string) (*User, error) {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	data := request{Email: email, Password: password}
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	resp, err := lm.makeRequest(http.MethodPost, "/llu/auth/login", bytes.NewReader(body), nil)
	if err != nil {
		return nil, err
	}
	return decodeResponse[User](resp)
}

func (lm LibreLinkUpManager) FetchConnections(token string) ([]LibreConnection, error) {
	resp, err := lm.makeRequest(http.MethodGet, "/llu/connections", nil, &token)
	if err != nil {
		return nil, err
	}
	data, err := decodeResponse[[]LibreConnection](resp)
	if err != nil {
		return nil, err
	}
	return *data, nil
}

func (lm LibreLinkUpManager) FetchData(patientId uuid.UUID, token string) (*GraphData, error) {
	url := fmt.Sprintf("/llu/connections/%s/graph", patientId.String())
	resp, err := lm.makeRequest(http.MethodGet, url, nil, &token)
	if err != nil {
		return nil, err
	}
	return decodeResponse[GraphData](resp)
}
