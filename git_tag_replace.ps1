# GitTagReplace.ps1
# 功能：根据输入参数，删除旧标签并创建新标签（支持标签更新模式）
param(
    [string]$InputTags
)

# 设置执行策略允许继续执行（主要影响非终止性错误）
$ErrorActionPreference = "Continue"

Write-Host "=== Git 标签替换脚本 ===" -ForegroundColor Cyan
Write-Host "操作模式说明：" -ForegroundColor Yellow
Write-Host "• 输入一个标签名: 将该标签重新创建（删除旧标签后新建）"
Write-Host "• 输入两个标签名（用空格分隔）: 删除第一个标签，创建第二个标签"
Write-Host ""

# 处理输入
if ([string]::IsNullOrWhiteSpace($InputTags)) {
    $InputTags = Read-Host "请输入标签（格式: 要删除的标签 要创建的标签 或 单个标签名）"
}

if ([string]::IsNullOrWhiteSpace($InputTags)) {
    Write-Host "错误：输入不能为空。" -ForegroundColor Red
    exit 1
}

# 使用更智能的分割方法，处理包含斜杠的标签名[2,4](@ref)
# 使用正则表达式分割，确保正确处理包含特殊字符的标签
$TagArray = @()
if ($InputTags -match '^\S+$') {
    # 单个标签（不包含空格）
    $TagArray = @($InputTags.Trim())
} else {
    # 多个标签，使用智能分割[2](@ref)
    $parts = $InputTags.Trim() -split '\s+', 2
    foreach ($part in $parts) {
        if (-not [string]::IsNullOrWhiteSpace($part)) {
            $TagArray += $part
        }
    }
}

# 根据输入数量确定操作模式
switch ($TagArray.Count) {
    1 {
        $DeleteTag = $TagArray[0]
        $CreateTag = $TagArray[0]
        Write-Host "检测到单个标签输入，执行标签更新模式" -ForegroundColor Green
        Write-Host "将重新创建标签: $DeleteTag" -ForegroundColor Green
    }
    2 {
        $DeleteTag = $TagArray[0]
        $CreateTag = $TagArray[1]
        Write-Host "检测到两个标签输入，执行标签替换模式" -ForegroundColor Green
        Write-Host "将删除标签: $DeleteTag，创建标签: $CreateTag" -ForegroundColor Green
    }
    default {
        Write-Host "错误：请输入1个或2个标签名（用空格分隔）。" -ForegroundColor Red
        exit 1
    }
}

Write-Host ""
Write-Host "待执行操作详情：" -ForegroundColor Cyan
Write-Host "  - 删除标签: $DeleteTag" -ForegroundColor Red
Write-Host "  - 创建标签: $CreateTag" -ForegroundColor Green
Write-Host ""

# 确认操作
$Confirm = Read-Host "确定要执行以上操作吗? (y/N)"
if ($Confirm -ne 'y' -and $Confirm -ne 'Y') {
    Write-Host "操作已取消。" -ForegroundColor Yellow
    exit 0
}

Write-Host "开始执行操作..." -ForegroundColor Cyan
Write-Host ""

# 记录操作结果
$Operations = @()

# 1. 删除本地标签（使用引号确保路径正确传递）[1,2](@ref)
try {
    Write-Host "[1/4] 删除本地标签 '$DeleteTag'..." -NoNewline
    # 使用引号包裹标签名，防止路径解析错误[1](@ref)
    git tag -d "$DeleteTag" 2>$null
    if ($LASTEXITCODE -eq 0) {
        Write-Host " 成功" -ForegroundColor Green
        $Operations += "删除本地标签: 成功"
    } else {
        Write-Host " 标签不存在" -ForegroundColor DarkYellow
        $Operations += "删除本地标签: 标签不存在"
    }
} catch {
    Write-Host " 错误: $($_.Exception.Message)" -ForegroundColor Red
    $Operations += "删除本地标签: 失败"
}

# 2. 删除远程标签（同样使用引号）
try {
    Write-Host "[2/4] 删除远程标签 '$DeleteTag'..." -NoNewline
    git push --delete origin "$DeleteTag" 2>$null
    if ($LASTEXITCODE -eq 0) {
        Write-Host " 成功" -ForegroundColor Green
        $Operations += "删除远程标签: 成功"
    } else {
        Write-Host " 标签不存在" -ForegroundColor DarkYellow
        $Operations += "删除远程标签: 标签不存在"
    }
} catch {
    Write-Host " 错误: $($_.Exception.Message)" -ForegroundColor Red
    $Operations += "删除远程标签: 失败"
}

# 3. 创建新标签（使用引号）
try {
    Write-Host "[3/4] 创建新标签 '$CreateTag'..." -NoNewline
    git tag -a "$CreateTag" -m "upgrade version" 2>$null
    if ($LASTEXITCODE -eq 0) {
        Write-Host " 成功" -ForegroundColor Green
        $Operations += "创建新标签: 成功"
    } else {
        Write-Host " 失败" -ForegroundColor Red
        $Operations += "创建新标签: 失败"
        exit 1
    }
} catch {
    Write-Host " 错误: $($_.Exception.Message)" -ForegroundColor Red
    $Operations += "创建新标签: 失败"
    exit 1
}

# 4. 推送新标签到远程（使用引号）
try {
    Write-Host "[4/4] 推送新标签到远程仓库..." -NoNewline
    git push origin "$CreateTag" 2>$null
    if ($LASTEXITCODE -eq 0) {
        Write-Host " 成功" -ForegroundColor Green
        $Operations += "推送新标签: 成功"
    } else {
        Write-Host " 失败" -ForegroundColor Red
        $Operations += "推送新标签: 失败"
        exit 1
    }
} catch {
    Write-Host " 错误: $($_.Exception.Message)" -ForegroundColor Red
    $Operations += "推送新标签: 失败"
    exit 1
}

# 显示操作结果摘要
Write-Host ""
Write-Host "=== 操作完成 ===" -ForegroundColor Cyan
Write-Host "操作摘要:" -ForegroundColor Yellow
$Operations | ForEach-Object { Write-Host "  ✓ $_" -ForegroundColor Green }

Write-Host ""
Write-Host "标签 '$CreateTag' 已成功创建并推送到远程仓库。" -ForegroundColor Green

# 验证新标签
Write-Host ""
Write-Host "验证新标签状态:" -ForegroundColor Cyan
git show "$CreateTag" --quiet 2>$null
if ($LASTEXITCODE -eq 0) {
    Write-Host "  ✓ 新标签验证成功" -ForegroundColor Green
} else {
    Write-Host "  ⚠ 新标签验证失败" -ForegroundColor Red
}