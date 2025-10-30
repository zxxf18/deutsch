package consts

const (
	// JWT相关
	JwtTokenKey string = "jwt_token" // JWT token字符串key

	// 用户信息
	UserInfoKey string = "userinfo" // *models.UserInfo key
	UserIDKey   string = "userID"   // string (从userinfo提取)
	RoleKey     string = "role"     // string (从userinfo提取)

	// 追踪
	TraceIDKey string = "traceID"
)
