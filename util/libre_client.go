package util

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type LibreTimestamp struct {
	time.Time
}

func (t *LibreTimestamp) UnmarshalJSON(data []byte) error {
	const libreTimestampLayout = "1/2/2006 3:04:05 PM"
	timeString := strings.ReplaceAll(string(data), "\"", "")
	parsedTime, err := time.Parse(libreTimestampLayout, timeString)
	if err != nil {
		return err
	}
	t.Time = parsedTime
	return nil
}

type (
	User struct {
		Id                    uuid.UUID  `json:"id"`
		FirstName             string     `json:"firstName"`
		LastName              string     `json:"lastName"`
		Email                 string     `json:"email"`
		Country               string     `json:"country"`
		UiLanguage            string     `json:"uiLanguage"`
		CommunicationLanguage string     `json:"communicationLanguage"`
		AccountType           string     `json:"accountType"`
		Uom                   string     `json:"uom"`
		DateFormat            string     `json:"dateFormat"`
		TimeFormat            string     `json:"timeFormat"`
		AuthTicket            AuthTicket `json:"authTicket"`
	}

	AuthTicket struct {
		Token    string    `json:"token"`
		Expires  Timestamp `json:"expires"`
		Duration Timestamp `json:"duration"`
	}

	LibreConnection struct {
		Id         uuid.UUID `json:"id"`
		PatientId  uuid.UUID `json:"patientId"`
		Country    string    `json:"country"`
		Status     int       `json:"status"`
		FirstName  string    `json:"firstName"`
		LastName   string    `json:"lastName"`
		TargetLow  int       `json:"targetLow"`
		TargetHigh int       `json:"targetHigh"`
	}

	PatientDevice struct {
		DeviceId uuid.UUID `json:"did"`
		Version  string    `json:"v"`
		Created  Timestamp `json:"created"`
	}

	Connections struct {
		Status int               `json:"status"`
		Data   []LibreConnection `json:"data"`
		Ticket AuthTicket        `json:"ticket"`
	}

	GlucoseMeasurement struct {
		FactoryTimestamp LibreTimestamp `json:"FactoryTimestamp"`
		Timestamp        string         `json:"Timestamp"`
		Type             int            `json:"type"`
		ValueInMgPerDl   int            `json:"ValueInMgPerDl"`
		TrendArrow       int            `json:"TrendArrow"`
		MeasurementColor int            `json:"MeasurementColor"`
		GlucoseUnits     int            `json:"GlucoseUnits"`
		Value            float64        `json:"Value"`
		IsHigh           bool           `json:"isHight"`
		IsLow            bool           `json:"isLow"`
	}

	GraphData struct {
		Connection  LibreConnection      `json:"connection"`
		Mesurements []GlucoseMeasurement `json:"graphData"`
	}

	Graph struct {
		Status int        `json:"status"`
		Data   GraphData  `json:"data"`
		Ticket AuthTicket `json:"ticket"`
	}
)

const host = "https://api.libreview.ru"

type LibreLinkUpManager struct{}

func NewLibreLinkUpManager() *LibreLinkUpManager {
	return &LibreLinkUpManager{}
}

func setDefaultHeaders(r *http.Request) {
	// r.Header.Set("accept-encoding", "gzip")
	r.Header.Set("cache-control", "no-cache")
	r.Header.Set("connection", "Keep-Alive")
	r.Header.Set("content-type", "application/json")
	r.Header.Set("product", "llu.android")
	r.Header.Set("version", "4.7")
	// User-Agent: python-requests/2.31.0
	// r.Header.Set("User-Agent", "python-requests/2.31.0")
}

func (lm *LibreLinkUpManager) makeRequest(method, path string, body io.Reader, token *string) (*http.Response, error) {
	url := host + path
	println("url", url)
	// buf := new(strings.Builder)
	// _, err := io.Copy(buf, body)
	// // check errors
	// println(buf.String())
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

func (lm *LibreLinkUpManager) Login(email, password string) (*User, error) {
	println("here 2")
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	data := request{Email: email, Password: password}
	println("marshalling")
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	println("making request")
	resp, err := lm.makeRequest(http.MethodPost, "/llu/auth/login", bytes.NewReader(body), nil)
	if err != nil {
		return nil, err
	}
	println("here")
	println("response", resp.Status)

	// b, err := httputil.DumpResponse(resp, true)
	// if err != nil {
	// 	return nil, err
	// }
	//
	// println(string(b))

	type response struct {
		Status int  `json:"status"`
		Data   User `json:"data"`
	}

	var repons response
	err = json.NewDecoder(resp.Body).Decode(&repons)
	if err != nil {
		return nil, err
	}
	return &repons.Data, nil
}

func (lm *LibreLinkUpManager) FetchConnections(token string) (*Connections, error) {
	resp, err := lm.makeRequest(http.MethodGet, "/llu/connections", nil, &token)
	if err != nil {
		return nil, err
	}
	var connections Connections
	err = json.NewDecoder(resp.Body).Decode(&connections)
	return &connections, err
}

func (lm *LibreLinkUpManager) FetchData(patientId uuid.UUID, token string) (*Graph, error) {
	url := "/llu/connections/" + patientId.String() + "/graph"
	resp, err := lm.makeRequest(http.MethodGet, url, nil, &token)
	if err != nil {
		return nil, err
	}
	var graph Graph
	err = json.NewDecoder(resp.Body).Decode(&graph)
	return &graph, err
}
