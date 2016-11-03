package sms

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"
)

// Base API url
const API_URL = "http://sms.ru"

var codeStatus map[int]string = map[int]string{
	-1:  "Not found",
	100: "Success",
	101: "The messege is passed to operator",
	102: "The message sent (in transit)",
	103: "The message was delivered",
	104: "Cannot be delivered: Time of life expired",
	105: "Cannot be delivered: deleted by operator",
	106: "Cannot be delivered: phone failure",
	107: "Cannot be delivered: unknown reason",
	108: "Cannot be delivered: rejected",
	130: "Cannot be delivered: Daily message limit on this number was exceeded",
	131: "Cannot be delivered: Same messages limit on this phone number in a minute was exceeded",
	132: "Cannot be delivered: Same messages limit on this phone number in a day was exceeded",
	200: "Wrong api_id",
	201: "Too low balance",
	202: "Wrong recipient",
	203: "The message has no text",
	204: "Sender name did not approve with administartion",
	205: "The message is too long (more than 8 sms)",
	206: "Daily message limit exceeded",
	207: "On this phone number (or one of them) must not send the messages, or you indicated more than 100 phone numbers",
	208: "Wrong time value",
	209: "You added this phone number (or one of them) in the stop-list",
	210: "You must use a POST, not a GET",
	211: "Method not found",
	212: "Text of message must be in UTF-8",
	220: "The service is not availiable now, try again later",
	230: "Daily message limit on this number was exceeded",
	231: "Same messages limit on this phone number in a minute was exceeded",
	232: "Same messages limit on this phone number in a day was exceeded",
	300: "Wrong token (maybe it was expired or your IP was changed)",
	301: "Wrong password, or user is not exist",
	302: "User was authorized, but account is not activate",
	901: "Wrong Url (should begin with 'http://')",
	902: "Callback is not defined",
}

var error_internal = errors.New("Internal Error")
var error_no_response = errors.New("Something went wrong")

// NewClient creates a new SmsClient instance.
//
// id is your api_id
func NewClient(id string) *SmsClient {
	return NewClientWithHttp(id, &http.Client{})
}

// NewClientWithHttp creates a new SmsClient instance
//
// and allows you to pass a http.Client.
func NewClientWithHttp(id string, client *http.Client) *SmsClient {
	c := &SmsClient{
		ApiId: id,
		Http:  client,
	}

	return c
}

// NewSms creates a new message
//
// to is where to send it (phone number), text is the message text.
func NewSms(to string, text string) *Sms {
	return &Sms{
		To:   to,
		Text: text,
	}
}

// NewMulti creates a one request for multiple messages
func NewMulti(sms ...*Sms) *Sms {
	arr := make(map[string]string)
	for _, o := range sms {
		arr[o.To] = o.Text
	}

	return &Sms{
		Multi: arr,
	}
}

func (c *SmsClient) makeRequest(endpoint string, params url.Values) (Response, []string, error) {
	params.Set("api_id", c.ApiId)
	url := API_URL + endpoint + "?" + params.Encode()

	resp, err := c.Http.Get(url)
	if err != nil {
		return Response{}, nil, err
	}
	defer resp.Body.Close()

	sc := bufio.NewScanner(resp.Body)
	var lines []string
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}

	if err := sc.Err(); err != nil {
		return Response{}, nil, error_internal
	}

	if len(lines) == 0 {
		return Response{}, nil, error_no_response
	}

	status, _ := strconv.Atoi(lines[0])

	if status >= 200 {
		msg := fmt.Sprintf("Code: %d; Status: %s", status, codeStatus[status])
		return Response{}, nil, errors.New(msg)
	}

	res := Response{Status: status}
	return res, lines, nil
}

// SmsSend will send a Sms item to Service
func (c *SmsClient) SmsSend(p *Sms) (Response, error) {
	var params = url.Values{}

	if len(p.Multi) > 0 {
		for to, text := range p.Multi {
			key := fmt.Sprintf("multi[%s]", to)
			params.Add(key, text)
		}
	} else {
		params.Set("to", p.To)
		params.Set("text", p.Text)
	}

	if len(p.From) > 0 {
		params.Set("from", p.From)
	}

	if p.PartnerId > 0 {
		val := strconv.Itoa(p.PartnerId)
		params.Set("partner_id", val)
	}

	if p.Test {
		params.Set("test", "1")
	}

	if p.Time.After(time.Now()) {
		val := strconv.FormatInt(p.Time.Unix(), 10)
		params.Set("time", val)
	}

	if p.Translit {
		params.Set("translit", "1")
	}

	res, lines, err := c.makeRequest("/sms/send", params)
	if err != nil {
		return Response{}, err
	}

	var ids []string
	re := regexp.MustCompile("^balance=")

	for i := 1; i < len(lines); i++ {
		isBalance := re.MatchString(lines[i])

		if isBalance {
			str := re.ReplaceAllString(lines[i], "")
			balance, err := strconv.ParseFloat(str, 32)
			if err != nil {
				return Response{}, error_internal
			}
			res.Balance = float32(balance)
		} else {
			ids = append(ids, lines[i])
		}
	}

	res.Ids = ids
	return res, nil
}

// SmsStatus will get a status of message
func (c *SmsClient) SmsStatus(id string) (Response, error) {
	params := url.Values{}
	params.Set("id", id)

	res, _, err := c.makeRequest("/sms/status", params)
	if err != nil {
		return Response{}, err
	}

	return res, nil
}

// SmsStatus will get a status of message
func (c *SmsClient) SmsCost(p *Sms) (Response, error) {
	var params = url.Values{}
	params.Set("to", p.To)
	params.Set("text", p.Text)
	if p.Translit {
		params.Set("translit", "1")
	}

	res, lines, err := c.makeRequest("/sms/cost", params)
	if err != nil {
		return Response{}, err
	}

	cost, err := strconv.ParseFloat(lines[1], 32)
	if err != nil {
		return Response{}, error_internal
	}

	count, err := strconv.Atoi(lines[2])
	if err != nil {
		return Response{}, error_internal
	}

	res.Cost = float32(cost)
	res.Count = count

	return res, nil
}

// MyBalance checks the balance
func (c *SmsClient) MyBalance() (Response, error) {
	res, lines, err := c.makeRequest("/my/balance", url.Values{})
	if err != nil {
		return Response{}, err
	}

	balance, err := strconv.ParseFloat(lines[1], 32)
	if err != nil {
		return Response{}, error_internal
	}

	res.Balance = float32(balance)
	return res, nil
}

// MyLimit checks the limit
func (c *SmsClient) MyLimit() (Response, error) {
	res, lines, err := c.makeRequest("/my/limit", url.Values{})
	if err != nil {
		return Response{}, err
	}

	limit, err := strconv.Atoi(lines[1])
	if err != nil {
		return Response{}, error_internal
	}

	limitSent, err := strconv.Atoi(lines[2])
	if err != nil {
		return Response{}, error_internal
	}

	res.Limit = limit
	res.LimitSent = limitSent
	return res, nil
}

// MySenders recieves the list of senders
func (c *SmsClient) MySenders() (Response, error) {
	res, lines, err := c.makeRequest("/my/senders", url.Values{})
	if err != nil {
		return Response{}, err
	}

	var senders []string
	for i := 1; i < len(lines); i++ {
		senders = append(senders, lines[i])
	}

	res.Senders = senders
	return res, nil
}

// StoplistGet recieves the stoplist
func (c *SmsClient) StoplistGet() (Response, error) {
	res, lines, err := c.makeRequest("/stoplist/get", url.Values{})
	if err != nil {
		return Response{}, err
	}

	stoplist := make(map[string]string)
	for i := 1; i < len(lines); i++ {
		re := regexp.MustCompile(";")
		str := re.Split(lines[i], 2)

		stoplist[str[0]] = str[1]
	}

	res.Stoplist = stoplist
	return res, nil
}

// StoplistAdd will add the phone number to stoplist
//
// phone is phone number, text is the additional information.
func (c *SmsClient) StoplistAdd(phone, text string) (Response, error) {
	params := url.Values{}
	params.Set("stoplist_phone", phone)
	params.Set("stoplist_text", text)

	res, _, err := c.makeRequest("/stoplist/add", params)
	if err != nil {
		return Response{}, err
	}

	return res, nil
}

// StoplistDel will delete the phone number from stoplist
//
// phone is phone number
func (c *SmsClient) StoplistDel(phone string) (Response, error) {
	params := url.Values{}
	params.Set("stoplist_phone", phone)

	res, _, err := c.makeRequest("/stoplist/del", params)
	if err != nil {
		return Response{}, err
	}

	return res, nil
}

// CallbackGet recieves the callbacks from service
func (c *SmsClient) CallbackGet() (Response, error) {
	res, lines, err := c.makeRequest("/callback/get", url.Values{})
	if err != nil {
		return Response{}, err
	}

	var callbacks []string
	for i := 1; i < len(lines); i++ {
		callbacks = append(callbacks, lines[i])
	}

	res.Callbacks = callbacks
	return res, nil
}

// CallbackAdd will add the callback url to service
//
// cbUrl is your callback url
func (c *SmsClient) CallbackAdd(cbUrl string) (Response, error) {
	params := url.Values{}
	params.Set("url", cbUrl)

	res, lines, err := c.makeRequest("/callback/add", params)
	if err != nil {
		return Response{}, err
	}

	var callbacks []string
	for i := 1; i < len(lines); i++ {
		callbacks = append(callbacks, lines[i])
	}

	res.Callbacks = callbacks
	return res, nil
}

// CallbackDel will delete the callback url from service
//
// cbUrl is your callback url
func (c *SmsClient) CallbackDel(cbUrl string) (Response, error) {
	params := url.Values{}
	params.Set("url", cbUrl)

	res, lines, err := c.makeRequest("/callback/del", params)
	if err != nil {
		return Response{}, err
	}

	var callbacks []string
	for i := 1; i < len(lines); i++ {
		callbacks = append(callbacks, lines[i])
	}

	res.Callbacks = callbacks
	return res, nil
}
