import apiClient from '../client';
import type { ApiResponse } from '@/types/common';
import type { Statistics, OverviewStatistics, StrategyComparisonResponse } from '@/types/statistics';

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

export interface StrategyCompareParams {
  strategies: string[];
  period: '24h' | '7d' | '30d' | 'all';
  symbols?: string[];
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

  // 策略对比
  compareStrategies: async (params: StrategyCompareParams): Promise<ApiResponse<StrategyComparisonResponse>> => {
    const searchParams = new URLSearchParams();
    searchParams.append('period', params.period);
    
    // Append strategies as repeated parameters: strategies=A&strategies=B
    params.strategies.forEach(s => searchParams.append('strategies', s));
    
    if (params.symbols) {
      params.symbols.forEach(s => searchParams.append('symbols', s));
    }

    return apiClient.get('/statistics/compare', { params: searchParams });
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
