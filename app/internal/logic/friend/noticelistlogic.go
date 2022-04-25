package friend

import (
	"context"
	"encoding/json"

	"github.com/wuyan94zl/go-zero-blog/app/internal/svc"
	"github.com/wuyan94zl/go-zero-blog/app/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type NoticeListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewNoticeListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NoticeListLogic {
	return &NoticeListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *NoticeListLogic) NoticeList() (resp *types.NoticeListResponse, err error) {
	id, _ := l.ctx.Value("id").(json.Number).Int64()
	notice, err := l.svcCtx.NoticeModel.GetListByUserId(id)
	if err != nil {
		return nil, err
	}
	var rlt types.NoticeListResponse
	for _, v := range notice {
		rlt.List = append(rlt.List, types.Notice{
			Id: v.Id, Tp: v.Tp, IsAgree: v.IsAgree, Status: v.Status,
			NickName: v.NickName, Content: v.Content, CreateTime: v.CreateTime.Format("2006-01-06 15:04:01"),
		})
	}
	return &rlt, nil
}