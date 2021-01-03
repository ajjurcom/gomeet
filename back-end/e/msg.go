package e

var MsgFlags = map[int]string {
	SUCCESS:        "ok",
	BACK_ERROR:     "background err, please look the log",
	INVALID_PARAMS: "请求参数错误",
	NOT_POWER:      "没有权限",
	NOT_EXITS:      "No Exists",
	URL_ERROR:      "请求链接错误，自动跳转到首页",
}

// GetMsg 根据返回码获取提示信息
//
// 参数：code 返回码
//
// 返回值：string 提示信息
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return "无信息"
}
