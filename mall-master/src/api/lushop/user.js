import request from '@/utils/request'

// 用户登录
export function login(data) {
  return request({
    url: 'user/pwd_login',
    method: 'post',
    data: data
  })
}

// 用户注册
export function register(data) {
  return request({
    url: 'user/register',
    method: 'post',
    data: data
  })
}

// 获取用户列表
export function fetchUserList(params) {
  return request({
    url: 'user/list',
    method: 'get',
    params: params
  })
}

// 获取用户详情
export function getUserDetail() {
  return request({
    url: 'user/detail',
    method: 'get'
  })
}

// 更新用户信息
export function updateUser(data) {
  return request({
    url: 'user/update',
    method: 'patch',
    data: data
  })
}

// 刷新token
export function refreshToken() {
  return request({
    url: 'user/refresh',
    method: 'get'
  })
}

// 获取用户地址列表
export function getAddressList() {
  return request({
    url: 'address',
    method: 'get'
  })
}

// 创建地址
export function createAddress(data) {
  return request({
    url: 'address',
    method: 'post',
    data: data
  })
}

// 更新地址
export function updateAddress(id, data) {
  return request({
    url: `address/${id}`,
    method: 'put',
    data: data
  })
}

// 删除地址
export function deleteAddress(id) {
  return request({
    url: `address/${id}`,
    method: 'delete'
  })
}

// 获取轮播图列表
export function getBannerList() {
  return request({
    url: 'banners',
    method: 'get'
  })
}

// 创建轮播图
export function createBanner(data) {
  return request({
    url: 'banners',
    method: 'post',
    data: data
  })
}

// 更新轮播图
export function updateBanner(id, data) {
  return request({
    url: `banners/${id}`,
    method: 'put',
    data: data
  })
}

// 删除轮播图
export function deleteBanner(id) {
  return request({
    url: `banners/${id}`,
    method: 'delete'
  })
} 