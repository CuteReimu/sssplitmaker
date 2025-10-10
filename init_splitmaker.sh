#!/bin/bash
set -e

# 删除并新建 hk-split-maker 目录
rm -rf hk-split-maker
mkdir hk-split-maker
cd hk-split-maker

# 初始化 git 仓库并设置远程
git init
git remote add origin git@github.com:slaurent22/hk-split-maker.git
git config core.sparseCheckout true
git sparse-checkout set --no-cone

# 配置稀疏检出规则
echo "src/asset/silksong/categories/*" > .git/info/sparse-checkout
echo "!src/asset/silksong/categories/category-directory.json" >> .git/info/sparse-checkout
# 拉取主分支
git pull --depth=1 origin main

cd ..

# 删除并复制 splitmaker 目录
rm -rf splitmaker
cp -avf hk-split-maker/src/asset/silksong/categories splitmaker