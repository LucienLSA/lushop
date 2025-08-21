import request from '@/utils/request'

// 品牌列表
export function fetchList(params) {
  return request({
    url: 'brands',
    method: 'get',
    params: params
  })
}

// 创建品牌
export function createBrand(data) {
  return request({
    url: 'brands',
    method: 'post',
    data: data
  })
}

// 更新品牌
export function updateBrand(id, data) {
  return request({
    url: `brands/${id}`,
    method: 'put',
    data: data
  })
}

// 删除品牌
export function deleteBrand(id) {
  return request({
    url: `brands/${id}`,
    method: 'delete'
  })
}

// 获取品牌详情
export function getBrand(id) {
  return request({
    url: `brands/${id}`,
    method: 'get'
  })
} 