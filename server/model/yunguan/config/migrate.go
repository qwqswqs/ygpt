package config

import (
	"bytes"
	"encoding/binary"
	"time"
	"ygpt/server/global"
	"ygpt/server/utils"
)

func Migrate() {
	err := global.GVA_DB.AutoMigrate(&LicenseConfig{})
	if err != nil {
		panic(err)
	} else {
		var total int64
		err := global.GVA_DB.Model(LicenseConfig{}).Count(&total).Error
		if err != nil {
			return
		}
		if total == 0 {
			global.GVA_DB.Model(LicenseConfig{}).Create(&LicenseConfig{
				StartTime:   time.Now().Unix(),
				EndTime:     time.Date(2050, 6, 1, 12, 0, 0, 0, time.UTC).Unix(),
				DeviceNum:   ToBigEndianBytes(10),
				RenterNum:   ToBigEndianBytes(1),
				AdminNum:    ToBigEndianBytes(1),
				LicenseUser: "admin",
				LicensePwd:  utils.MD5V([]byte("123456")),
			})
		}
	}
	err = global.GVA_DB.AutoMigrate(&AlertConfig{})
	if err != nil {
		panic(err)
	}
}

// 字节序转换函数
func ToBigEndianBytes(value uint32) []byte {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.BigEndian, value)
	return buf.Bytes()
}

func FromBigEndianBytes(data []byte) uint32 {
	var value uint32
	buf := bytes.NewReader(data)
	_ = binary.Read(buf, binary.BigEndian, &value)
	return value
}
