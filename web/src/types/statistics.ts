export interface Statistics {
  strategy_name: string;
  symbol?: string;
  period_label: string;
  period_start: string;
  period_end: string;

  // Signal counts
  total_signals: number;
  confirmed_signals: number;
  invalidated_signals: number;

  // Outcome counts
  profitable_signals: number;
  losing_signals: number;
  neutral_signals: number;

  // Performance metrics
  win_rate?: string;
  avg_profit_pct?: string;
  avg_loss_pct?: string;
  avg_holding_hours?: string;
  best_signal_pct?: string;
  worst_signal_pct?: string;
  profit_factor?: string;

  // K-line metrics
  kline_theoretical_win_rate?: string;
  kline_close_win_rate?: string;
  total_kline_hours: number;
  profitable_kline_hours_high: number;
  profitable_kline_hours_close: number;

  // Hourly return statistics
  avg_hourly_return_pct?: string;
  max_hourly_return_pct?: string;
  min_hourly_return_pct?: string;

  // Theoretical maximum profit/loss
  avg_max_potential_profit_pct?: string;
  avg_max_potential_loss_pct?: string;

  calculated_at: string;
}

export type StatisticsResponse = Statistics;

export interface StrategyPerformance24h {
  strategy_name: string;
  signal_count: number;
  win_rate?: string;
  avg_return_pct?: string;
  profitable_count: number;
  losing_count: number;
}

export interface SignalStatusDistribution {
  pending: number;
  confirmed: number;
  tracking: number;
  closed: number;
  invalidated: number;
}

export interface OverviewStatistics {
  total_signals_today: number;
  active_signals: number;
  overall_win_rate_24h?: string;
  avg_return_pct_24h?: string;
  strategy_breakdown?: StrategyPerformance24h[];
  top_performing_pair?: string;
  worst_performing_pair?: string;
  status_distribution?: SignalStatusDistribution;
}

export interface ComparisonMetrics {
  win_rates: Record<string, string>;
  avg_returns: Record<string, string>;
  total_signals: Record<string, number>;
  profit_factors: Record<string, string>;
  best_win_rate: string;
  best_avg_return: string;
  most_signals: string;
}

export interface StrategyComparisonResponse {
  period: string;
  strategies: string[];
  comparison: ComparisonMetrics;
  detailed_stats: StatisticsResponse[];
}
