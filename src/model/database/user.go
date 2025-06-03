package model_database

import (
	"time"
)

type User struct {
	IDUser       string    `gorm:"column:idUser;primaryKey" json:"idUser"`
	NameUser     string    `gorm:"column:nameUser" json:"nameUser"`
	EmailUser    string    `gorm:"column:emailUser" json:"emailUser"`
	PasswordUser string    `gorm:"column:passwordUser" json:"-"`
	CreateUser   time.Time `gorm:"column:createUser;autoCreateTime" json:"createUser"`
	UpdateUser   time.Time `gorm:"column:updateUser;autoUpdateTime" json:"updateUser"`

	Wallets []Wallet `gorm:"foreignKey:IDUser"`
}

func (User) TableName() string {
	return "tblUser"
}
