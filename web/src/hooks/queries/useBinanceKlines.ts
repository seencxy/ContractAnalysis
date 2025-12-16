import { useQuery, type UseQueryOptions } from '@tanstack/react-query';
import axios from 'axios';
import type { KlineData } from '@/types/kline';

interface BinanceKlineParams {
  symbol: string;
  interval: string;
  startTime?: number;
  endTime?: number;
  limit?: number;
}

// Binance K-line format:
// [
//   1499040000000,      // Open time
//   "0.01634790",       // Open
//   "0.80000000",       // High
//   "0.01575800",       // Low
//   "0.01577100",       // Close
//   "148976.11427815",  // Volume
//   1499644799999,      // Close time
//   ...
// ]
type BinanceKlineResponse = [
  number, // Open time
  string, // Open
  string, // High
  string, // Low
  string, // Close
  string, // Volume
  number, // Close time
  string, // Quote asset volume
  number, // Number of trades
  string, // Taker buy base asset volume
  string, // Taker buy quote asset volume
  string  // Ignore
][];

const binanceClient = axios.create({
  baseURL: '/binance-api',
  timeout: 10000,
});

export const fetchBinanceKlines = async (params: BinanceKlineParams): Promise<KlineData[]> => {
  const { data } = await binanceClient.get<BinanceKlineResponse>('/fapi/v1/klines', {
    params: {
      symbol: params.symbol.toUpperCase(), // Ensure symbol is uppercase
      interval: params.interval,
      startTime: params.startTime,
      endTime: params.endTime,
      limit: params.limit || 500,
    },
  });

  return data.map((item) => ({
    time: item[0],
    open: parseFloat(item[1]),
    high: parseFloat(item[2]),
    low: parseFloat(item[3]),
    close: parseFloat(item[4]),
    volume: parseFloat(item[5]),
  }));
};

export interface FundingRateData {
  symbol: string;
  fundingTime: number;
  fundingRate: string;
}

export interface LSRatioData {
  symbol: string;
  longShortRatio: string;
  longAccount: string;
  shortAccount: string;
  timestamp: number;
}

export const fetchBinanceFundingRate = async (params: BinanceKlineParams): Promise<FundingRateData[]> => {
  const { data } = await binanceClient.get<FundingRateData[]>('/fapi/v1/fundingRate', {
    params: {
      symbol: params.symbol.toUpperCase(),
      startTime: params.startTime,
      endTime: params.endTime,
      limit: params.limit || 100,
    },
  });
  return data;
};

export const fetchBinanceLSRatios = async (params: BinanceKlineParams) => {
  const commonParams = {
    symbol: params.symbol.toUpperCase(),
    period: params.interval, // Binance uses 'period' for LS ratios, usually 5m, 15m, 30m, 1h...
    startTime: params.startTime,
    endTime: params.endTime,
    limit: params.limit || 100,
  };

  const [topTraderRes, globalRes] = await Promise.all([
    binanceClient.get<LSRatioData[]>('/futures/data/topLongShortPositionRatio', { params: commonParams }),
    binanceClient.get<LSRatioData[]>('/futures/data/globalLongShortAccountRatio', { params: commonParams }),
  ]);

  return {
    topTrader: topTraderRes.data,
    global: globalRes.data,
  };
};

export function useBinanceFundingRate(
  params: BinanceKlineParams,
  options?: Omit<UseQueryOptions<FundingRateData[]>, 'queryKey' | 'queryFn'>
) {
  return useQuery({
    queryKey: ['binance', 'fundingRate', params],
    queryFn: () => fetchBinanceFundingRate(params),
    ...options,
  });
}

export function useBinanceLSRatios(
  params: BinanceKlineParams,
  options?: Omit<UseQueryOptions<{ topTrader: LSRatioData[], global: LSRatioData[] }>, 'queryKey' | 'queryFn'>
) {
  return useQuery({
    queryKey: ['binance', 'lsRatios', params],
    queryFn: () => fetchBinanceLSRatios(params),
    ...options,
  });
}

export function useBinanceKlines(
  params: BinanceKlineParams,
  options?: Omit<UseQueryOptions<KlineData[]>, 'queryKey' | 'queryFn'>
) {
  return useQuery({
    queryKey: ['binance', 'klines', params],
    queryFn: () => fetchBinanceKlines(params),
    ...options,
  });
}
