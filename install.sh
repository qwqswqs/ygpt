#!/bin/bash

set -e  # 发生错误时立即退出

### 配置参数 ###
SERVICE_NAME="cloud"          # 服务名称
EXEC_FILE_NAME="server"       # 可执行文件名（安装目录下）
INSTALL_DIR="/usr/local/$SERVICE_NAME"  # 安装目录
SERVICE_FILE="/etc/systemd/system/$SERVICE_NAME.service"

echo ">>> 创建目录 $INSTALL_DIR"
mkdir -p "$INSTALL_DIR"

echo ">>> 复制文件到 $INSTALL_DIR"
cp -r ./* "$INSTALL_DIR/"

echo ">>> 赋予执行权限"
chmod +x "$INSTALL_DIR/$EXEC_FILE_NAME"

echo ">>> 创建 systemd 服务文件"
cat <<EOF >"$SERVICE_FILE"
[Unit]
Description=$SERVICE_NAME Service
After=network.target

[Service]
Type=simple
ExecStart=$INSTALL_DIR/$EXEC_FILE_NAME
WorkingDirectory=$INSTALL_DIR
Restart=always
RestartSec=5s

[Install]
WantedBy=multi-user.target
EOF

echo ">>> 重新加载 systemd 配置"
systemctl daemon-reload

echo ">>> 设置开机自启"
systemctl enable "$SERVICE_NAME"

echo ">>> 创建完成"
echo ">>> 使用 'systemctl start $SERVICE_NAME' 启动服务"
echo ">>> 使用 'systemctl status $SERVICE_NAME' 查看服务状态"
echo ">>> 使用 'systemctl stop $SERVICE_NAME' 停止服务"
