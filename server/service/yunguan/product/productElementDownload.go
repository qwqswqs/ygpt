package product

import (
	"ygpt/server/global"
	"ygpt/server/model/yunguan/product"
)

type ProductElemDownloadService struct {
}

func (i *ProductElemDownloadService) ListElementDownload(privateIP string) ([]product.ProductElementDownload, error) {
	var info []product.ProductElementDownload
	//TODO : 筛选下载完成的文件，目前是全部返回
	err := global.GVA_DB.Model(product.ProductElementDownload{}).Debug().Where("private_ip = ?", privateIP).Find(&info).Error
	if err != nil {
		return nil, err
	}
	return info, nil
}
