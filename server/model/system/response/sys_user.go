package response

import (
	"time"
	"ygpt/server/global"
	"ygpt/server/model/system"
)

type SysUserResponse struct {
	User system.SysUser `json:"user"`
}

type LoginResponse struct {
	User      system.SysUser `json:"user"`
	Token     string         `json:"token"`
	ExpiresAt int64          `json:"expiresAt"`
}
type GetOpsUserListRes struct {
	global.GVA_MODEL
	Username  string    `json:"userName"`
	Password  string    `json:"password"`
	NickName  string    `json:"nickName"`
	HeaderImg string    `json:"headerImg"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	Enable    int       `json:"enable"`
	LoginTime time.Time `json:"loginTime"`
	Source    int       `json:"source"`
}
