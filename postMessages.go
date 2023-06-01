package go253

import (
	"bytes"
	"net/http"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

// client 的标准返回，下游可以根据这个返回来判断是否发送成功
type SendRecord struct {
	Message   *message      `json:"message"`
	Response  *sendResponse `json:"response"`
	Timestamp int64         `json:"timestamp"`
	Error     error         `json:"-"`
}

func (r *SendRecord) reset() {
	r.Message.reset()
	r.Response.reset()
	r.Error = nil
	r.Timestamp = 0
}

type sendResponse struct {
	Code       string    `json:"code"`
	Message    string    `json:"message"`
	MessageID  string    `json:"msgid"`
	ErrorPhone []*string `json:"errorPhone"`
	ErrorNum   *int      `json:"errorNum"`
	SuccessNum *int      `json:"successNum"`
}

func (r *sendResponse) reset() {
	r.Code = ""
	r.Message = ""
	r.MessageID = ""
	r.ErrorPhone = nil
	r.ErrorNum = nil
	r.SuccessNum = nil
}

type normalMessageResponse struct {
	Code  string `json:"code"`
	Error string `json:"error"`
	MsgID string `json:"msgid"`
}

func (r *normalMessageResponse) reset() {
	r.Code = ""
	r.Error = ""
	r.MsgID = ""
}

func (r *normalMessageResponse) toResponse(sr *sendResponse) {
	sr.Code = r.Code
	sr.Message = r.Error
	sr.MessageID = r.MsgID
}

type varMessageResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    struct {
		MessageID  string   `json:"messageId"`
		ErrorPhone []string `json:"errorPhone"`
		ErrorNum   int      `json:"errorNum"`
		SuccessNum int      `json:"successNum"`
	} `json:"data"`
}

func (r *varMessageResponse) reset() {
	r.Code = ""
	r.Message = ""
	r.Data.MessageID = ""
	r.Data.ErrorPhone = nil
	r.Data.ErrorNum = 0
	r.Data.SuccessNum = 0
}

func (r *varMessageResponse) toResponse(sr *sendResponse) {
	sr.Code = r.Code
	sr.Message = r.Message
	sr.MessageID = r.Data.MessageID
	sr.ErrorPhone = make([]*string, len(r.Data.ErrorPhone))
	for i, phone := range r.Data.ErrorPhone {
		sr.ErrorPhone[i] = &phone
	}
	sr.ErrorNum = &r.Data.ErrorNum
	sr.SuccessNum = &r.Data.SuccessNum
}

func (c *client) postMessage(msg *message) {
	defer func() {
		if err := recover(); err != nil {
			// todo: 处理panic
		}
	}()
	// 对象池中取处发送记录*SendRecord, 用于返回给下游
	record := sendRecordPool.Get().(*SendRecord)
	defer sendRecordPool.Put(record)
	defer record.reset()
	// 记录发送时间和发送的消息
	record.Timestamp = time.Now().Unix()
	record.Message = msg
	// http client
	hc := http.DefaultClient
	// endpoint 选择
	var e string
	switch msg.msgType() {
	case SMSTypeNormal:
		switch c.nodeType {
		case NodeShanghai:
			e = NormalEndpointSH
		case NodeSingapore:
			e = NormalEndpointSG
		default:
			// 因为已经在client中做了校验，所以这里不可能出现这种情况
		}
	case SMSTypeVariable:
		switch c.nodeType {
		case NodeShanghai:
			e = VarEndpointSH
		case NodeSingapore:
			e = VarEndpointSG
		default:
			// 同上
		}
	default:
		// 同上
	}
	// 消息序列化
	body, err := jsoniter.Marshal(msg)
	if err != nil {
		record.Error = errors.Wrap(err, "marshal message failed")
		c.recordChan <- *record
		return
	}
	// 创建http请求并发送
	request, err := http.NewRequest(http.MethodPost, e, bytes.NewReader(body))
	if err != nil {
		record.Error = errors.Wrap(err, "create request failed")
		c.recordChan <- *record
		return
	}
	request.Header.Set("Content-Type", "application/json")
	resp, err := hc.Do(request)
	if err != nil {
		record.Error = errors.Wrap(err, "post message failed")
		c.recordChan <- *record
		return
	}
	defer resp.Body.Close()
	// 解析响应
	// 对象池中取出发送响应*sendResponse, 用于返回给下游
	sr := sendResponsePool.Get().(*sendResponse)
	defer sendResponsePool.Put(sr)
	defer sr.reset()
	switch msg.msgType() {
	case SMSTypeNormal:
		res := normalMessageResponsePool.Get().(*normalMessageResponse)
		defer normalMessageResponsePool.Put(res)
		defer res.reset()
		if err := jsoniter.NewDecoder(resp.Body).Decode(&res); err != nil {
			record.Error = errors.Wrap(err, "decode response failed")
			c.recordChan <- *record
			return
		}
		res.toResponse(sr)
		record.Response = sr
		// check response code, 判断是否可以继续发送，并返回错误类型
		ok, err := checkResponseCode(sr.Code)
		if err != nil {
			record.Error = errors.Wrap(err, "check response code failed")
		}
		c.recordChan <- *record
		if !ok {
			defer c.Close()
		}
	case SMSTypeVariable:
		res := varMessageResponsePool.Get().(*varMessageResponse)
		defer varMessageResponsePool.Put(res)
		defer res.reset()
		if err := jsoniter.NewDecoder(resp.Body).Decode(&res); err != nil {
			record.Error = errors.Wrap(err, "decode response failed")
			c.recordChan <- *record
			return
		}
		res.toResponse(sr)
		record.Response = sr
		// check response code
		ok, err := checkResponseCode(sr.Code)
		if err != nil {
			record.Error = errors.Wrap(err, "check response code failed")
		}
		c.recordChan <- *record
		if !ok {
			defer c.Close()
		}
	default:
		// 因为前面已经判断过了，所以这里不会走到
	}
}
