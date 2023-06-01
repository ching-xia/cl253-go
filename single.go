package go253

import (
	"bytes"
	"net/http"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

func (c *client) SingleMessage(msg *message) SendRecord {
	record := sendRecordPool.Get().(*SendRecord)
	defer sendRecordPool.Put(record)
	defer record.reset()
	record.Message = msg
	sr := sendResponsePool.Get().(*sendResponse)
	defer sendResponsePool.Put(sr)
	defer sr.reset()
	hc := http.DefaultClient

	var e string
	// endpoint
	switch msg.msgType() {
	case SMSTypeNormal:
		switch c.nodeType {
		case NodeShanghai:
			e = NormalEndpointSH
		case NodeSingapore:
			e = NormalEndpointSG
		default:
			record.Error = errors.Errorf("unknown node number: %d", c.nodeType)
			return *record
		}
	case SMSTypeVariable:
		switch c.nodeType {
		case NodeShanghai:
			e = VarEndpointSH
		case NodeSingapore:
			e = VarEndpointSG
		default:
			record.Error = errors.Errorf("unknown node number: %d", c.nodeType)
			return *record
		}
	default:
		record.Error = errors.Errorf("unknown message type: %d", msg.msgType())
		return *record
	}
	body, err := jsoniter.Marshal(msg)
	if err != nil {
		record.Error = errors.Wrap(err, "marshal message failed")
		return *record
	}
	request, err := http.NewRequest(http.MethodPost, e, bytes.NewReader(body))
	if err != nil {
		record.Error = errors.Wrap(err, "create request failed")
		return *record
	}
	request.Header.Set("Content-Type", "application/json")
	resp, err := hc.Do(request)
	if err != nil {
		record.Error = errors.Wrap(err, "post message failed")
		return *record
	}
	defer resp.Body.Close()
	switch msg.msgType() {
	case SMSTypeNormal:
		res := normalMessageResponsePool.Get().(*normalMessageResponse)
		defer normalMessageResponsePool.Put(res)
		defer res.reset()
		if err := jsoniter.NewDecoder(resp.Body).Decode(&res); err != nil {
			record.Error = errors.Wrap(err, "decode response failed")
			return *record
		}
		res.toResponse(sr)
		record.Response = sr
		record.Timestamp = time.Now().Unix()
		// check response code
		_, err := checkResponseCode(sr.Code)
		if err != nil {
			record.Error = errors.Wrap(err, "check response code failed")
		}
		return *record
	case SMSTypeVariable:
		res := varMessageResponsePool.Get().(*varMessageResponse)
		defer varMessageResponsePool.Put(res)
		defer res.reset()
		if err := jsoniter.NewDecoder(resp.Body).Decode(&res); err != nil {
			record.Error = errors.Wrap(err, "decode response failed")
			return *record
		}
		res.toResponse(sr)
		record.Response = sr
		record.Timestamp = time.Now().Unix()
		// check response code
		_, err := checkResponseCode(sr.Code)
		if err != nil {
			record.Error = errors.Wrap(err, "check response code failed")
		}
		return *record
	default:
		record.Error = errors.Errorf("unknown message type: %d", msg.msgType())
		return *record
	}
}
