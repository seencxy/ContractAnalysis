import { useQuery, type UseQueryOptions } from '@tanstack/react-query';
import { signalsApi, type SignalFilters } from '@/api/endpoints/signals';
import type { ApiResponse, PaginatedData } from '@/types/common';
import type { Signal, SignalTracking, SignalKlineTracking } from '@/types/signal';

export function useSignals(
  filters: SignalFilters = {},
  options?: Omit<UseQueryOptions<ApiResponse<PaginatedData<Signal>>>, 'queryKey' | 'queryFn'>
) {
  return useQuery({
    queryKey: ['signals', filters],
    queryFn: () => signalsApi.getSignals(filters),
    staleTime: 30000, // 30秒
    ...options,
  });
}

export function useActiveSignals(
  options?: Omit<UseQueryOptions<ApiResponse<Signal[]>>, 'queryKey' | 'queryFn'>
) {
  return useQuery({
    queryKey: ['signals', 'active'],
    queryFn: () => signalsApi.getActiveSignals(),
    refetchInterval: 30000, // 30秒自动刷新
    ...options,
  });
}

export function useSignalDetail(
  signalId: string,
  options?: Omit<UseQueryOptions<ApiResponse<Signal>>, 'queryKey' | 'queryFn'>
) {
  return useQuery({
    queryKey: ['signals', signalId],
    queryFn: () => signalsApi.getSignalById(signalId),
    enabled: !!signalId,
    ...options,
  });
}

export function useSignalTracking(
  signalId: string,
  options?: Omit<UseQueryOptions<ApiResponse<SignalTracking[]>>, 'queryKey' | 'queryFn'>
) {
  return useQuery({
    queryKey: ['signals', signalId, 'tracking'],
    queryFn: () => signalsApi.getSignalTracking(signalId),
    enabled: !!signalId,
    ...options,
  });
}

export function useSignalKlines(
  signalId: string,
  options?: Omit<UseQueryOptions<ApiResponse<SignalKlineTracking[]>>, 'queryKey' | 'queryFn'>
) {
  return useQuery({
    queryKey: ['signals', signalId, 'klines'],
    queryFn: () => signalsApi.getSignalKlines(signalId),
    enabled: !!signalId,
    ...options,
  });
}
