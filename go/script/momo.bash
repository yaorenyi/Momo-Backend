#!/bin/bash

# ==========================================================
# Momo Backend 交互式控制脚本
# ==========================================================

APP_NAME="momo-backend-linux-amd64"
APP_PATH="./${APP_NAME}"
LOG_FILE="./momo.log"
PID_FILE="./momo.pid"

# 权限与文件检查
if [ ! -f "$APP_PATH" ]; then
    echo "错误：程序文件 $APP_PATH 不存在！"
    echo "请确保脚本与 $APP_NAME 放在同一目录下。"
    exit 1
fi

if [ ! -x "$APP_PATH" ]; then
    echo "正在尝试为程序添加执行权限..."
    chmod +x "$APP_PATH"
fi

# 环境变量设置 (根据 momo 的需要可自行调整)
export GIN_MODE=release

# --- 功能函数 ---

start() {
    if [ -f "$PID_FILE" ]; then
        PID=$(cat "$PID_FILE")
        if kill -0 "$PID" 2>/dev/null; then
            echo "Momo 已在运行中，PID: $PID"
            return 0
        else
            echo "发现过期的 PID 文件，正在清理..."
            rm -f "$PID_FILE"
        fi
    fi

    echo "正在启动 $APP_NAME ..."
    # 使用 nohup 后台运行，并将输出重定向到日志
    nohup "$APP_PATH" > "$LOG_FILE" 2>&1 &
    
    # 记录 PID
    echo $! > "$PID_FILE"
    
    # 稍等片刻检查是否启动成功
    sleep 1
    if kill -0 $(cat "$PID_FILE") 2>/dev/null; then
        echo "程序已启动，PID: $(cat $PID_FILE)"
        echo "日志文件: $LOG_FILE"
    else
        echo "启动失败，请检查日志 $LOG_FILE"
    fi
}

stop() {
    if [ ! -f "$PID_FILE" ]; then
        echo "程序未运行（找不到 PID 文件）"
        # 兜底检查：防止 PID 文件丢失但进程还在
        PIDS=$(pgrep -f "$APP_NAME")
        if [ -n "$PIDS" ]; then
            echo "发现匹配进程 ($PIDS)，请手动 kill 或清理。"
        fi
        return 1
    fi

    PID=$(cat "$PID_FILE")
    if kill -0 "$PID" 2>/dev/null; then
        echo "正在停止 Momo (PID: $PID)..."
        kill "$PID"
        
        # 优雅停机等待
        for i in {1..10}; do
            if ! kill -0 "$PID" 2>/dev/null; then
                break
            fi
            echo -n "."
            sleep 1
        done

        # 如果还没停，暴力强制结束
        if kill -0 "$PID" 2>/dev/null; then
            echo -e "\n进程无响应，强制终止..."
            kill -9 "$PID"
        fi

        rm -f "$PID_FILE"
        echo -e "\n程序已安全停止"
    else
        echo "检测到残留 PID 文件，已清理。"
        rm -f "$PID_FILE"
    fi
}

status() {
    if [ -f "$PID_FILE" ]; then
        PID=$(cat "$PID_FILE")
        if kill -0 "$PID" 2>/dev/null; then
            echo "状态：运行中"
            echo "PID：$PID"
            echo "运行时间：$(ps -o etime= -p "$PID")"
        else
            echo "状态：进程异常（PID 文件存在但进程已消失）"
        fi
    else
        echo "状态：未运行"
    fi
}

show_menu() {
    echo
    echo "=============================="
    echo "    Momo Backend 管理面板"
    echo "=============================="
    echo "1) 启动服务"
    echo "2) 停止服务"
    echo "3) 重启服务"
    echo "4) 查看状态"
    echo "5) 查看实时日志 (Ctrl+C 退出日志)"
    echo "6) 退出面板"
    echo "------------------------------"
    read -p "请输入操作编号 [1-6]: " choice
}

# --- 主循环 ---

while true; do
    show_menu
    case "$choice" in
        1) start ;;
        2) stop ;;
        3) 
            stop
            sleep 1
            start 
            ;;
        4) status ;;
        5) 
            if [ -f "$LOG_FILE" ]; then
                tail -f "$LOG_FILE"
            else
                echo "日志文件尚未生成"
            fi
            ;;
        6) 
            echo "退出管理面板"
            exit 0 
            ;;
        *) 
            echo "无效选项，请输入 1-6" 
            ;;
    esac
    echo
    read -p "按回车键继续..." temp
done