import os
from datetime import timedelta

class Config:
    # 服务器基础配置
    NAME = "v2lushopapi"
    HOST = os.getenv("HOST", "10.101.171.193")  # 优先从环境变量获取
    PORT = int(os.getenv("PORT", 8101))
    VERSION = "v2"
    TAGS = ["flask", "lushop", "api"]
    DEBUG = True  # 对应 debug 模式

    # OAuth2 配置
    OAUTH2_ACCESS_TOKEN_EXP = timedelta(hours=2)
    OAUTH2_REFRESH_TOKEN_EXP = timedelta(hours=12)
    OAUTH2_JWT_SIGNED_KEY = "k2bjI75JJHolp0i"
    OAUTH2_CLIENTS = [
        {
            "id": "15383660176",
            "secret": "123456",
            "name": "测试应用1",
            "domain": "http://10.101.178.90:8101",
            "scope": [{"id": "all", "title": "用户账号、手机、权限、角色等信息"}]
        },
        {
            "id": "test_client_2",
            "secret": "test_secret_2",
            "name": "测试应用2",
            "domain": "http://10.101.178.90:8101",
            "scope": [{"id": "all", "title": "用户账号, 手机, 权限, 角色等信息"}]
        }
    ]

    # Session 配置
    SESSION_COOKIE_NAME = "session_id"
    SECRET_KEY = "kkoiybh1ah6rbh0"  # 用于 session 加密
    PERMANENT_SESSION_LIFETIME = timedelta(seconds=1200)  # 20分钟

    # JWT 配置（对应原 jwt 配置）
    JWT_SECRET_KEY = "s7DRD35xGlStUOFjsjSS4sbqg0azszYg"
    JWT_ACCESS_TOKEN_EXPIRES = timedelta(seconds=2592000)  # 30天
    JWT_REFRESH_TOKEN_EXPIRES = timedelta(seconds=5184000)  # 60天

    # 阿里云短信配置
    ALI_SMS = {
        "api_key": "Ali_ApiKey",
        "api_secret": "Ali_ApiSecret",
        "sign_name": "阿里云短信测试",
        "template_code": "SMS_154950909",
        "phone_number": "19821216806",
        "region_id": "cn-hangzhou",
        "expire": 600,
        "cooldown": 60
    }

    # Redis 配置
    REDIS_URL = "redis://:{password}@{host}:{port}/{db}".format(
        host="127.0.0.1",
        port=6379,
        db=4,
        password=""  # 无密码则留空
    )
    REDIS_POOL_SIZE = 10

    # 服务注册与监控
    CONSUL = {"host": "192.168.226.140", "port": 8500}
    JAEGER = {
        "service_name": "goods_web",
        "jaeger_gin_endpoint": "192.168.226.140:4318"
    }
    SENTINEL = {
        "app": {"name": "v2lushop_api", "type": 0},
        "log": {
            "dir": "./temp/csp",
            "pid": False,
            "metric": {"maxFileCount": 14, "flushIntervalSec": 1}
        },
        "stat": {
            "globalStatisticIntervalMsTotal": 6000,
            "system": {"collectIntervalMs": 1000}
        }
    }

    # 日志配置
    LOG = {
        "level": "debug",
        "filepath": "./temp/logs/",
        "filename": "v2lushopapi.log",
        "max_size": 200,
        "max_age": 30,
        "max_backups": 7
    }

    # OSS 配置
    OSS = {
        "api_key": "OSS_ACCESS_KEY_ID",
        "api_secret": "OSS_ACCESS_KEY_SECRET",
        "host": "http://lushop666.oss-cn-shanghai.aliyuncs.com",
        "callback_url": "https://3e69ee01b766.ngrok-free.app/g/v2/oss/callback",
        "upload_dir": "lushop_images/",
        "expired_time": 3000,
        "bucket": "lushop666",
        "endpoint": "oss-cn-shanghai.aliyuncs.com"
    }

    # 支付宝配置
    ALIPAY = {
        "app_id": "AliPay_Id",
        "private_key": "Lushop_Private_Key",
        "ali_public_key": "AliPay_Public_Key",
        "notify_url": "https://3e69ee01b766.ngrok-free.app/g/v2/oss/callback",
        "return_url": "",
        "product_code": "FAST_INSTANT_TRADE_PAY"
    }

    # 微服务配置
    SERVICES = {
        "user_srv": {"name": "user_srv"},
        "userop_srv": {"name": "userop_srv"},
        "goods_srv": {"name": "goods_srv"},
        "order_srv": {"name": "order_srv"},
        "inventory_srv": {"name": "inventory_srv"}
    }