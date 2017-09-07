# Whalebrew for Windows

$BINARY_URL = "https://github.com/bfirsh/whalebrew/releases/download/0.0.5/whalebrew-Windows-x86_64.exe"
$USER_PATH = [Environment]::GetEnvironmentVariable("Path", [System.EnvironmentVariableTarget]::User)
$WHALEBREW_INSTALL_PATH = [Environment]::GetEnvironmentVariable("WHALEBREW_INSTALL_PATH", [System.EnvironmentVariableTarget]::User)

if ($WHALEBREW_INSTALL_PATH -eq $null -or $WHALEBREW_INSTALL_PATH -eq '') {
  $WHALEBREW_INSTALL_PATH = 'C:\whalebrew';
  [Environment]::SetEnvironmentVariable("WHALEBREW_INSTALL_PATH", $WHALEBREW_INSTALL_PATH, [System.EnvironmentVariableTarget]::User)
}

if (![System.IO.Directory]::Exists($WHALEBREW_INSTALL_PATH)) {
  [System.IO.Directory]::CreateDirectory($WHALEBREW_INSTALL_PATH)
}

if ($($USER_PATH).ToLower().Contains($($WHALEBREW_INSTALL_PATH).ToLower()) -eq $false) {
  [Environment]::SetEnvironmentVariable("Path", "$USER_PATH;%WHALEBREW_INSTALL_PATH%", [System.EnvironmentVariableTarget]::User)
}

$WHALEBREW_PATH = Join-Path $WHALEBREW_INSTALL_PATH "whalebrew.exe"
(New-Object System.Net.WebClient).DownloadFile($BINARY_URL, "$WHALEBREW_PATH")

Write-Output "Installed whalebrew to `"$WHALEBREW_PATH`""
