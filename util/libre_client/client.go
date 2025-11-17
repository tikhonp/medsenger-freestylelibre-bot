// Package libreclient provides a client for the LibreLinkUp API.
package libreclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

var (
	ErrIncorrectUsernameOrPassword = errors.New("incorrect username/password")
	ErrInvalidAuthSession          = errors.New("invalid auth session")
)

// host set for rus region. Default is "https://api.libreview.io"
const host = "https://api.libreview.ru"

// LibreLinkUpManager is a client for the LibreLinkUp API.
//
// SDK client written inspired by this dump: https://gist.github.com/khskekec/6c13ba01b10d3018d816706a32ae8ab2
type LibreLinkUpManager struct{}

func NewLibreLinkUpManager() *LibreLinkUpManager {
	return &LibreLinkUpManager{}
}

func setDefaultHeaders(r *http.Request) {
	r.Header.Set("Accept", "*/*")
	r.Header.Set("Accept-Encoding", "gzip, deflate, br")
	r.Header.Set("Connection", "keep-alive")
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("product", "llu.ios")
	r.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU OS 16_2 like Mac OS X) AppleWebKit/536.26 (KHTML, like Gecko) Version/16.2 Mobile/10A5355d Safari/8536.25")
	r.Header.Set("version", "4.16.0")
}

func (lm LibreLinkUpManager) makeRequest(method, path string, body io.Reader, token, accountID *string) (*http.Response, error) {
	url := host + path
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	setDefaultHeaders(req)
	if token != nil {
		req.Header.Set("authorization", "Bearer "+*token)
	}
	if accountID != nil {
		// Account-Id is the SHA-256 hex digest of the user.id UUID
		//  exactly as returned in the login response (dashes kept, lowercase)
		req.Header.Set("Account-Id", *accountID)
	}
	return http.DefaultClient.Do(req)
}

func decodeResponse[Response any](resp *http.Response) (*Response, error) {
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("status code: %d", resp.StatusCode)
		} else {
			respStr := string(body)
			if strings.Contains(respStr, "invalid auth session") {
				return nil, ErrInvalidAuthSession
			}
			return nil, fmt.Errorf("status code: %d response string: %s", resp.StatusCode, respStr)
		}
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

func (lm LibreLinkUpManager) Login(email, password string) (*LoginRespose, error) {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	data := request{Email: email, Password: password}
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	resp, err := lm.makeRequest(http.MethodPost, "/llu/auth/login", bytes.NewReader(body), nil, nil)
	if err != nil {
		return nil, err
	}
	return decodeResponse[LoginRespose](resp)
}

func (lm LibreLinkUpManager) FetchConnections(token, accountID string) ([]LibreConnection, error) {
	resp, err := lm.makeRequest(http.MethodGet, "/llu/connections", nil, &token, &accountID)
	if err != nil {
		return nil, err
	}
	data, err := decodeResponse[[]LibreConnection](resp)
	if err != nil {
		return nil, err
	}
	return *data, nil
}

func (lm LibreLinkUpManager) FetchData(patientID uuid.UUID, token, accountID string) (*GraphData, error) {
	url := fmt.Sprintf("/llu/connections/%s/graph", patientID.String())
	resp, err := lm.makeRequest(http.MethodGet, url, nil, &token, &accountID)
	if err != nil {
		return nil, fmt.Errorf("fetch data request err: %w", err)
	}
	graphData, err := decodeResponse[GraphData](resp)
	if err != nil {
		return nil, fmt.Errorf("fetch data decode err: %w", err)
	}
	return graphData, nil
}
