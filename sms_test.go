package sms_test

import (
	"github.com/dm-borisov/go-smsru"
	"os"
	"testing"
)

const (
	CALLBACK_URL = "http://sms.ru/callback"
)

func getPhone() string {
	return os.Getenv("phone")
}

func getClient(t *testing.T) *sms.SmsClient {
	apiId := os.Getenv("api_id")
	return sms.NewClient(apiId)
}

/* Test Sms
---------------------------------------------*/
func TestSmsSend(t *testing.T) {
	c := getClient(t)

	msg := sms.NewSms(getPhone(), "Sample")
	msg.Test = true

	_, err := c.SmsSend(msg)

	if err != nil {
		t.Fail()
	}
}

func TestSmsMultiSend(t *testing.T) {
	c := getClient(t)

	msg := sms.NewSms(getPhone(), "Sample")
	multi := sms.NewMulti(msg)
	multi.Test = true

	_, err := c.SmsSend(multi)

	if err != nil {
		t.Fail()
	}
}

func TestSmsStatus(t *testing.T) {
	c := getClient(t)
	id := "201600-1000000"

	_, err := c.SmsStatus(id)

	if err != nil {
		t.Fail()
	}
}

func TestSmsCost(t *testing.T) {
	c := getClient(t)
	msg := sms.NewSms(getPhone(), "Sample")

	_, err := c.SmsCost(msg)

	if err != nil {
		t.Fail()
	}
}

/* Test My
---------------------------------------------*/
func TestMyBalance(t *testing.T) {
	c := getClient(t)

	_, err := c.MyBalance()
	if err != nil {
		t.Fail()
	}
}

func TestMyLimit(t *testing.T) {
	c := getClient(t)

	_, err := c.MyLimit()
	if err != nil {
		t.Fail()
	}
}

func TestMySenders(t *testing.T) {
	c := getClient(t)

	_, err := c.MySenders()
	if err != nil {
		t.Fail()
	}
}

/* Test Stoplist
---------------------------------------------*/
func TestStoplistAdd(t *testing.T) {
	c := getClient(t)

	_, err := c.StoplistAdd(getPhone(), "TestAdd")
	if err != nil {
		t.Fail()
	}
}

func TestStoplistGet(t *testing.T) {
	c := getClient(t)

	_, err := c.StoplistGet()
	if err != nil {
		t.Fail()
	}
}

func TestStoplistDel(t *testing.T) {
	c := getClient(t)

	_, err := c.StoplistDel(getPhone())
	if err != nil {
		t.Fail()
	}
}

/* Test Callback
---------------------------------------------*/
func TestCallbackAdd(t *testing.T) {
	c := getClient(t)

	_, err := c.CallbackAdd(CALLBACK_URL)
	if err != nil {
		t.Fail()
	}
}

func TestCallbackGet(t *testing.T) {
	c := getClient(t)

	_, err := c.CallbackGet()
	if err != nil {
		t.Fail()
	}
}

func TestCallbackDel(t *testing.T) {
	c := getClient(t)

	_, err := c.CallbackDel(CALLBACK_URL)
	if err != nil {
		t.Fail()
	}
}
