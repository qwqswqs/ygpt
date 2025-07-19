package system

import (
	"errors"
	"fmt"
	"github.com/gofrs/uuid/v5"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"ygpt/server/global"
	systemModel "ygpt/server/model/system"
	"ygpt/server/model/yunguan/res"
	"ygpt/server/model/yunguan/system"
	"ygpt/server/service/compute/websocket/protocol"
	"ygpt/server/utils"
)

type ModifyService struct{}

func (m *ModifyService) YgModify(message *protocol.YgModifyRequest) error {
	modifyRet := protocol.YgModifyResponse{}
	modifyRet.ID = message.ID
	modifyRet.Command = protocol.YgModifyRetCmd
	modifyRet.Method = message.Method
	var systemConfig system.SystemConfig
	switch message.Method {
	case 1:
		//TODO: 更新系统配置表,配置具体Key的值待定，注意更改
		err := global.GVA_DB.Where("key", "access").First(&systemConfig).Error
		if err != nil {
			return err
		} else {
			systemConfig.Value = fmt.Sprintf("%v", message.Access)
			err = global.GVA_DB.Model(&systemConfig).Where("id", systemConfig.ID).Updates(&systemConfig).Error
			if err != nil {
				return err
			} else {
				modifyRet.RetCode = 0
				modifyRet.RetMsg = "修改成功"
			}
		}
	case 2:
		//TODO: 更新系统配置表，配置具体Key的值待定，注意更改
		err := global.GVA_DB.Where("key", "autoUpgrade").First(&systemConfig).Error
		if err != nil {
			return err
		} else {
			systemConfig.Value = fmt.Sprintf("%v", message.AutoUpgrade)
			err = global.GVA_DB.Model(&systemConfig).Where("id", systemConfig.ID).Updates(&systemConfig).Error
			if err != nil {
				return err
			} else {
				modifyRet.RetCode = 0
				modifyRet.RetMsg = "修改成功"
			}
		}
	case 3:
		//更新用户表，更改管理员
		var authorities []systemModel.SysAuthority
		authorities = append(authorities, systemModel.SysAuthority{
			AuthorityId: 3,
		})
		var user systemModel.SysUser

		err := global.GVA_DB.Where("authority_id = ?", 3).First(&user).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 处理未找到记录的情况
				// 否则 附加uuid 注册
				user.UUID = uuid.Must(uuid.NewV4())
				err = global.GVA_DB.Create(&user).Error
				if err != nil {
					return err
				}
			} else if err != nil {
				// 处理其他类型的错误
				global.GVA_LOG.Error("Error querying database:", zap.Error(err))
				return err
			}
		} else {
			newUserInfo := &systemModel.SysUser{
				UUID:        uuid.UUID{},
				Username:    message.AdminModify.Contact,
				NickName:    message.AdminModify.NickName,
				Password:    utils.BcryptHash(message.AdminModify.Password),
				HeaderImg:   "",
				AuthorityId: 3,
				Authorities: authorities,
				Enable:      1,
				Phone:       message.AdminModify.Contact,
				Email:       message.AdminModify.Email,
				Source:      2,
			}
			err = global.GVA_DB.Model(&user).Where("id", user.ID).Updates(&newUserInfo).Error
			if err != nil {
				return err
			}
		}
		modifyRet.RetCode = 0
		modifyRet.RetMsg = "修改成功"
	case 4:
		//TODO: 查询运管统计信息
		var resInfo res.ResInfo
		global.GVA_DB.Where("resource_type = ?", 1).Find(&resInfo)
	case 5:
		//TODO: 查询运管实例运行信息
	}
	return global.GVA_WS.WsCli.SendMessage(modifyRet)
}
