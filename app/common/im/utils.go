package im

import (
	"github.com/wuyan94zl/go-zero-blog/app/common/utils"
	"strconv"
)

func GenChannelIdByFriend(userId, friendId int64) string {
	num := userId + friendId
	return genChannelId(strconv.FormatInt(num, 10))
}

func GenChannelIdByGroup(groupName string) string {
	return genChannelId(groupName)
}

func genChannelId(key string) string {
	return utils.Md5ByString(key)

}


