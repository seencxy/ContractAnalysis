// 基础API响应结构(所有API返回都使用此结构)
export interface ApiResponse<T = unknown> {
  code: number;
  message: string;
  data?: T;
  error?: {
    type: string;
    details: string;
  };
  timestamp: number;
}

// 分页响应
export interface PaginationMeta {
  page: number;
  limit: number;
  total: number;
  total_pages: number;
}

export interface PaginatedData<T> {
  items: T[];
  pagination: PaginationMeta;
}
