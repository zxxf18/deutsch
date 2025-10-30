package model

// UserInfo 用户信息（JWT payload / response DTO）
type UserInfo struct {
	ID   string `json:"id"`   // UUID
	Name string `json:"name"` // 用户名
	Role string `json:"role"` // 角色: user/admin/guest
}
