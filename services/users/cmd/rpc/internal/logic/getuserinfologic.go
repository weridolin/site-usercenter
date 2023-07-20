package logic

import (
	"context"

	"github.com/weridolin/site-gateway/services/users/cmd/rpc/internal/svc"
	"github.com/weridolin/site-gateway/services/users/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserInfoLogic) GetUserInfo(in *pb.GetUserInfoReq) (*pb.GetUserInfoResp, error) {
	user, err := l.svcCtx.UserModel.QueryUser(map[string]interface{}{"id": in.UserId}, l.svcCtx.DB)
	if err != nil {
		return nil, err
	} else {
		return &pb.GetUserInfoResp{
			UserId:       int64(user.ID),
			UserName:     user.Username,
			UserMobile:   user.Phone,
			UserEmail:    user.Email,
			UserAvatar:   user.Avatar,
			Age:          int64(user.Age),
			Gender:       int32(user.Gender),
			IsSuperAdmin: user.IsSuperAdmin,
		}, nil
	}
}
