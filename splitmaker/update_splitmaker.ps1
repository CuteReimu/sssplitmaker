$response = Invoke-WebRequest -Uri "https://api.github.com/repos/slaurent22/hk-split-maker/contents/src/asset/silksong/categories"
$objects = $response | ConvertFrom-Json
$objects | ForEach-Object { Invoke-WebRequest -Uri "https://raw.githubusercontent.com/slaurent22/hk-split-maker/refs/heads/main/src/asset/silksong/categories/$($_.name)" -OutFile $_.name -v }