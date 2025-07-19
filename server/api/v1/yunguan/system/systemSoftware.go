package system

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
	"log"
	"net"
	"net/http"
	"os/exec"
	"strings"
	"time"
	"ygpt/server/global"
	"ygpt/server/model/common/request"
	"ygpt/server/model/common/response"
	"ygpt/server/model/yunguan/system"
	"ygpt/server/utils"
)

type SystemSoftwareApi struct {
}

// 添加系统运行工具
func (d *SystemSoftwareApi) AddSystemSoftware(c *gin.Context) {
	var r system.SystemSoftware
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(r, utils.RegisterVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	serviceReturn, err := systemSoftwareService.AddSystemSoftwareService(r)
	if err != nil {
		global.GVA_LOG.Error("添加失败!", zap.Error(err))
		response.FailWithDetailed(global.GVA_MODEL{ID: serviceReturn.ID}, "添加失败", c)
		return
	}
	response.OkWithMessage("添加成功", c)
}

// 修改系统运行工具
func (d *SystemSoftwareApi) UpdateSystemSoftware(c *gin.Context) {
	var model system.SystemSoftware
	err := c.ShouldBindJSON(&model)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = systemSoftwareService.UpdateSystemSoftwareService(model)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithDetailed(gin.H{"updateAt": time.Now().Unix()}, "修改成功", c)
}

// 获取系统运行工具列表
func (b *SystemSoftwareApi) GetSystemSoftwareList(c *gin.Context) {
	var r request.PageInfo
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(r, utils.RegisterVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := systemSoftwareService.GetSystemSoftwareListService(r)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     r.Page,
		PageSize: r.PageSize,
	}, "获取成功", c)
}

// 删除系统运行工具
func (b *SystemSoftwareApi) DeleteSystemSoftware(c *gin.Context) {
	var model system.SystemSoftware
	err := c.ShouldBindJSON(&model)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = systemSoftwareService.DeleteSystemSoftwareService(model.ID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithDetailed(gin.H{"deleteAt": time.Now().Unix()}, "删除成功", c)
}

func (b *SystemSoftwareApi) StartScan(c *gin.Context) {
	// 从 URL 参数获取扫描目录，默认为 /home
	directoryToScan := c.Query("url")
	sshUser := c.Query("user")
	sshHost := c.Query("host")
	sshPassword := c.Query("password")
	filePath := fmt.Sprintf("/home/%s/clamav_setup.sh", sshUser)
	// 创建 SSH 配置
	config := &ssh.ClientConfig{
		User: sshUser,
		Auth: []ssh.AuthMethod{
			ssh.Password(sshPassword), // 使用密码进行身份验证
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // 忽略主机密钥校验
		Timeout:         30 * time.Second,            // 设置超时
	}

	// 连接到远程服务器
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", sshHost), config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to connect to host: %s", err.Error()),
		})
		return
	}
	defer client.Close()

	createFlie := fmt.Sprintf("touch %s", filePath)
	createFileCommand := `
echo "#!/bin/bash
echo "传递的扫描路径：$1"
echo "扫描路径: $SCAN_PATH"
echo "创建并覆盖 clamav_setup.sh 文件..."

SCAN_PATH=$1

# 更新系统软件包
echo \"更新系统软件包...\"
sudo apt update && sudo apt upgrade -y  

# 安装 ClamAV 和 ClamAV-Daemon
echo \"安装 ClamAV 和 ClamAV-Daemon...\"
sudo apt install clamav clamav-daemon -y

# 停止 ClamAV-Freshclam 服务，以便更新病毒数据库
echo \"停止 ClamAV-Freshclam 服务以更新病毒数据库...\"
sudo systemctl stop clamav-freshclam

# 更新病毒数据库
echo \"更新病毒数据库...\"
sudo freshclam

# 重新启动 ClamAV-Freshclam 服务
echo \"重新启动 ClamAV-Freshclam 服务...\"
sudo systemctl start clamav-freshclam

# 扫描指定目录中的病毒
echo \"扫描 $SCAN_PATH 目录中的病毒...\"
sudo clamscan -r \"$SCAN_PATH\"

# 扫描并自动删除指定目录中的感染文件
echo \"扫描并删除 $SCAN_PATH 目录中的感染文件...\"
sudo clamscan -r --remove \"$SCAN_PATH\"

# 将扫描结果保存到日志文件中
echo \"保存扫描结果到 scan_log.txt 文件中...\"
sudo clamscan -r \"$SCAN_PATH\" | tee scan_log.txt

# 自动化扫描设置：添加到 Cron 作业
echo \"设置每天凌晨 2 点自动扫描...\"
cron_job=\"0 2 * * * /usr/bin/clamscan -r $SCAN_PATH --log=/var/log/clamav/scan.log\"
(crontab -l 2>/dev/null; echo \"$cron_job\") | crontab -

# 配置 ClamAV-Daemon 实现实时扫描
echo \"配置 ClamAV-Daemon 进行实时扫描...\"
sudo sed -i '/^Example/d' /etc/clamav/clamd.conf
sudo sed -i 's|#LocalSocket .*|LocalSocket /var/run/clamav/clamd.ctl|' /etc/clamav/clamd.conf

# 重新启动 ClamAV-Daemon 服务
echo \"重新启动 ClamAV-Daemon 服务...\"
sudo systemctl restart clamav-daemon

echo \"ClamAV 安装、病毒扫描和定时任务配置完成。\"

echo \"提示：可通过安装 Xming 和 Putty 实现图形化界面访问。\"
" > %s
`
	// 创建新的会话执行文件创建命令
	session, err := client.NewSession()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to create session: %s", err.Error()),
		})
		return
	}
	defer session.Close()
	err = session.Run(createFlie)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Failed to create  clamav_setup.sh file",
			"output": err.Error(),
		})
		return
	}
	// 创建新的会话执行文件写入命令
	session, err = client.NewSession()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to create session: %s", err.Error()),
		})
		return
	}
	// 使用 session 执行命令
	createFileCommand = fmt.Sprintf(createFileCommand, filePath)

	// 执行命令：将文件内容输出到目标文件
	var stdout, stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	err = session.Run(createFileCommand)
	if err != nil {
		log.Fatalf("Failed to create and write file: %v", err)
	}
	// 创建新的会话执行文件权限设置命令
	session, err = client.NewSession()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to create session for chmod: %s", err.Error()),
		})
		return
	}
	defer session.Close()

	// 设置文件为可执行
	setExecutableCommand := fmt.Sprintf("chmod +x %s", filePath)
	err = session.Run(setExecutableCommand)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Failed to set file executable",
			"output": err.Error(),
		})
		return
	}

	// 创建新的会话执行扫描命令
	session, err = client.NewSession()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to create session for scan: %s", err.Error()),
		})
		return
	}
	defer session.Close()

	// 执行脚本
	// 假设 directoryToScan 是你从请求中获取到的扫描路径，filePath 是脚本路径
	sshCommand := fmt.Sprintf("clamscan  -r %s", directoryToScan)
	//sshCommand := fmt.Sprintf("ls  %s", directoryToScan)

	fmt.Println("SSH Command:++++++++++++++++++++++++++++++++++++++++++", sshCommand)
	fmt.Printf("传递给脚本的扫描路径-----------------------:%s\n", directoryToScan)

	session, err = client.NewSession()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to create SSH session: %s", err.Error()),
		})
		return
	}
	defer session.Close()

	session.Stdout = &stdout
	session.Stderr = &stderr
	err = session.Run(sshCommand)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  fmt.Sprintf("Scan failed: %s", err.Error()),
			"output": stderr.String(),
		})
		return
	}

	// 返回扫描结果
	c.JSON(http.StatusOK, gin.H{
		"message": "Scan started successfully",
		"output":  stdout.String(),
	})
}

func (b *SystemSoftwareApi) StopScan(c *gin.Context) {
	// 从 URL 参数获取远程服务器连接信息
	sshUser := c.Query("user")
	sshHost := c.Query("host")
	sshPassword := c.Query("password")

	// 创建 SSH 配置
	config := &ssh.ClientConfig{
		User: sshUser,
		Auth: []ssh.AuthMethod{
			ssh.Password(sshPassword), // 使用密码进行身份验证
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // 忽略主机密钥校验
		Timeout:         30 * time.Second,            // 设置超时
	}

	// 连接到远程服务器
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", sshHost), config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to connect to host: %s", err.Error()),
		})
		return
	}
	defer client.Close()

	// 查找并停止 clamscan 进程
	session, err := client.NewSession()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to create session: %s", err.Error()),
		})
		return
	}
	defer session.Close()

	// 查找 clamscan 进程的 PID
	psCommand := "pgrep clamscan" // 使用 pgrep 查找 clamscan 进程
	var stdout, stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr
	err = session.Run(psCommand)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to find clamscan process: %s", stderr.String()),
		})
		return
	}

	// 获取 PID 列表
	pids := strings.Split(stdout.String(), "\n")
	if len(pids) == 0 || pids[0] == "" {
		c.JSON(http.StatusOK, gin.H{
			"message": "No scanning process found to stop",
		})
		return
	}

	// 杀掉所有正在运行的 clamscan 进程
	for _, pid := range pids {
		if pid != "" {
			killCommand := fmt.Sprintf("kill -9 %s", pid)
			session, err := client.NewSession()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": fmt.Sprintf("Failed to create session for killing process: %s", err.Error()),
				})
				return
			}
			defer session.Close()
			err = session.Run(killCommand)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": fmt.Sprintf("Failed to kill process with PID %s: %s", pid, err.Error()),
				})
				return
			}
		}
	}

	// 返回停止扫描成功的响应
	c.JSON(http.StatusOK, gin.H{
		"message": "Scan stopped successfully",
	})
}

func (b *SystemSoftwareApi) UninstallClamAV(c *gin.Context) {
	sshUser := c.Query("user")
	sshHost := c.Query("host")
	sshPassword := c.Query("password")

	config := &ssh.ClientConfig{
		User: sshUser,
		Auth: []ssh.AuthMethod{
			ssh.Password(sshPassword),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         30 * time.Second,
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", sshHost), config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to connect to host: %s", err.Error()),
		})
		return
	}
	defer client.Close()

	// Step 1: Check if ClamAV services are running
	checkServiceCommand := "systemctl is-active clamav-daemon"
	session, err := client.NewSession()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to create session: %s", err.Error()),
		})
		return
	}
	defer session.Close()

	var stdout, stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	// 检查服务状态
	err = session.Run(checkServiceCommand)
	if err == nil && stdout.String() == "active\n" {
		// 如果服务正在运行，调用 stopScan API
		b.StopScan(c) // 假设 StopScan 已实现并且停止了 ClamAV 服务
	}

	// Step 2: Uninstall ClamAV
	uninstallCommand := "apt-get purge -y clamav clamav-daemon"
	session, err = client.NewSession()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to create session for uninstall: %s", err.Error()),
		})
		return
	}
	defer session.Close()

	session.Stdout = &stdout
	session.Stderr = &stderr

	err = session.Run(uninstallCommand)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Failed to uninstall ClamAV",
			"output": stderr.String(),
		})
		return
	}

	// 返回卸载成功信息
	c.JSON(http.StatusOK, gin.H{
		"message": "ClamAV uninstalled successfully",
		"output":  stdout.String(),
	})
}

// 检测系统资源信息，支持执行自定义 Shell 脚本或平台默认资源检测
func (b *SystemSoftwareApi) SystemCheck(c *gin.Context) {
	// 获取用户自定义 Shell 脚本路径（通过请求参数传递）
	scriptPath := c.Query("script_path")

	if scriptPath != "" {
		// 如果用户提供了自定义 Shell 脚本路径，则执行该脚本
		cmd := exec.Command("/bin/sh", scriptPath)
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":  "Failed to execute custom script",
				"detail": err.Error(),
			})
			return
		}

		// 返回脚本的输出结果
		c.JSON(http.StatusOK, gin.H{
			"result": out.String(),
		})
		return
	}

	// 否则，执行平台默认的资源检测（使用 gopsutil 获取系统信息）
	cpuUsage, err := cpu.Percent(1*time.Second, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get CPU usage"})
		return
	}

	memory, err := mem.VirtualMemory()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get memory info"})
		return
	}

	diskUsage, err := disk.Usage("/")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get disk usage"})
		return
	}

	netInterfaces, err := net.Interfaces()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get network info"})
		return
	}

	systemInfo := gin.H{
		"cpu_usage":          cpuUsage,
		"memory_total":       memory.Total,
		"memory_used":        memory.Used,
		"memory_used_pct":    memory.UsedPercent,
		"disk_total":         diskUsage.Total,
		"disk_free":          diskUsage.Free,
		"disk_used_pct":      diskUsage.UsedPercent,
		"network_interfaces": netInterfaces,
	}

	c.JSON(http.StatusOK, systemInfo)
}
