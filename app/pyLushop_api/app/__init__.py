from flask import Flask
from flask_cors import CORS
from flask_jwt_extended import JWTManager
from flask_redis import FlaskRedis
from config import Config

# 初始化扩展
redis_client = FlaskRedis()
jwt = JWTManager()

def create_app(config_class=Config):
    app = Flask(__name__)
    app.config.from_object(config_class)

    # 初始化扩展
    CORS(app, resources={r"/*": {"origins": "*"}})  # 替代原 CORS 中间件
    redis_client.init_app(app)
    jwt.init_app(app)

    # 初始化日志（替代原 zap 日志）
    init_logger(app)

    # 初始化服务注册（Consul）
    init_consul(app)

    # 注册蓝图（对应原路由初始化）
    register_blueprints(app)

    return app

def init_logger(app):
    """初始化日志系统"""
    import logging
    from logging.handlers import RotatingFileHandler
    import os

    log_config = app.config["LOG"]
    log_dir = log_config["filepath"]
    if not os.path.exists(log_dir):
        os.makedirs(log_dir)

    log_file = os.path.join(log_dir, log_config["filename"])
    handler = RotatingFileHandler(
        log_file,
        maxBytes=log_config["max_size"] * 1024 * 1024,  # 转换为 MB
        backupCount=log_config["max_backups"]
    )
    formatter = logging.Formatter('%(asctime)s %(levelname)s: %(message)s [in %(pathname)s:%(lineno)d]')
    handler.setFormatter(formatter)
    handler.setLevel(log_config["level"].upper())

    app.logger.addHandler(handler)
    app.logger.setLevel(log_config["level"].upper())
    app.logger.info("Logger initialized")

def init_consul(app):
    """初始化 Consul 服务注册"""
    from consul import Consul
    consul_config = app.config["CONSUL"]
    try:
        consul = Consul(host=consul_config["host"], port=consul_config["port"])
        app.consul = consul
        app.logger.info("Consul client initialized")
    except Exception as e:
        app.logger.error(f"Consul initialization failed: {str(e)}")

def register_blueprints(app):
    """注册路由蓝图（router 包）"""
    from app.routes.user import user_bp
    from app.routes.goods import goods_bp
    from app.routes.order import order_bp
    from app.routes.oauth2 import oauth2_bp
    from app.routes.oss import oss_bp

    # 注册蓝图，前缀对应原版本号
    api_prefix = f"/{app.config['VERSION']}"
    app.register_blueprint(user_bp, url_prefix=f"{api_prefix}/user")
    app.register_blueprint(goods_bp, url_prefix=f"{api_prefix}/goods")
    app.register_blueprint(order_bp, url_prefix=f"{api_prefix}/order")
    app.register_blueprint(oauth2_bp, url_prefix="/oauth2")
    app.register_blueprint(oss_bp, url_prefix=f"{api_prefix}/oss")

    # 健康检查路由（对应原 /health）
    @app.route("/health")
    def health_check():
        return {"code": 200, "success": True}