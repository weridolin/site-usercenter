package logic

import (
	"context"

	"github.com/weridolin/site-gateway/services/users/cmd/rpc/internal/svc"
	"github.com/weridolin/site-gateway/services/users/cmd/rpc/pb"
	"github.com/weridolin/site-gateway/services/users/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type TokenValidateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTokenValidateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TokenValidateLogic {
	return &TokenValidateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TokenValidateLogic) TokenValidate(in *pb.TokenValidateReq) (*pb.TokenValidateResp, error) {
	_, err := models.ParseToken(in.Token, l.svcCtx.Config.JwtAuth.AccessSecret)
	if err != nil {
		return &pb.TokenValidateResp{
			IsValid: false,
		}, nil
	} else {
		return &pb.TokenValidateResp{
			IsValid: true,
		}, nil
	}

}
