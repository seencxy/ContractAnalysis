import apiClient from '../client';
import type { ApiResponse, PaginatedData } from '@/types/common';
import type { Signal, SignalTracking, SignalKlineTracking } from '@/types/signal';

export interface SignalFilters {
  page?: number;
  limit?: number;
  status?: string;
  symbol?: string;
  strategy_name?: string;
  type?: string;
  start_time?: string;
  end_time?: string;
}

export const signalsApi = {
  // 获取信号列表(分页)
  getSignals: async (filters: SignalFilters): Promise<ApiResponse<PaginatedData<Signal>>> => {
    return apiClient.get('/signals', { params: filters });
  },

  // 获取活跃信号
  getActiveSignals: async (): Promise<ApiResponse<Signal[]>> => {
    return apiClient.get('/signals/active');
  },

  // 获取信号详情
  getSignalById: async (signalId: string): Promise<ApiResponse<Signal>> => {
    return apiClient.get(`/signals/${signalId}`);
  },

  // 获取信号追踪记录
  getSignalTracking: async (signalId: string): Promise<ApiResponse<SignalTracking[]>> => {
    return apiClient.get(`/signals/${signalId}/tracking`);
  },

  // 获取信号K线数据
  getSignalKlines: async (signalId: string): Promise<ApiResponse<SignalKlineTracking[]>> => {
    return apiClient.get(`/signals/${signalId}/klines`);
  },
};
