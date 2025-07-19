package minio

import (
	"bytes"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"mime/multipart"
)

// DownloadFile 下载文件返回流和错误
func DownloadFile(endpoint, accessKeyID, secretAccessKey, bucketName, ObjectName string) (io.Reader, error) {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false,
	})
	object, err := minioClient.GetObject(context.Background(), bucketName, ObjectName, minio.GetObjectOptions{})
	return object, err
}

func GetBucketSize(client *minio.Client, bucketName string) (int64, error) {
	var totalSize int64
	// 列出存储桶中的对象
	objectsCh := client.ListObjects(context.Background(), bucketName, minio.ListObjectsOptions{Recursive: true})
	for object := range objectsCh {
		if object.Err != nil {
			return 0, object.Err
		}
		totalSize += object.Size
	}

	return totalSize, nil
}

func CreateBucket(endpoint, accessKeyID, secretAccessKey, bucketName string) (err error) {

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false,
	})
	err = minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{Region: "us-east-1"})
	if err != nil {
		// 检查存储桶是否已经存在
		if minio.ToErrorResponse(err).Code == "BucketAlreadyOwnedByYou" {
			fmt.Printf("存储桶 %s 已存在\n", bucketName)
			return err
		} else {
			fmt.Printf("创建存储桶失败:%s %s", bucketName, err)
			return err
		}
	} else {
		fmt.Printf("存储桶 %s 创建成功\n", bucketName)
		return err
	}
	return err
}
func RemoveBucket(endpoint, accessKeyID, secretAccessKey, bucketName string) (err error) {
	// 检查存储桶是否存在
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false,
	})
	exists, err := minioClient.BucketExists(context.Background(), bucketName)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if !exists {
		fmt.Printf("存储桶: %s 不存在.\n", bucketName)
		return err
	}

	// 列出存储桶中的所有对象
	objectCh := minioClient.ListObjects(context.Background(), bucketName, minio.ListObjectsOptions{})
	// 循环删除每个对象
	for obj := range objectCh {
		if obj.Err != nil {
			fmt.Println(obj.Err)
		}
		// 删除对象
		err := minioClient.RemoveObject(context.Background(), bucketName, obj.Key, minio.RemoveObjectOptions{})
		if err != nil {
			fmt.Printf("删除文件失败 %s: %v\n", obj.Key, err)
			return err
		} else {
			fmt.Printf("删除文件: %s\n", obj.Key)
		}
	}

	fmt.Printf("存储桶 %s 中所有文件已被删除.", bucketName)
	err = minioClient.RemoveBucket(context.Background(), bucketName)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Printf("存储桶: %s 已被成功删除.\n", bucketName)
	return err
}
func UploadFile(endpoint, accessKeyID, secretAccessKey, bucketName, fileName string, file *multipart.FileHeader) (err error) {
	buf := new(bytes.Buffer)
	open, err := file.Open()
	if err != nil {
		return err
	}
	_, err = io.Copy(buf, open)
	if err != nil {
		return err
	}

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false,
	})
	_, err = minioClient.PutObject(context.Background(), bucketName, fileName, buf, file.Size, minio.PutObjectOptions{})
	if err != nil {
		fmt.Printf("上传文件失败: %s", err)
		return err
	}
	fmt.Println("上传文件成功")
	return err
}
func RemoveFile(endpoint, accessKeyID, secretAccessKey, bucketName, objectName string) (err error) {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false,
	})
	err = minioClient.RemoveObject(context.Background(), bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("删除文件失败: %v", err)
	}
	return err
}
