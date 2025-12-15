import { useQuery, type UseQueryOptions } from '@tanstack/react-query';
import { statisticsApi, type StatisticsFilters, type StatisticsHistoryFilters } from '@/api/endpoints/statistics';
import type { ApiResponse } from '@/types/common';
import type { Statistics, OverviewStatistics } from '@/types/statistics';

export function useOverviewStatistics(
  options?: Omit<UseQueryOptions<ApiResponse<OverviewStatistics>>, 'queryKey' | 'queryFn'>
) {
  return useQuery({
    queryKey: ['statistics', 'overview'],
    queryFn: () => statisticsApi.getOverview(),
    staleTime: 60000, // 1分钟
    ...options,
  });
}

export function useStrategyStatistics(
  filters: StatisticsFilters = {},
  options?: Omit<UseQueryOptions<ApiResponse<Statistics[]>>, 'queryKey' | 'queryFn'>
) {
  return useQuery({
    queryKey: ['statistics', 'strategies', filters],
    queryFn: () => statisticsApi.getStrategies(filters),
    staleTime: 60000,
    ...options,
  });
}

export function useSymbolStatistics(
  filters: StatisticsFilters = {},
  options?: Omit<UseQueryOptions<ApiResponse<Statistics[]>>, 'queryKey' | 'queryFn'>
) {
  return useQuery({
    queryKey: ['statistics', 'symbols', filters],
    queryFn: () => statisticsApi.getSymbols(filters),
    staleTime: 60000,
    ...options,
  });
}

export function useStatisticsHistory(
  filters: StatisticsHistoryFilters,
  options?: Omit<UseQueryOptions<ApiResponse<Statistics[]>>, 'queryKey' | 'queryFn'>
) {
  return useQuery({
    queryKey: ['statistics', 'history', filters],
    queryFn: () => statisticsApi.getHistory(filters),
    enabled: !!filters.start_time && !!filters.end_time, // Only fetch when dates are set
    staleTime: 300000, // 5 minutes cache for historical data
    ...options,
  });
}
