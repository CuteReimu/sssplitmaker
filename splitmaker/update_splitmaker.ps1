$response = curl -s "https://api.github.com/repos/slaurent22/hk-split-maker/contents/src/asset/silksong/categories"
$objects = $response | ConvertFrom-Json
$objects | ForEach-Object { $_.name }