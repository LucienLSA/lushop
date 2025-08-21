from app import create_app
import os
import signal
import sys

app = create_app()

def signal_handler(sig, frame):
    """优雅退出处理"""
    app.logger.info("Server is shutting down...")
    sys.exit(0)

if __name__ == "__main__":
    # 注册信号处理器
    signal.signal(signal.SIGINT, signal_handler)
    signal.signal(signal.SIGTERM, signal_handler)

    # 启动服务
    app.run(
        host=app.config["HOST"],
        port=app.config["PORT"],
        debug=app.config["DEBUG"]
    )