package model_database

import (
	"time"
)

type Wallet struct {
	IDWallet      string    `gorm:"column:idWallet;primaryKey" json:"idWallet"`
	IDUser        string    `gorm:"column:idUser" json:"idUser"`
	AddressWallet string    `gorm:"column:addressWallet" json:"addressWallet"`
	CreateWallet  time.Time `gorm:"column:createWallet;autoCreateTime" json:"createWallet"`
	UpdateWallet  time.Time `gorm:"column:updateWallet;autoUpdateTime" json:"updateWallet"`

	// User User `gorm:"foreignKey:IDUser;references:IDUser"`
}

func (Wallet) TableName() string {
	return "tblWallet"
}
