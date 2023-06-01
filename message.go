package go253

import (
	"strings"

	"github.com/pkg/errors"
)

func NewMessage(ops ...msgOption) (*message, error) {
	m := msgPool.Get().(*message)
	m.reset()
	for _, op := range ops {
		op(m)
	}
	if err := m.validate(); err != nil {
		return nil, errors.Wrap(err, "message validation failed")
	}
	return m, nil
}

// WithMessage sets the moble  message content. and params is optional.
func WithMessage(mobile, msg string, params ...string) msgOption {
	return func(m *message) {
		if len(params) > 0 {
			m.Msg = msg
			ts := make([]string, 0, len(params)+1)
			ts = append(ts, mobile)
			ts = append(ts, params...)
			s := strings.Join(ts, ";")
			m.Params = &s
		} else {
			m.Msg = msg
			m.Mobile = &mobile
		}
	}
}

// WithUid sets the uid for the message.
func WithUid(uid string) msgOption {
	return func(m *message) {
		m.Uid = &uid
	}
}

// WithSenderID sets the sender id for the message.
func WithSenderID(id string) msgOption {
	return func(m *message) {
		m.SenderID = &id
	}
}

type msgOption func(*message)

// account 和 password 在发送前由 client 设置
type message struct {
	Account  string  `json:"account" validate:"required"`
	Password string  `json:"password" validate:"required"`
	Msg      string  `json:"msg" validate:"required"` // message content
	Mobile   *string `json:"mobile"`                  // message receiver mobile
	Params   *string `json:"params"`                  // message params
	Uid      *string `json:"uid"`                     // message uid
	SenderID *string `json:"senderId"`                // message sender id
}

func (m *message) reset() {
	m.Account = ""
	m.Password = ""
	m.Msg = ""
	m.Params = nil
	m.Uid = nil
	m.SenderID = nil
}

func (m *message) msgType() int {
	// 如果mobile不为空，则为普通短信
	if m.Mobile != nil {
		return SMSTypeNormal
	}
	return SMSTypeVariable
}

func (m *message) validate() error {
	if m.Msg == "" {
		return errors.New("message content is empty")
	}
	if m.Mobile != nil && m.Params != nil {
		return errors.New("mobile and params cannot be set at the same time")
	}
	return nil
}
