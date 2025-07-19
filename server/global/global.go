package global

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/CDUwenbojin/websocket"
	"github.com/gin-gonic/gin"
	"github.com/qiniu/qmgo"
	"github.com/redis/go-redis/v9"
	"math/big"
	"regexp"
	"strconv"
	"sync"
	"time"
	"ygpt/server/service/compute/websocket/ws"
	"yunion.io/x/onecloud/pkg/mcclient"

	"github.com/songzhibin97/gkit/cache/local_cache"
	"ygpt/server/utils/timer"

	"golang.org/x/sync/singleflight"

	"go.uber.org/zap"

	"ygpt/server/config"

	gmap "github.com/doraemonkeys/sync-gmap"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var (
	GVA_WS             *ws.WsClient
	GVA_MSG_TRANSFER   gmap.SyncMap[websocket.SessionID, gmap.SyncMap[int64, chan any]] //云管平台作为服务器与实例中的代理程序的消息通道
	GVA_SLMSG_TRANSFER gmap.SyncMap[int64, chan any]                                    //云管平台作为客户端与算力平台通信的消息通道
	GVA_MUTEX          sync.Mutex
	GVA_DB             *gorm.DB
	GVA_DBList         map[string]*gorm.DB
	GVA_REDIS          redis.UniversalClient
	GVA_REDISList      map[string]redis.UniversalClient
	GVA_MONGO          *qmgo.QmgoClient
	GVA_CONFIG         config.Server
	GVA_CONFIG_PATH    string
	GVA_VP             *viper.Viper
	// GVA_LOG    *oplogging.Logger
	GVA_LOG                 *zap.Logger
	GVA_Timer               timer.Timer = timer.NewTimerTask()
	GVA_Concurrency_Control             = &singleflight.Group{}
	GVA_ROUTERS             gin.RoutesInfo
	GVA_ACTIVE_DBNAME       *string
	BlackCache              local_cache.Cache
	lock                    sync.RWMutex
	CloudUrl                = "https://10.10.1.23:30888/servers"
	CloudHeaders            = map[string]string{
		"Accept":          "*/*",
		"Accept-Encoding": "*",
		"Content-Length":  "0",
		"Content-Type":    "application/json",
		"User-Agent":      "yunioncloud-go/201708",
		"X-Auth-Token":    "a68007dbf9152695b679b4418de6eb82f1a82b2e9476fdd9925a708524a28012",
	}
)

// GetGlobalDBByDBName 通过名称获取db list中的db
func GetGlobalDBByDBName(dbname string) *gorm.DB {
	lock.RLock()
	defer lock.RUnlock()
	return GVA_DBList[dbname]
}

// MustGetGlobalDBByDBName 通过名称获取db 如果不存在则panic
func MustGetGlobalDBByDBName(dbname string) *gorm.DB {
	lock.RLock()
	defer lock.RUnlock()
	db, ok := GVA_DBList[dbname]
	if !ok || db == nil {
		panic("db no init")
	}
	return db
}

func GetRedis(name string) redis.UniversalClient {
	redis, ok := GVA_REDISList[name]
	if !ok || redis == nil {
		panic(fmt.Sprintf("redis `%s` no init", name))
	}
	return redis
}

func GetCloudClientSession() (*mcclient.ClientSession, error) {
	client := mcclient.NewClient(GVA_CONFIG.Cloudpods.AuthURL,
		GVA_CONFIG.Cloudpods.Timeout,
		GVA_CONFIG.Cloudpods.Debug,
		GVA_CONFIG.Cloudpods.Insecure,
		"",
		"")
	token, err := client.Authenticate(GVA_CONFIG.Cloudpods.Uname, GVA_CONFIG.Cloudpods.Passwd, GVA_CONFIG.Cloudpods.DomainName, GVA_CONFIG.Cloudpods.TenantName, GVA_CONFIG.Cloudpods.TenantDomain)
	if err != nil {
		GVA_LOG.Error("cloudpods auth failed", zap.Error(err))
		return nil, err
	}
	s := client.NewSession(context.Background(),
		GVA_CONFIG.Cloudpods.Region,
		"",
		GVA_CONFIG.Cloudpods.EndpointType,
		token,
	)

	return s, err
	//params := jsonutils.NewDict()
	//params.Set("scope", jsonutils.NewString("system"))
	//params.Set("limit", jsonutils.NewString("2"))
	//
	//result, err := modules.Disks.List(s, params)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("%s", jsonutils.Marshal(result).PrettyString())
}

func HandleMsgTransfer(demand string, sessionId websocket.SessionID, id int64, newMsgCh chan any) (msgCh chan any, ok bool) {
	GVA_MUTEX.Lock()
	defer GVA_MUTEX.Unlock()

	msgTransfer, isExisted := GVA_MSG_TRANSFER.Load(sessionId)
	if !isExisted {
		msgTransfer = *gmap.NewSyncMap[int64, chan any]()
	}

	switch demand {
	case "store":
		if newMsgCh != nil {
			msgTransfer.Store(id, newMsgCh)
		}
		GVA_MSG_TRANSFER.Store(sessionId, msgTransfer)
	case "delete":
		msgTransfer.Delete(id)
		GVA_MSG_TRANSFER.Store(sessionId, msgTransfer)
	case "get":
		msgCh, ok = msgTransfer.Load(id)
		return msgCh, ok
	default:
	}
	return nil, false
}

func AddTransferChannel(msgID int64, ch chan any) {
	if _, ok := GVA_SLMSG_TRANSFER.Load(msgID); !ok {
		GVA_SLMSG_TRANSFER.Store(msgID, ch)
	}
}

func DeleteTransferChannel(msgID int64) {
	if _, ok := GVA_SLMSG_TRANSFER.Load(msgID); ok {
		GVA_SLMSG_TRANSFER.Delete(msgID)
	}
}

func GenerateRandomLetter() string {
	// 随机生成一个字母（小写字母）
	letters := "abcdefghijklmnopqrstuvwxyz"

	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
	if err != nil {
		return "a"
	}

	return string(letters[n.Int64()])
}
func TimestampToLetters(timestamp int64) string {
	// 将时间戳转换为字符串，然后通过 ASCII 码转换为字母
	timestampStr := strconv.FormatInt(timestamp, 10)
	// 确保至少有8位数字，不足的部分在前面补零
	if len(timestampStr) < 8 {
		timestampStr = fmt.Sprintf("%08s", timestampStr)
	}
	// 截取后8位数字
	if len(timestampStr) > 8 {
		timestampStr = timestampStr[len(timestampStr)-8:]
	}
	result := ""
	for _, char := range timestampStr {
		// 将字符转换为整数
		digit, _ := strconv.Atoi(string(char))
		// 加上97得到新的ASCII值
		asciiVal := digit + 97
		// 转换为对应的ASCII字符
		result += string(rune(asciiVal))
	}
	return result
}

func GenerateUniqueName(renterId int) string {
	// 生成纳秒级时间戳并转换为字符串，取前 8 位以控制长度
	timeStr := TimestampToLetters(time.Now().UnixNano())

	// 将租户 ID 转换为字符串
	renterIdStr := strconv.Itoa(renterId)

	// 组合时间戳、随机字符串和租户 ID
	uniqueName := timeStr + "-" + renterIdStr + "-"

	// 确保名称不超过 20 个字符
	if len(uniqueName) > 20 {
		uniqueName = uniqueName[:20]
	}

	// 随机生成一个字母作为结尾
	if !regexp.MustCompile(`[a-zA-Z]$`).MatchString(uniqueName) {
		uniqueName = uniqueName[:len(uniqueName)-1] + GenerateRandomLetter()
	}
	return uniqueName
}
