package business

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"time"
	"ygpt/server/api/v1/compute/udp"
	"ygpt/server/global"
	"ygpt/server/model/yunguan/product"
	"ygpt/server/model/yunguan/renter"
	"ygpt/server/model/yunguan/res"
	"ygpt/server/service/compute/websocket/protocol"
)

type ElementService struct{}

// SyncElementInfo 同步所有未同步的元素信息
//func (e *ElementService) SyncElementInfo() error {
//	var elementInfo product.ProductElementInfo
//	var productIDs []int
//	if err := global.GVA_DB.Model(&elementInfo).
//		Where("is_synced = ?", 1).
//		Pluck("id", &productIDs).Error; err != nil {
//		return err
//	} else {
//		for _, productID := range productIDs {
//			err = e.ElementRelease(productID)
//			if err != nil {
//				return err
//			}
//		}
//	}
//	return nil
//}

func (e *ElementService) HandleCopyAiProductToCloudResource(message *protocol.CopyAiProductToCloudResource) error {
	var renterRes renter.RenterRes
	if err := global.GVA_DB.Model(&renterRes).
		Where("resource_id = ?", message.Payload.CloudResourceID).First(&renterRes).Error; err != nil {
		return err
	}
	var details []res.DownloadDetail
	for _, detail := range message.Payload.DownloadDetails {
		details = append(details, res.DownloadDetail{
			DownloadSource: detail.DownloadSource,
			ProductID:      detail.ProductID,
			ProductName:    detail.ProductName,
			ProductType:    detail.ProductType,
			FileName:       detail.FileName,
			Url:            detail.Url,
		})
	}
	var copyAiProductToCloudResourceResponse protocol.CopyAiProductToCloudResourceResponse
	copyAiProductToCloudResourceResponse.ID = message.ID
	copyAiProductToCloudResourceResponse.Command = protocol.CopyAiRetCmd
	copyAiProductToCloudResourceResponse.RetCode = 0
	copyAiProductToCloudResourceResponse.RetMsg = "success"

	err := udp.InstanceService.Download(renterRes.SshHost, details)
	if err != nil {
		copyAiProductToCloudResourceResponse.RetCode = -400
		copyAiProductToCloudResourceResponse.RetMsg = err.Error()
	}
	//返回应答
	return global.GVA_WS.WsCli.SendMessage(copyAiProductToCloudResourceResponse)
}

// ElementRelease 向平台发送同步元素请求
func (e *ElementService) ElementRelease(elementID int) error {
	if elementID == 0 {
		return fmt.Errorf("productID is empty")
	}
	elementMsg := protocol.ElementRelease{}
	elementMsg.ID = global.GVA_WS.Counter.GetID()
	elementMsg.Command = protocol.ElementUploadCmd
	elementMsg.NodeCode = global.GVA_CONFIG.System.NodeCode
	var productEle *product.ProductElem
	var productFile *[]product.ProductFile
	if err := global.GVA_DB.Model(&product.ProductElem{}).Where("id = ?", elementID).First(&productEle).Error; err != nil {
		return err
	}
	if err := global.GVA_DB.Model(&product.ProductFile{}).Where("product_id = ?", elementID).Find(&productFile).Error; err != nil {
		return err
	}
	if productEle != nil && productFile != nil {
		elementMsg.ProductEle = *productEle
		elementMsg.ProductFiles = *productFile
	}
	//先创建通道，再发送消息，等待接收结果
	resultCh := make(chan any, 1)
	global.AddTransferChannel(int64(elementMsg.ID), resultCh)
	defer func() {
		global.DeleteTransferChannel(int64(elementMsg.ID))
	}()

	if err := global.GVA_WS.SendingWithRetry(elementMsg); err != nil {
		return err
	}
	deadLine := time.After(time.Second * 6)
	var result *protocol.ElementReleaseResponse
	select {
	case <-deadLine:
		global.GVA_LOG.Error("ElementReleaseResponse timeout")
		return fmt.Errorf("ElementReleaseResponse timeout")
	case r := <-resultCh:
		err, ok := r.(error)
		if ok {
			global.GVA_LOG.Error(err.Error())
			return err
		}
		result, ok = r.(*protocol.ElementReleaseResponse)
		if !ok {
			global.GVA_LOG.Error("error type of ElementReleaseResponse receive from task channel")
			return fmt.Errorf("error type of ElementReleaseResponse receive from task channel")
		}
		if result.RetCode != 0 {
			return fmt.Errorf("ElementReleaseResponse error：%v %v", result.RetCode, result.RetMsg)
		}
		return nil
	}
}

// HandleElementReleaseResult 处理元素同步结果应答
func (e *ElementService) HandleElementReleaseResult(message *protocol.ElementReleaseResponse) error {
	elementID := message.ElementID
	var elementInfo product.ProductElem
	if message.RetCode != 0 {
		if message.RetCode == protocol.ExistCode {
			elementInfo.IsSynced = 1
			global.GVA_DB.Where("id = ?", elementID).Updates(&elementInfo)
			return nil
		}
		global.GVA_DB.Where("id = ?", elementID).First(&elementInfo)
		elementInfo.IsSynced = 0
		global.GVA_DB.Updates(&elementInfo)
		err := fmt.Errorf("element release ret id error：%v %v", message.ID, message.RetCode)
		return err
	} else {
		global.GVA_DB.Where("id = ?", elementID).First(&elementInfo)
		elementInfo.IsSynced = 1
		global.GVA_DB.Updates(&elementInfo)

		//更新产品基本信息不用每次更新文件
		if message.RetCode == protocol.SuccessCodeCmd {
			//将文件同步到平台,失败退出，产品视为未同步
			err := e.elementFileSync(elementID)
			if err != nil {
				global.GVA_LOG.Error(err.Error())
				return err
			}
		}
	}
	return nil
}

// ElementManage 上架、下架元素请求
func (e *ElementService) ElementManage(method int8, elementID ...int) error {
	if method == 0 {
		return fmt.Errorf("报文参数不全")
	}
	manageMsg := protocol.ElementManageRequest{}
	manageMsg.ID = global.GVA_WS.Counter.GetID()
	manageMsg.Command = protocol.ElementManageCmd
	manageMsg.NodeCode = global.GVA_CONFIG.System.NodeCode
	//method : 1.查询，2.上架，3.下架
	manageMsg.Method = method
	//要上架、下架的产品ID
	if method == 2 || method == 3 {
		manageMsg.ElementID = elementID[0]
	}
	return global.GVA_WS.SendingWithRetry(manageMsg)
}

// HandleElementManageResult 处理元素上架、下架结果应答
func (e *ElementService) HandleElementManageResult(message *protocol.ElementManageResponse) error {
	var elementInfo product.ProductElementInfo
	elementID := message.ElementID
	global.GVA_DB.Where("id = ?", elementID).First(&elementInfo)
	if message.RetCode != 0 {
		return fmt.Errorf("element manage ret id error：%v %v", message.ID, message.RetCode)
	} else {
		if message.Method == 1 {
			//解析使用次数的结果，存数据库!!!暂时不用这个功能
		}
		if message.RetCode == 2 {
			elementInfo.Status = 1
			global.GVA_DB.Save(&elementInfo)
		}
		if message.RetCode == 2 {
			elementInfo.Status = 0
			global.GVA_DB.Save(&elementInfo)
		}
	}
	return nil
}

func (e *ElementService) elementFileSync(elementID int) error {
	// 查询数据库获取产品文件信息
	var productFile []product.ProductFile
	if result := global.GVA_DB.Model(&product.ProductElem{}).Where("id = ?", elementID).Find(&productFile); result.Error != nil || len(productFile) == 0 {
		return fmt.Errorf("element info not found")
	}
	// 判断云管算力是否用同一个minio
	if global.GVA_CONFIG.Minio.EndPoint == global.GVA_CONFIG.SLMinio.EndPoint {
		if global.GVA_CONFIG.Minio.BucketName == global.GVA_CONFIG.SLMinio.BucketName {
			global.GVA_LOG.Info("云管算力使用同一个minio，不需要同步文件")
			return nil
		}
	}

	// 初始化 MinIO 客户端
	cloudMinio, err := newMinioClient(global.GVA_CONFIG.Minio.EndPoint, global.GVA_CONFIG.Minio.AccessKey, global.GVA_CONFIG.Minio.SecretKey, global.GVA_CONFIG.Minio.UseSSL)
	if err != nil {
		return err
	}
	slMinio, err := newMinioClient(global.GVA_CONFIG.SLMinio.EndPoint, global.GVA_CONFIG.SLMinio.AccessKey, global.GVA_CONFIG.SLMinio.SecretKey, global.GVA_CONFIG.SLMinio.UseSSL)
	if err != nil {
		return err
	}
	go func() {
		// 将文件从云管 MinIO 上传到算力 MinIO
		for _, file := range productFile {
			if file.Name == "" || *file.Size <= 0 {
				continue // 跳过无效文件
			}
			dl := 60 * 5
			key := 0
			for {
				//确定是否存在文件
				if _, err := cloudMinio.StatObject(context.Background(), global.GVA_CONFIG.SLMinio.BucketName, file.Name, minio.StatObjectOptions{}); err == nil {
					break
				}
				if key >= dl {
					break
				}
				time.Sleep(5 * time.Second)
				key += 5
			}

			// 从云管 MinIO 获取对象
			obj, err := cloudMinio.GetObject(context.Background(), global.GVA_CONFIG.Minio.BucketName, file.Name, minio.GetObjectOptions{})
			if err != nil {
				if obj != nil {
					obj.Close()
				}
				continue
			}
			// 将对象上传到算力 MinIO
			_, err = slMinio.PutObject(context.Background(), global.GVA_CONFIG.SLMinio.BucketName, file.Name, obj, *file.Size, minio.PutObjectOptions{ContentType: "application/octet-stream"})
			if err != nil {
				obj.Close() // 确保资源释放
				continue
			}
			obj.Close() // 确保资源释放
		}
	}()

	return nil
}

// 提取 MinIO 客户端初始化逻辑
func newMinioClient(endpoint, accessKey, secretKey string, useSSL int) (*minio.Client, error) {
	var secure bool
	if useSSL == 1 {
		secure = true
	} else {
		secure = false
	}
	return minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: secure,
	})
}
