package libreclient

import (
	"strconv"

	"github.com/tikhonp/medsenger-freestylelibre-bot/util"
	"github.com/google/uuid"
)

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
		Token    string         `json:"token"`
		Expires  util.Timestamp `json:"expires"`
		Duration util.Timestamp `json:"duration"`
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
		DeviceId uuid.UUID      `json:"did"`
		Version  string         `json:"v"`
		Created  util.Timestamp `json:"created"`
	}

	GlucoseMeasurement struct {
		FactoryTimestamp LibreTimeFormat `json:"FactoryTimestamp"`
		Timestamp        string          `json:"Timestamp"`
		Type             int             `json:"type"`
		ValueInMgPerDl   int             `json:"ValueInMgPerDl"`
		TrendArrow       int             `json:"TrendArrow"`
		MeasurementColor int             `json:"MeasurementColor"`
		GlucoseUnits     int             `json:"GlucoseUnits"`
		Value            float64         `json:"Value"`
		IsHigh           bool            `json:"isHight"`
		IsLow            bool            `json:"isLow"`
	}

	GraphData struct {
		Connection  LibreConnection      `json:"connection"`
		Mesurements []GlucoseMeasurement `json:"graphData"`
	}
)

func (gm *GlucoseMeasurement) ValueAsString() string {
	return strconv.FormatFloat(gm.Value, 'f', 2, 64)
}
