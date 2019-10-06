package base

type DataCode int

const (
	SUCCESS                DataCode = 200
	USE_NOT_EXIST          DataCode = 201
	USE_OR_PASSSWORD_ERROR DataCode = 202
	PARAMETER_ERROR              DataCode = 1100
	INVALID_TOKEN          DataCode = 401
	NOT_LOGIN              DataCode = 402
	SERVER_ERROR           DataCode = 500
)

func (code DataCode) String() string {
	switch code {
	case SUCCESS:
		return ""
	case USE_NOT_EXIST:
		return "用户不存在"
	case USE_OR_PASSSWORD_ERROR:
		return "用户或密码错误"
	case INVALID_TOKEN:
		return "无效的登录用户"
	case NOT_LOGIN:
		return "请先登录"
	case SERVER_ERROR:
		return "服务异常"
	case PARAMETER_ERROR:
		return "参数错误"
	default:
		return "服务异常"
	}
}
