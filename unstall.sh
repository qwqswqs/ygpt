#!/bin/bash

set -e  # 发生错误时立即退出

### 配置参数 ###
SERVICE_NAME="cloud"          # 仅需修改此处，即可卸载不同的服务
INSTALL_DIR="/usr/local/$SERVICE_NAME"  # 安装目录
SERVICE_FILE="/etc/systemd/system/$SERVICE_NAME.service"

echo ">>> 停止并禁用服务 $SERVICE_NAME"
if systemctl is-active --quiet "$SERVICE_NAME"; then
    systemctl stop "$SERVICE_NAME"
fi

if systemctl is-enabled --quiet "$SERVICE_NAME"; then
    systemctl disable "$SERVICE_NAME"
fi

echo ">>> 删除 systemd 服务文件 $SERVICE_FILE"
rm -f "$SERVICE_FILE"

echo ">>> 重新加载 systemd 配置"
systemctl daemon-reload
systemctl reset-failed

echo ">>> 删除安装目录 $INSTALL_DIR"
rm -rf "$INSTALL_DIR"

echo ">>> 卸载完成"
echo ">>> 你可以运行 'systemctl list-units --failed' 检查是否有残留的失败服务"
