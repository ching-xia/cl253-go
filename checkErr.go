package go253

import "github.com/pkg/errors"

// 检查返回码，并校验是否可以继续发送
func checkResponseCode(code string) (bool, error) {
	switch code {
	case "0":
		return true, nil
	case "101":
		return false, errors.New("账号不存在")
	case "102":
		return false, errors.New("账号密码错误")
	case "106":
		return true, errors.New("短信内容过长")
	case "108":
		return true, errors.New("手机号码格式错误")
	case "110":
		return false, errors.New("余额不足")
	case "112":
		return false, errors.New("产品配置错误")
	case "114":
		return false, errors.New("IP地址认证错误")
	case "115":
		return false, errors.New("未开通产品权限")
	case "123":
		return true, errors.New("短信内容为空")
	case "128":
		return false, errors.New("账号长度错误")
	case "129":
		return false, errors.New("产品价格配置错误")
	default:
		return false, errors.Errorf("未知错误: %s", code)
	}
}

func CheckStatusCode(code string) error {
	switch code {
	case "DELIVRD":
		return nil
	case "UNKNOWN":
		return errors.New("未知短信状态")
	case "REJECTD":
		return errors.New("短信被短信中心拒绝")
	case "MBBLACK":
		return errors.New("目的号码是黑名单号码")
	case "SM11":
		return errors.New("网关验证号码格式错误")
	case "SM12":
		return errors.New("253平台验证号码格式错误")
	default:
		return errors.Errorf("未知错误: %s", code)
	}
}
