#!/bin/bash

# 设置项目根目录
PROJECT_ROOT=$(pwd)

# 定义打包目录
PUBLIC_DIR="/home/suanliquan/sl-yun"
PACKAGE_DIR="/home/package/slyun"

echo "清理上次残留"
cd $PACKAGE_DIR && rm -rf *

echo "开始拉取最新代码"
cd $PROJECT_ROOT && git pull

echo "清除残留"
# 删除上一次打包的文件夹
rm -rf $PUBLIC_DIR
mkdir -p $PUBLIC_DIR

# 前端构建命令
echo "开始打包 前端"
cd $PROJECT_ROOT/web && npm install && npm run build

# 检查 dist 目录是否存在
if [ -d "dist" ]; then
    mv dist $PUBLIC_DIR/web
else
    echo "Error: 'dist' directory not found after build."
    exit 1
fi

# 脚本
cd $PROJECT_ROOT && cp install.sh unstall.sh config.yaml.template $PUBLIC_DIR/

# 后端构建命令
echo "开始打包 后端"
cd $PROJECT_ROOT/server && go build -o cloud_end main.go

# 检查 cloud_end 文件是否存在
if [ -f "cloud_end" ]; then
    mv cloud_end $PUBLIC_DIR/server
    cp -r $PROJECT_ROOT/server/resource $PUBLIC_DIR/
else
    echo "Error: 'cloud_end' binary not found after build."
    exit 1
fi

# git 提交
cd $PROJECT_ROOT
DATE_TODAY=$(date +"%Y-%m-%d")
TAR_FILENAME=yun$DATE_TODAY.tar
echo "开始归档"
tar -cvf $TAR_FILENAME  $PUBLIC_DIR
mv $TAR_FILENAME $PACKAGE_DIR
cd $PACKAGE_DIR && git add . && git commit -m "update[$DATE_TODAY]" && git push
echo "归档完成"
echo "流程结束"