from flask import Blueprint, request, jsonify
from flask_jwt_extended import create_access_token, jwt_required, get_jwt_identity
from pydantic import BaseModel, EmailStr, field_validator
from app.services.user_service import UserService  # 微服务调用层
from app.utils.exceptions import BusinessError  # 自定义异常
import datetime

user_bp = Blueprint('user', __name__)

# 数据验证模型（Pydantic）
class RegisterRequest(BaseModel):
    mobile: str
    password: str
    nickname: str

    @field_validator('mobile')
    def mobile_validate(cls, v):
        if not v.isdigit() or len(v) != 11:
            raise ValueError('手机号格式错误')
        return v

class LoginRequest(BaseModel):
    mobile: str
    password: str

@user_bp.route('/register', methods=['POST'])
def register():
    """用户注册接口"""
    try:
        data = RegisterRequest(**request.json)
        # 调用用户微服务注册逻辑
        user = UserService.register(
            mobile=data.mobile,
            password=data.password,
            nickname=data.nickname
        )
        return jsonify(code=200, message="注册成功", data=user)
    except BusinessError as e:
        return jsonify(code=e.code, message=e.message)
    except Exception as e:
        return jsonify(code=500, message=f"注册失败：{str(e)}")

@user_bp.route('/login', methods=['POST'])
def login():
    """用户登录接口（生成JWT）"""
    try:
        data = LoginRequest(**request.json)
        user = UserService.check_password(data.mobile, data.password)
        # 生成双Token（访问令牌和刷新令牌）
        access_token = create_access_token(
            identity=user['id'],
            additional_claims={"mobile": user['mobile'], "nickname": user['nickname']},
            expires_delta=datetime.timedelta(seconds=app.config['JWT_ACCESS_TOKEN_EXPIRES'].total_seconds())
        )
        refresh_token = create_refresh_token(
            identity=user['id'],
            expires_delta=datetime.timedelta(seconds=app.config['JWT_REFRESH_TOKEN_EXPIRES'].total_seconds())
        )
        return jsonify(code=200, message="登录成功", data={
            "id": user['id'],
            "nickname": user['nickname'],
            "access_token": access_token,
            "refresh_token": refresh_token
        })
    except BusinessError as e:
        return jsonify(code=e.code, message=e.message)
    
@user_bp.route('/logout', methods=['POST'])
@jwt_required()
def logout():
    """用户登出接口（将Token加入黑名单）"""
    jti = get_jwt()['jti']  # 获取JWT的唯一标识
    redis_client.setex(f"jwt_blacklist:{jti}", app.config['JWT_ACCESS_TOKEN_EXPIRES'], 1)
    return jsonify(code=200, message="注销成功")

@user_bp.route('/refresh', methods=['POST'])
@jwt_required(refresh=True)
def refresh_token():
    """刷新访问令牌"""
    user_id = get_jwt_identity()
    user = UserService.get_user_by_id(user_id)
    access_token = create_access_token(
        identity=user_id,
        additional_claims={"mobile": user['mobile'], "nickname": user['nickname']},
        expires_delta=datetime.timedelta(seconds=app.config['JWT_ACCESS_TOKEN_EXPIRES'].total_seconds())
    )
    return jsonify(code=200, data={"access_token": access_token})

@user_bp.route('/info', methods=['GET'])
@jwt_required()  # JWT认证装饰器
def get_user_info():
    """获取当前登录用户信息"""
    user_id = get_jwt_identity()  # 从JWT中获取用户ID
    user = UserService.get_user_by_id(user_id)
    return jsonify(code=200, data=user)