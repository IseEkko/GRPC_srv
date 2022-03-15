package handler

import (
	"context"
	"crypto/sha512"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"mx_shop/goods_srv/global"
	"mx_shop/goods_srv/model"
	"mx_shop/goods_srv/proto"
	"strings"
	"time"
)

type UserServer struct {
}

func ModelToRsponse(user model.User) proto.UserInfoResponse {
	userInfoRsp := proto.UserInfoResponse{
		Id:       user.ID,
		Password: user.Password,
		NickName: user.NickName,
		Gender:   user.Gender,
		Mobile:   user.Mobile,
		Role:     int32(user.Role),
	}
	if user.Birthday != nil {
		userInfoRsp.Birthday = uint64(user.Birthday.Unix())
	}
	return userInfoRsp
}

func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

//获取用户列表
func (userserver UserServer) GetUserList(ctx context.Context, req *proto.PageInfo) (*proto.UserListRespons, error) {
	var users []model.User
	result := global.DB.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	rep := &proto.UserListRespons{}
	rep.Total = int32(result.RowsAffected)

	//对拿到的数据进行分页
	global.DB.Scopes(Paginate(int(req.Pn), int(req.PSize))).Find(&users)

	for _, user := range users {
		userInfoRsp := ModelToRsponse(user)
		rep.Data = append(rep.Data, &userInfoRsp)
	}
	return rep, nil
}

//通过手机号码查询用户
func (userserver UserServer) GetUserByMobile(ctx context.Context, req *proto.MobileRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	reslut := global.DB.Where(&model.User{Mobile: req.Mobile}).First(&user)
	if reslut.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")

	}
	if reslut.Error != nil {
		return nil, reslut.Error
	}
	userInfoRsp := ModelToRsponse(user)
	return &userInfoRsp, nil
}

//通过id查询用户
func (userserver UserServer) GetUserById(ctx context.Context, req *proto.IdRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	reslut := global.DB.First(&user, req.Id)
	if reslut.Error != nil {
		return nil, reslut.Error
	}
	if reslut.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	userinfo := ModelToRsponse(user)
	return &userinfo, nil
}

//创建用户接口
func (userserver UserServer) CreateUser(ctx context.Context, req *proto.CreateUserinfo) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.Where(&model.User{Mobile: req.Mobile}).First(&user)
	if result.RowsAffected == 1 {
		return nil, status.Errorf(codes.AlreadyExists, "用户已存在")
	}
	user.Mobile = req.Mobile
	user.NickName = req.NickName

	//密码加密
	options := &password.Options{16, 100, 32, sha512.New}
	salt, encodedPwd := password.Encode(req.PassWord, options)
	Newpassword := fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)
	user.Password = Newpassword
	result = global.DB.Create(&user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	userInfoRsp := ModelToRsponse(user)
	return &userInfoRsp, nil
}

//更新用户
func (userserver UserServer) UpdateUser(ctx context.Context, req *proto.UpdateUserInfo) (rsp *proto.Empty, err error) {
	var user model.User
	reslut := global.DB.First(&user, req.Id)
	if reslut.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	//int 转time
	birthday := time.Unix(int64(req.Birthday), 0)
	user.NickName = req.NickName
	user.Birthday = &birthday
	user.Gender = req.Gender

	reslut = global.DB.Save(&user)
	if reslut.Error != nil {
		return nil, status.Errorf(codes.Internal, reslut.Error.Error())
	}
	rsp = &proto.Empty{}
	return rsp, nil
}

//校验密码
func (userserver UserServer) CheckPassword(ctx context.Context, req *proto.PasswordCheckInfo) (*proto.CheckReponse, error) {
	options := &password.Options{16, 100, 32, sha512.New}
	pswwordInfo := strings.Split(req.EncrytedPassword, "$")
	check := password.Verify(req.Password, pswwordInfo[2], pswwordInfo[3], options)
	return &proto.CheckReponse{Success: check}, nil
}
