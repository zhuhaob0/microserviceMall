package model

import (
	"errors"
	"final-design/pkg/mysql"
	"log"

	"github.com/gohouse/gorose/v2"
	"golang.org/x/crypto/bcrypt"
)

type AdminUser struct {
	UserId   int64  `json:"user_id"`   // Id
	Username string `json:"user_name"` // 用户名称
	Password string `json:"password"`  // 密码
	Age      int    `json:"age"`       // 年龄
}

type AdminUserModel struct{}

func NewAdminUserModel() *AdminUserModel {
	return &AdminUserModel{}
}

func (p *AdminUserModel) getTableName() string {
	return "admin_user"
}

func (p *AdminUserModel) GetUserList() ([]gorose.Data, error) {
	conn := mysql.DB()
	list, err := conn.Table(p.getTableName()).Get()
	if err != nil {
		log.Printf("Error: %v", err)
		return nil, err
	}
	return list, nil
}

func (p *AdminUserModel) CheckUser(username string, password string) (*AdminUser, error) {
	conn := mysql.DB()
	data, err := conn.Table(p.getTableName()).Where(map[string]interface{}{
		"user_name": username,
	}).First()
	if err != nil {
		log.Printf("Error: %v", err)
		return nil, err
	}
	if data == nil {
		return &AdminUser{}, errors.New("用户名或密码错误")
	}
	user := &AdminUser{
		UserId:   data["user_id"].(int64),
		Username: data["user_name"].(string),
		Password: data["password"].(string),
		Age:      int(data["age"].(int64)),
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("密码错误")
	}
	return user, nil
}

func (p *AdminUserModel) CreateUser(user *AdminUser) error {
	conn := mysql.DB()
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	_, err := conn.Table(p.getTableName()).Data(map[string]interface{}{
		"user_name": user.Username,
		"password":  string(hashPassword),
		"age":       user.Age,
	}).Insert()
	if err != nil {
		log.Printf("Error: %v", err)
		return err
	}
	return nil
}
