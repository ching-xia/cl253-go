package go253

import (
	"bytes"
	"fmt"
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

func (c *client) Balance() (balance float64, err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("panic:", r)
		}
	}()
	hc := http.DefaultClient
	type balanceRequest struct {
		Account  string `json:"account"`
		Password string `json:"password"`
	}
	body, err := jsoniter.Marshal(&balanceRequest{
		Account:  c.account,
		Password: c.password,
	})
	if err != nil {
		return 0, errors.Wrap(err, "marshal balance request failed")
	}
	request, err := http.NewRequest(http.MethodPost, BalanceEndpoint, bytes.NewReader(body))
	if err != nil {
		return
	}
	request.Header.Set("Content-Type", "application/json")
	resp, err := hc.Do(request)
	if err != nil {
		return 0, errors.Wrap(err, "get balance failed")
	}
	defer resp.Body.Close()
	type balanceResponse struct {
		Code    int     `json:"code"`
		Error   string  `json:"error"`
		Balance float64 `json:"balance"`
	}
	var br balanceResponse
	if err := jsoniter.NewDecoder(resp.Body).Decode(&br); err != nil {
		return 0, errors.Wrap(err, "decode balance response failed")
	}
	if br.Code != 0 {
		return 0, errors.Errorf("get balance failed: %s", br.Error)
	}
	return br.Balance, nil
}
