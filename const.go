package go253

import "net/url"

const (
	NodeShanghai = iota + 1
	NodeSingapore
)
const (
	SMSTypeNormal = iota + 1
	SMSTypeVariable
)

const (
	NormalEndpointSH = "http://intapi.253.com/send/json"
	VarEndpointSH    = "http://intapi.253.com/send/var/json"
	NormalEndpointSG = "http://intapi.sgap.253.com/send/json"
	VarEndpointSG    = "http://intapi.sgap.253.com/var/json"
	BalanceEndpoint  = "http://intapi.253.com/balance/json"
)

var (
	normalSH, _   = url.Parse(NormalEndpointSH)
	varSH, _      = url.Parse(VarEndpointSH)
	normalSG, _   = url.Parse(NormalEndpointSG)
	varSG, _      = url.Parse(VarEndpointSG)
	balanceURL, _ = url.Parse(BalanceEndpoint)
)

var (
	nodeName = map[int]string{
		NodeShanghai:  "shanghai",
		NodeSingapore: "singapore",
	}
	nodeNameReverse = map[string]int{
		"shanghai":  NodeShanghai,
		"singapore": NodeSingapore,
	}
)

var (
	SMSTypeName = map[int]string{
		SMSTypeNormal:   "普通短信",
		SMSTypeVariable: "变量短信",
	}
)
