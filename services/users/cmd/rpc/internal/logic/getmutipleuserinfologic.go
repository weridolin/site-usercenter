package logic

import (
	"context"

	"github.com/weridolin/site-gateway/services/users/cmd/rpc/internal/svc"
	"github.com/weridolin/site-gateway/services/users/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMutipleUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMutipleUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMutipleUserInfoLogic {
	return &GetMutipleUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMutipleUserInfoLogic) GetMutipleUserInfo(in *pb.GetMutipleUserInfoReq) (*pb.GetMutipleUserInfoResp, error) {
	var userInfos []*pb.GetUserInfoResp
	err := l.svcCtx.DB.Table("users").Where("id in (?)", in.UserIds).Find(&userInfos).Error
	if err != nil {
		return nil, err
	}
	return &pb.GetMutipleUserInfoResp{
		UserInfos: userInfos,
	}, nil
}
