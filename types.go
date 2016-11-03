package sms

import (
	"net/http"
	"time"
)

// Basic
type SmsClient struct {
	ApiId string       `json:"api_id"`
	Http  *http.Client `json:"-"`
	Debug bool         `json:"-"`
}

type Response struct {
	Status    int               `json:"status"`
	Ids       []string          `json:"id"`
	Cost      float32           `json:"cost"`
	Count     int               `json:"count"`
	Balance   float32           `json:"balance"`
	Limit     int               `json:"limit"`
	LimitSent int               `json:"limit_sent"`
	Senders   []string          `json:"senders"`
	Stoplist  map[string]string `json:"stoplist"`
	Callbacks []string          `json:"callbacks"`
}

type Sms struct {
	To        string            `json:"to"`
	Text      string            `json:"text"`
	Translit  bool              `json:"translit"`
	Multi     map[string]string `json:"multi"`
	From      string            `json:"from"`
	Time      time.Time         `json:"time"`
	Test      bool              `json:"test"`
	PartnerId int               `json:"partner_id"`
}
