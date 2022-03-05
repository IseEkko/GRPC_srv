package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"mx_shop/user_srv/proto"
)

var userClient proto.UserClient
var conn *grpc.ClientConn

func Init() {
	var err error
	conn, err = grpc.Dial("127.0.0.1:50051", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	userClient = proto.NewUserClient(conn)
}

//获取用户列表
func TestGetUserList() {
	rsp, err := userClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    1,
		PSize: 5,
	})
	if err != nil {
		panic(err)
	}
	for _, user := range rsp.Data {
		fmt.Println(user.Mobile, user.NickName, user.Password)
		checkrsp, err := userClient.CheckPassword(context.Background(), &proto.PasswordCheckInfo{
			Password:         "admin123",
			EncrytedPassword: user.Password,
		})
		if err != nil {
			panic(err)
		}
		fmt.Println(checkrsp.Success)
	}
}

//测试创建用户
func TestCreatUser() {
	for i := 0; i < 10; i++ {
		rsp, err := userClient.CreateUser(context.Background(), &proto.CreateUserinfo{
			NickName: fmt.Sprintf("lzqqqz%d", i),
			PassWord: "admin123",
			Mobile:   fmt.Sprintf("1918203123%d", i),
		})
		if err != nil {
			panic(err)
		}
		fmt.Println(rsp.Id)
	}

}

//通过用户电话进行查询
func TestByMobileUser() {
	rsp, err := userClient.GetUserByMobile(context.Background(), &proto.MobileRequest{Mobile: "19182031524"})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp)
}

//通过用户的id进行查询
func TestByIdUser() {
	rsp, err := userClient.GetUserById(context.Background(), &proto.IdRequest{Id: 12})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp)
}

//更改用户信息
func TestUpdatesUser() {
	rsp, err := userClient.UpdateUser(context.Background(), &proto.UpdateUserInfo{
		Id:       12,
		NickName: "ceshi23",
		Birthday: 0,
		Gender:   "man",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp)
}
func main() {
	Init()
	//
	TestByIdUser()
	defer conn.Close()

}
