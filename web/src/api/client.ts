import axios from 'axios';

const apiClient = axios.create({
  baseURL: '/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// 请求拦截器
apiClient.interceptors.request.use(
  (config) => {
    // 可以在这里添加认证token等
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// 响应拦截器
apiClient.interceptors.response.use(
  (response) => {
    // API返回格式: { code, message, data, timestamp }
    // 直接返回整个response.data(包含code, message, data等)
    return response.data;
  },
  (error) => {
    // 统一错误处理
    const status = error.response?.status;
    const message = error.response?.data?.message || '请求失败';
    
    if (status === 422) {
      const errorMsg = error.response?.data?.error || '参数验证失败';
      console.error('Validation Error:', errorMsg);
      // Optional: You could trigger a UI notification here if you have a global notification system
    } else {
      console.error('API Error:', message);
    }
    
    return Promise.reject(error);
  }
);

export default apiClient;
