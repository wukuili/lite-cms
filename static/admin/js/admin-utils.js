/**
 * LiteCMS Admin 共用工具脚本
 * 统一处理认证过期和 401 响应
 */

// 检查登录状态
const token = localStorage.getItem('token');
if (!token) window.location.href = '/admin/login';

// 默认请求头
const headers = { 'Content-Type': 'application/json', 'Authorization': 'Bearer ' + token };

/**
 * 封装的 fetch，自动处理 401 认证过期
 */
async function authFetch(url, options = {}) {
    if (!options.headers) {
        options.headers = headers;
    }
    const res = await fetch(url, options);
    if (res.status === 401) {
        localStorage.removeItem('token');
        alert('登录已过期，请重新登录');
        window.location.href = '/admin/login';
        throw new Error('Unauthorized');
    }
    return res;
}

function logout() {
    localStorage.removeItem('token');
    window.location.href = '/admin/login';
}
