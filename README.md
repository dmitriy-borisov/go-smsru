# Golang-client for the SMS.ru API #

[![Build Status](https://travis-ci.org/dmitriy-borisov/go-smsru.svg?branch=master)](https://travis-ci.org/dmitriy-borisov/go-smsru)
[![GoDoc](https://godoc.org/github.com/dmitriy-borisov/go-smsru?status.svg)](https://godoc.org/github.com/dmitriy-borisov/go-smsru)

Supports:
- sms/send, sms/status, sms/cost
- my/balance, my/limit, my/senders
- stoplist/get, stoplist/add, stoplist/del
- callback/get, callback/add, callback/del

## Installation ##
Install:
```go
go get github.com/dmitriy-borisov/go-smsru
```
Import:
```go
import "github.com/dmitriy-borisov/go-smsru"
```

## Examples ##

```go
package main

import (
    "log"
    "github.com/dmitriy-borisov/go-smsru"
)

const API_ID = "MY_API_ID"

func main() {
    client := sms.NewClient(API_ID)
    
    // Send one message
    msg := sms.NewSms("79250001122", "Sample text")
    
    res, err := client.SmsSend(msg)
    if err != nil {
        log.Panic(err)
    } else {
        log.Printf("Status = %d, Id = %s, Balance = %f", res.Status, res.Ids[0], res.Balance)
    }
    
    // Send multiple messages
    msg := sms.NewSms("79250001122", "Sample text")
    msg2 := sms.NewSms("79251112233", "Sample text")
    multi := sms.NewMulti(msg, msg2)
    
    res, err := client.SmsSend(multi)
    if err != nil {
        log.Panic(err)
    } else {
        log.Printf("Status = %d, Ids = %v, Balance = %f", res.Status, res.Ids, res.Balance)
    }
}
```

## Tests ##

```bash
phone=YOUR_PHONE api_id=YOUR_API_ID go test
 ```
