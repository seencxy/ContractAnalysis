export type SignalType = 'LONG' | 'SHORT';
export type SignalStatus = 'PENDING' | 'CONFIRMED' | 'TRACKING' | 'CLOSED' | 'INVALIDATED';

export interface Signal {
  signal_id: string;
  symbol: string;
  type: SignalType;
  strategy_name: string;
  generated_at: string;
  price_at_signal: string;
  long_account_ratio: string;
  short_account_ratio: string;
  long_position_ratio: string;
  short_position_ratio: string;
  long_trader_count: number;
  short_trader_count: number;
  status: SignalStatus;
  is_confirmed: boolean;
  confirmed_at?: string;
  reason?: string;
  created_at: string;
  updated_at: string;
}

export interface SignalTracking {
  id: number;
  signal_id: string;
  tracked_at: string;
  current_price: string;
  price_change_pct: string;
  highest_price?: string;
  lowest_price?: string;
  highest_change_pct?: string;
  lowest_change_pct?: string;
  hours_tracked: number;
  is_profit_target_hit: boolean;
  is_stop_loss_hit: boolean;
}

export interface SignalKlineTracking {
  id: number;
  signal_id: string;
  kline_open_time: string;
  kline_close_time: string;
  open_price: string;
  high_price: string;
  low_price: string;
  close_price: string;
  volume: string;
  open_change_pct: string;
  high_change_pct: string;
  low_change_pct: string;
  close_change_pct: string;
  hourly_return_pct?: string;
  is_profitable_at_high: boolean;
  is_profitable_at_close: boolean;
}
