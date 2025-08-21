package main

// import (
// 	"context"
// 	"fmt"
// 	"usersrv/proto"

// 	"google.golang.org/grpc"
// )

// // func main() {
// // 	options := &password.Options{16, 100, 32, sha512.New}
// // 	salt, encodedPwd := password.Encode("generic password", options)
// // 	newPassword := fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, &encodedPwd)
// // 	fmt.Println(len(newPassword))
// // 	fmt.Println(newPassword)

// //		passwordInfo := strings.Split(newPassword, "$")
// //		fmt.Println(passwordInfo)
// //		check := password.Verify("generic password", passwordInfo[2], passwordInfo[3], options)
// //		fmt.Println(check)
// //	}
// // const (
// // 	PasswordCost = 12 //密码加密难度
// // )

// // func main() {
// // 	password := "123456"
// // 	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), PasswordCost)
// // 	bcryptPW := string(bytes)
// // 	fmt.Println(bcryptPW)

// // 	err := bcrypt.CompareHashAndPassword([]byte(bcryptPW), []byte(password))
// // 	if err != nil {
// // 		fmt.Println(err)
// // 	}
// // 	fmt.Println("success")
// // }

// // var (
// // 	_db *gorm.DB
// // )

// // func TestUserCreate() {
// // 	// initialize.MySQL()
// // 	// _db = global.NewDBClient(context.Background())
// // 	// db := global.DB
// // 	rawPW := "admin123"
// // 	var user model.User
// // 	_ = user.SetPassword(rawPW)
// // 	bytes, err := bcrypt.GenerateFromPassword([]byte(rawPW), model.PasswordCost)
// // 	if err != nil {
// // 		panic("用户密码加密失败")
// // 	}
// // 	newPW := string(bytes)
// // 	for i := 0; i < 10; i++ {
// // 		user := &model.User{
// // 			NickName: fmt.Sprintf("lucien%d", i),
// // 			Mobile:   fmt.Sprintf("176349228%d", i),
// // 			Password: newPW,
// // 		}
// // 		global.DB.Save(&user)
// // 	}
// // }

// var userClient proto.UserClient
// var conn *grpc.ClientConn

// func InitClient() {
// 	var err error
// 	conn, err = grpc.Dial("127.0.0.1:8083", grpc.WithInsecure())
// 	if err != nil {
// 		panic(err)
// 	}
// 	userClient = proto.NewUserClient(conn)
// }

// func TestCreateUser() {
// 	for i := 0; i < 10; i++ {
// 		rsp, err := userClient.CreateUser(context.Background(), &proto.CreateUserInfo{
// 			NickName: fmt.Sprintf("bobby%d", i),
// 			Mobile:   fmt.Sprintf("1878222222%d", i),
// 			PassWord: "admin123",
// 		})
// 		if err != nil {
// 			panic(err)
// 		}
// 		fmt.Println("创建用户成功", rsp.Id)
// 	}
// }

// // func TestCreateUser(t *testing.T) {

// // 	InitClient()
// // 	err := CreateUser()
// // 	var errStr string
// // 	errStr = "rpc error: code = AlreadyExists desc = 用户已存在 [recovered] panic: rpc error: code = AlreadyExists desc = 用户已存在"
// // 	if err.Error() == errStr {
// // 		t.Logf("用户已存在测试通过")
// // 	} else {
// // 		t.Logf(err.Error())
// // 	}

// // 	conn.Close()
// // }

// func TestGetUserList() {
// 	rsp, err := userClient.GetUserList(context.Background(), &proto.PageInfo{
// 		Pn:    1,
// 		PSize: 5,
// 	})
// 	if err != nil {
// 		fmt.Println("查询用户失败")
// 		panic(err)
// 	}
// 	// fmt.Println(rsp)
// 	for _, user := range rsp.Data {
// 		// fmt.Println(user.NickName, user.Mobile, user.PassWord)
// 		checkRsp, err := userClient.CheckPassWord(context.Background(), &proto.PasswordCheckInfo{
// 			PassWord:          "admin123",
// 			EncryptedPassWord: user.PassWord,
// 		})
// 		if err != nil {
// 			panic(err)
// 			// fmt.Println("获取用户信息失败")
// 		}
// 		fmt.Println(checkRsp.Success)
// 	}
// }

// // func TestGetUserList(t *testing.T) {
// // 	initialize.MySQL()
// // 	InitClient()
// // 	err := GetUserList()
// // 	errStr1 := "fail1"
// // 	errStr2 := "fail2"
// // 	if err.Error() == errStr1 {
// // 		t.Logf("fail1")
// // 	} else if err.Error() == errStr2 {
// // 		t.Logf("fail2")
// // 	}
// // 	conn.Close()
// // }

// func main() {
// 	InitClient()

// 	// TestCreateUser()
// 	TestGetUserList()

// 	conn.Close()
// }
