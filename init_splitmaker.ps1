# 删除并新建 hk-split-maker 目录
Remove-Item -Recurse -Force hk-split-maker -ErrorAction SilentlyContinue
New-Item -ItemType Directory -Path hk-split-maker | Out-Null
Set-Location hk-split-maker

# 初始化 git 仓库并设置远程
git init
git remote add origin git@github.com:slaurent22/hk-split-maker.git
git config core.sparseCheckout true
git sparse-checkout set --no-cone

# 配置稀疏检出规则
Set-Content .git/info/sparse-checkout "src/asset/silksong/categories/*"
Add-Content .git/info/sparse-checkout "!src/asset/silksong/categories/category-directory.json"
Add-Content .git/info/sparse-checkout "!src/asset/silksong/categories/every.json"
Add-Content .git/info/sparse-checkout "!src/asset/silksong/categories/room-timer.json"
Add-Content .git/info/sparse-checkout "!src/asset/silksong/categories/blank.json"

# 拉取主分支
git pull --depth=1 origin main

Set-Location ..

# 删除并复制 splitmaker 目录
Remove-Item -Recurse -Force splitmaker -ErrorAction SilentlyContinue
New-Item -ItemType Directory -Path splitmaker | Out-Null
Copy-Item -Force hk-split-maker/src/asset/silksong/categories/*.json splitmaker

git add -A splitmaker