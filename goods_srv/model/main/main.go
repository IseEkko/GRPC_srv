package main

import (
	"crypto/md5"
	"encoding/hex"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io"
	"log"
	"mx_shop/goods_srv/global"
	"mx_shop/goods_srv/model"
	"os"
	"time"
)

func ginMd5(code string) string {
	Md5 := md5.New()
	_, _ = io.WriteString(Md5, code)
	return hex.EncodeToString(Md5.Sum(nil))
}

var (
	DB *gorm.DB
)

func main() {
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	dsn := "mxshop:123456@tcp(120.55.71.155:3306)/mxshop?charset=utf8mb4&parseTime=True&loc=Local"
	/**
	进行日志配置，这里配置可以让他打印出sql
	*/
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Info, // 日志级别
			IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  true,        // 禁用彩色打印
		},
	)
	var err error
	//这里这个db就是生成的对象
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}
	////密码生成
	//options := &password.Options{16, 100, 32, sha512.New}
	//salt, encodedPwd := password.Encode("admin123", options)
	//Newpassword := fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)
	//for i := 0; i < 10; i++ {
	//	user := model.User{
	//		NickName: fmt.Sprintf("boby%d", i),
	//		Mobile:   fmt.Sprintf("1918203152%d", i),
	//		Password: Newpassword,
	//	}
	//	DB.Save(&user)
	//}
	var users []model.User
	result := global.DB.Find(&users)
	if result.Error != nil {
		panic(result.Error)
	}
	//_ = DB.AutoMigrate(&model.User{})
	//参数是原始密码，然后返回回来的是一个salt和加密后的密码
	//salt, encodedPwd := password.Encode("generic password", nil)
	//
	//// Using custom options

	//
	//pswwordInfo := strings.Split(Newpassword,"$")
	//check := password.Verify("gen eric password",pswwordInfo[2],pswwordInfo[3],options)
	//fmt.Println(check)
}
