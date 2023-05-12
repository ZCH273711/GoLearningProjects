package service

import (
	"examples/go-chat/internal/dao/pool"
	"examples/go-chat/internal/model"
	"examples/go-chat/pkg/common/request"
	"examples/go-chat/pkg/common/response"
	"examples/go-chat/pkg/errors"
	"examples/go-chat/pkg/global/log"
	"github.com/google/uuid"
	"time"
)

type userService struct {
}

var UserService = new(userService)

func (u *userService) Register(user *model.User) error {
	db := pool.GetDB()
	var userCount int64
	db.Model(user).Where("username", user.Username).Count(&userCount)
	if userCount > 0 {
		return errors.New("user already exists")
	}
	user.Uuid = uuid.New().String()
	user.CreatedAt = time.Now()
	user.DeletedAt = 0

	db.Create(&user)
	return nil
}

func (u *userService) Login(user *model.User) bool {
	pool.GetDB().AutoMigrate(&user)
	log.Logger.Debug("user", log.Any("user in service", user))
	db := pool.GetDB()

	var queryUser *model.User
	db.First(&queryUser, "username = ?", user.Username)
	log.Logger.Debug("queryUser", log.Any("queryUSer", queryUser))
	user.Uuid = queryUser.Uuid
	return queryUser.Password == user.Password
}

func (u *userService) ModifyUserInfo(user *model.User) error {
	var queryUser *model.User
	db := pool.GetDB()
	db.First(&queryUser, "username = ?", user.Username)
	log.Logger.Debug("queryUser", log.Any("queryUser", queryUser))
	var nullId int32 = 0
	if nullId == queryUser.Id {
		return errors.New("用户不存在")
	}
	queryUser.Nickname = user.Nickname
	queryUser.Email = user.Email
	queryUser.Password = user.Password
	db.Save(queryUser)
	return nil
}

func (u *userService) GetUserDetails(uuid string) model.User {
	var queryUser *model.User
	db := pool.GetDB()
	db.Select("uuid", "username", "nickname", "avator").First(&queryUser, "uuid = ?", uuid)
	return *queryUser
}

func (u *userService) GetUserOrGroupByName(name string) response.SearchResponse {
	var queryUser *model.User
	db := pool.GetDB()
	db.Select("uuid", "username", "nickname", "avator").First(&queryUser, "username = ?", name)
	var queryGroup *model.Group
	db.Select("uuid", "name").First(&queryGroup, "name = ?", name)
	search := response.SearchResponse{
		User:  *queryUser,
		Group: *queryGroup,
	}
	return search
}

func (u *userService) GetUserList(uuid string) []model.User {
	db := pool.GetDB()
	var queryUser *model.User
	db.First(&queryUser, "uuid = ?", uuid)
	var nullId int32 = 0
	if nullId == queryUser.Id {
		return nil
	}
	var queryUsers []model.User
	db.Raw("select u.username, u.uuid, u.avator, from user_friends as uf join users as u on uf.friend_id=u.id where uf.user_id = ?", queryUser.Id).Scan(&queryUsers)
	return queryUsers
}

func (u *userService) AddFriend(userFriendRequest *request.FriendRequest) error {
	var queryUser *model.User
	db := pool.GetDB()
	db.First(&queryUser, "uuid = ?", userFriendRequest.Uuid)
	log.Logger.Debug("queryUser", log.Any("queryUser", queryUser))
	var nullId int32 = 0
	if nullId == queryUser.Id {
		return errors.New("用户不存在")
	}
	var friend *model.User
	db.First(&friend, "username = ?", userFriendRequest.FriendUsername)
	if nullId == friend.Id {
		return errors.New("已添加该好友")
	}
	userFriend := model.UserFriend{
		UserId:   queryUser.Id,
		FriendId: friend.Id,
	}
	var userFriendQuery *model.UserFriend
	db.First(&userFriendQuery, "user_id = ? and friend_id = ?", queryUser.Id, friend.Id)
	if userFriendQuery.ID != nullId {
		return errors.New("该用户已经是你的好友")
	}
	db.AutoMigrate(&userFriend)
	db.Save(userFriend)
	log.Logger.Debug("userFriend", log.Any("userFriend", userFriend))
	return nil
}

func (u *userService) ModifyUserAvatar(avatar string, userUuid string) error {
	var queryUser *model.User
	db := pool.GetDB()
	db.First(&queryUser, "uuid = ?", userUuid)
	var nullId int32 = 0
	if nullId == queryUser.Id {
		return errors.New("用户不存在")
	}
	db.Model(&queryUser).Update("avatar", avatar)
	return nil
}
