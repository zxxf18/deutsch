package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"

	"deutsch/internal/code"
	"deutsch/internal/pkg/jwt"
	"deutsch/internal/svc"
	"deutsch/internal/types"
	"deutsch/model/gormdb"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type LoginLogic struct {
	logx.Logger
	ctx     context.Context
	svcCtx  *svc.ServiceContext
	jwtImpl jwt.JWT
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	l := &LoginLogic{
		Logger: logx.WithContext(logx.ContextWithFields(ctx, logx.Field("logic", "login"))),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
	jwtImpl, err := jwt.New(ctx, svcCtx.Config.JWTAuth.AccessSecret, svcCtx.Config.JWTAuth.AccessExpire)
	if err != nil {
		l.Errorf("failed to get jwt object: %+v", err)
	}
	l.jwtImpl = jwtImpl
	return l
}

func (l *LoginLogic) Login(req *types.LoginRequest) (resp *types.LoginResponse, err error) {
	resp = &types.LoginResponse{}

	if req.Email == "" && req.Phone == "" {
		return nil, code.NewCodeError(code.CodeValidationError)
	}
	if req.Password == "" {
		return nil, code.NewCodeError(code.CodeInvalidCredentials)
	}

	var user *gormdb.User
	if req.Email != "" {
		user, err = l.svcCtx.UserRepo.GetByEmail(l.ctx, req.Email)
	} else {
		user, err = l.svcCtx.UserRepo.GetByPhone(l.ctx, req.Phone)
	}
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, code.NewCodeError(code.CodeInvalidCredentials)
		}
		l.Errorf("failed to get user: %+v", err)
		return nil, code.NewCodeError(code.CodeDatabaseError)
	}
	if user == nil {
		return nil, code.NewCodeError(code.CodeInvalidCredentials)
	}

	if !user.IsEnabled {
		return nil, code.NewCodeError(code.CodeUserDisabled)
	}

	hash := sha256.Sum256([]byte(req.Password))
	passwordHash := hex.EncodeToString(hash[:])
	if user.PasswordHash != passwordHash {
		return nil, code.NewCodeError(code.CodeInvalidCredentials)
	}

	jwtData := map[string]string{
		"userID": user.ID,
		"role":   user.Role,
	}
	jwtResp, err := l.jwtImpl.GenerateJWT(jwtData)
	if err != nil {
		l.Errorf("failed to generate jwt: %+v", err)
		return nil, code.NewCodeError(code.CodeInternalServerError)
	}

	resp.Base = *code.BaseSuccessResp()
	resp.Data.ID = user.ID
	resp.Data.Username = user.Username
	resp.Data.Email = user.Email
	resp.Data.Role = user.Role
	resp.Data.Nickname = user.Nickname
	resp.Data.JwtToken = jwtResp.Token
	resp.Data.Expires = jwtResp.Expire.UnixMilli()
	resp.Data.MaxRefresh = jwtResp.MaxRefresh.UnixMilli()
	resp.Data.CreatedAt = user.CreatedAt.UnixMilli()
	resp.Data.UpdatedAt = user.UpdatedAt.UnixMilli()
	return
}
