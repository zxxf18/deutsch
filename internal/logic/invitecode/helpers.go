package invitecode

import (
	"context"
	"strings"

	"deutsch/model/gormdb"
	"deutsch/model/repository"
)

// getInviteByIDOrCode 支持用 UUID 或 Code 查询邀请码，path 参数可能是任一种
func getInviteByIDOrCode(ctx context.Context, repo repository.InviteCodeRepository, idOrCode string) (*gormdb.InviteCode, error) {
	// UUID 格式（含 '-'）用 GetByID，否则用 GetByCode
	if strings.Contains(idOrCode, "-") && len(idOrCode) == 36 {
		return repo.GetByID(ctx, idOrCode)
	}
	ic, err := repo.GetByCode(ctx, idOrCode)
	if err != nil {
		return nil, err
	}
	return ic, nil
}
