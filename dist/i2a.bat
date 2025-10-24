@echo off
setlocal enabledelayedexpansion

:: 获取当前脚本所在目录
set "SCRIPT_DIR=%~dp0"
set "EXE_PATH=%SCRIPT_DIR%image2avif.exe"

:: 处理拖放文件的情况（如果有）
set "ARGUMENTS=%*"
if "%ARGUMENTS%"=="" (
    :: 没有参数时显示帮助
    "%EXE_PATH%" --help
    pause
    exit /b 0
)

:: 执行主程序
"%EXE_PATH%" %*

:: 保持窗口打开（仅当通过双击运行时）
if "%~1"=="" (
    pause
)

endlocal
