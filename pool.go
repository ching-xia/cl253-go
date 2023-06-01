package go253

import (
	"sync"
)

// 短信
var msgPool = sync.Pool{
	New: func() interface{} {
		return new(message)
	},
}

// 标准发送短信记录
var sendRecordPool = sync.Pool{
	New: func() interface{} {
		return new(SendRecord)
	},
}

// 标准发送短信返回值
var sendResponsePool = sync.Pool{
	New: func() interface{} {
		return new(sendResponse)
	},
}

// 普通短信返回值
var normalMessageResponsePool = sync.Pool{
	New: func() interface{} {
		return new(normalMessageResponse)
	},
}

// 变量短信返回值
var varMessageResponsePool = sync.Pool{
	New: func() interface{} {
		return new(varMessageResponse)
	},
}
