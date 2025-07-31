import request from '@/utils/request'

// 分类列表
export function fetchList(params) {
  return request({
    url: 'categorys',
    method: 'get',
    params: params
  })
}

// 创建分类
export function createCategory(data) {
  return request({
    url: 'categorys',
    method: 'post',
    data: data
  })
}

// 更新分类
export function updateCategory(id, data) {
  return request({
    url: `categorys/${id}`,
    method: 'put',
    data: data
  })
}

// 删除分类
export function deleteCategory(id) {
  return request({
    url: `categorys/${id}`,
    method: 'delete'
  })
}

// 获取分类详情
export function getCategory(id) {
  return request({
    url: `categorys/${id}`,
    method: 'get'
  })
}

// 获取分类品牌列表
export function getCategoryBrands(id) {
  return request({
    url: `categorybrands/${id}`,
    method: 'get'
  })
} 