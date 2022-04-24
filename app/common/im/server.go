package im

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/wuyan94zl/chart"
	"github.com/wuyan94zl/go-zero-blog/app/internal/svc"
	"github.com/wuyan94zl/go-zero-blog/app/models/messages"
	"github.com/wuyan94zl/go-zero-blog/app/models/sendqueue"
	"net/http"
	"strconv"
	"time"
)

const (
	AesKey         = "wuyan94zl1asdfghjklqwertyuiopzas"
	publicChanelId = "wuyan94zl:im:public"
	sendMessage    = 100
)

type cliDetail struct {
	NickName  string `json:"nick_name"`
	Phone     string `json:"phone"`
	HeadProto string `json:"head_proto"`
}

func Run(ctx *svc.ServiceContext) {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		token, _ := strconv.Atoi(r.FormValue("_token"))
		info, err := ctx.UserModel.FindOne(context.Background(), int64(token))
		if err != nil {
			fmt.Println("ws 连接错误：", err)
			return
		}
		chart.NewServer(w, r, uint64(info.Id), cliDetail{NickName: info.NickName, Phone: info.Mobile, HeadProto: "test uri"}, &data{ctx: ctx})
	})
	fmt.Printf("Starting Ws Server at 0.0.0.0:9988...\n")
	err := http.ListenAndServe(":9988", mux)
	if err != nil {
		fmt.Println("ws err：", err)
		return
	}
}

type data struct {
	ctx *svc.ServiceContext
}

// 发送消息回调
func (d *data) SendMessage(msg chart.Message) {
	switch msg.Type {
	case sendMessage:
		sendTime, err := time.Parse("2006-01-02 15:01:05", "2022-04-15 22:12:12")
		if err != nil {
			return
		}
		message := messages.Messages{
			ChannelId:  msg.ChannelId,
			SendUserId: int64(msg.UserId),
			Message:    msg.Content,
			CreateTime: sendTime,
		}
		d.ctx.MessageModel.Insert(context.Background(), &message)
	}
}

func (d *data) DelaySendMessage(channelId string, msg chart.Message, sent []uint64) {
	users, err := d.ctx.UserUsersModel.AllChannelIdUsers(channelId)
	if err != nil {
		return
	}
	sentMap := make(map[int64]bool)
	for _, v := range sent {
		sentMap[int64(v)] = true
	}
	msgByte, err := json.Marshal(msg)
	if err != nil {
		return
	}
	for _, user := range users {
		if _, ok := sentMap[user.UserId]; !ok {
			d.ctx.SendQueueModel.Insert(context.Background(), &sendqueue.SendQueues{UserId: user.UserId, Message: string(msgByte), SendUserId: int64(msg.UserId)})
		}
	}
}

// LoginServer 登录成功后回调
func (d *data) LoginServer(uid uint64) {
	list, _ := d.ctx.UserModel.Friends(d.ctx.UserUsersModel, int64(uid))
	var channelIds []string
	for _, v := range list {
		channelIds = append(channelIds, GenChannelIdByFriend(int64(uid), v.Id))
	}
	//channelIds = append(channelIds, publicChanelId)
	chart.JoinChannelIds(uid, channelIds...)
	go func() {
		time.Sleep(time.Second * 2)
		queues, _ := d.ctx.SendQueueModel.FindByUserId(context.Background(), int64(uid))
		for _, queue := range queues {
			SendMessageToUid(uid, uid, queue.Message, 100)
			d.ctx.SendQueueModel.Delete(context.Background(), queue.Id)
		}
	}()
}
func (d *data) LogoutServer(uid uint64) {
	// 退出登陆回调
	fmt.Println("logout ", uid)
}
func (d *data) ErrorLogServer(err error) {
	// 错误消息回调
	fmt.Println("err: ", err)
}
