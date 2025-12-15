import apiClient from '../client';
import type { ApiResponse } from '@/types/common';
import type { Statistics, OverviewStatistics } from '@/types/statistics';

export interface StatisticsFilters {
  period?: '24h' | '7d' | '30d' | 'all';
  strategy?: string;
  symbol?: string;
}

export interface StatisticsHistoryFilters {
  start_time: string;  // ISO 8601 format
  end_time: string;
  strategy?: string;
  symbol?: string;
}

export const statisticsApi = {
  // 获取统计概览
  getOverview: async (): Promise<ApiResponse<OverviewStatistics>> => {
    return apiClient.get('/statistics/overview');
  },

  // 获取策略统计
  getStrategies: async (filters?: StatisticsFilters): Promise<ApiResponse<Statistics[]>> => {
    return apiClient.get('/statistics/strategies', { params: filters });
  },

  // 获取交易对统计
  getSymbols: async (filters?: StatisticsFilters): Promise<ApiResponse<Statistics[]>> => {
    return apiClient.get('/statistics/symbols', { params: filters });
  },

  // 获取历史统计
  getHistory: async (filters: StatisticsHistoryFilters): Promise<ApiResponse<Statistics[]>> => {
    return apiClient.get('/statistics/history', { params: filters });
  },
};
