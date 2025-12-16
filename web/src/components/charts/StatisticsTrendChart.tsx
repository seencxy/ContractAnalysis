import { useEffect, useRef, useMemo } from 'react';
import * as echarts from 'echarts';
import type { Statistics } from '@/types/statistics';
import dayjs from 'dayjs';

interface StatisticsTrendChartProps {
  data: Statistics[];
  metricType: 'win_rate' | 'avg_profit' | 'profit_factor' | 'total_signals';
}

interface MetricConfig {
  name: string;
  unit: string;
  dataKey: keyof Statistics;
  color: string;
}

export function StatisticsTrendChart({ data, metricType }: StatisticsTrendChartProps) {
  const chartRef = useRef<HTMLDivElement>(null);
  const chartInstance = useRef<echarts.ECharts | null>(null);

  const processedData = useMemo(() => {
    if (!data || data.length === 0) return [];

    // Group by calculated_at
    const groups = new Map<string, Statistics[]>();
    data.forEach((item) => {
      const key = item.calculated_at;
      if (!groups.has(key)) {
        groups.set(key, []);
      }
      groups.get(key)!.push(item);
    });

    // Sort by time
    const sortedKeys = Array.from(groups.keys()).sort((a, b) =>
      new Date(a).getTime() - new Date(b).getTime()
    );

    return sortedKeys.map((key) => {
      const items = groups.get(key)!;
      const totalSignals = items.reduce((sum, item) => sum + item.total_signals, 0);
      const totalProfitable = items.reduce((sum, item) => sum + item.profitable_signals, 0);

      // Helper to get float from string
      const getFloat = (val: string | undefined) => (val ? parseFloat(val) : 0);

      // Weighted Sums
      let weightedAvgProfit = 0;
      let weightedProfitFactor = 0;

      items.forEach((item) => {
        weightedAvgProfit += getFloat(item.avg_profit_pct) * item.total_signals;
        weightedProfitFactor += getFloat(item.profit_factor) * item.total_signals;
      });

      return {
        calculated_at: key,
        total_signals: totalSignals,
        win_rate: totalSignals > 0 ? ((totalProfitable / totalSignals) * 100).toFixed(2) : 0,
        avg_profit_pct: totalSignals > 0 ? (weightedAvgProfit / totalSignals).toFixed(2) : 0,
        profit_factor: totalSignals > 0 ? (weightedProfitFactor / totalSignals).toFixed(2) : 0,
      };
    });
  }, [data]);

  const getMetricConfig = (): MetricConfig => {
    switch (metricType) {
      case 'win_rate':
        return {
          name: '胜率',
          unit: '%',
          dataKey: 'win_rate',
          color: '#1677ff',
        };
      case 'avg_profit':
        return {
          name: '平均盈利',
          unit: '%',
          dataKey: 'avg_profit_pct',
          color: '#52c41a',
        };
      case 'profit_factor':
        return {
          name: '盈利因子',
          unit: '',
          dataKey: 'profit_factor',
          color: '#722ed1',
        };
      case 'total_signals':
        return {
          name: '信号总数',
          unit: '',
          dataKey: 'total_signals',
          color: '#fa8c16',
        };
    }
  };

  useEffect(() => {
    if (!chartRef.current || processedData.length === 0) return;

    if (!chartInstance.current) {
      chartInstance.current = echarts.init(chartRef.current);
    }

    const config = getMetricConfig();

    const option: echarts.EChartsOption = {
      tooltip: {
        trigger: 'axis',
        backgroundColor: 'rgba(255, 255, 255, 0.95)',
        borderColor: '#f0f0f0',
        borderWidth: 1,
        textStyle: {
            color: '#333'
        },
        formatter: (params: any) => {
          const param = params[0];
          return `<div style="margin-bottom: 4px; color: #888;">${dayjs(param.name).format('YYYY-MM-DD HH:mm')}</div>
                  <div style="font-weight: bold; color: ${config.color}">${config.name}: ${param.value}${config.unit}</div>`;
        },
      },
      xAxis: {
        type: 'category',
        data: processedData.map((item) => item.calculated_at),
        axisLabel: {
          formatter: (value: string) => dayjs(value).format('MM-DD HH:mm'),
          color: '#8c8c8c'
        },
        axisLine: {
            lineStyle: { color: '#f0f0f0' }
        },
        axisTick: { show: false }
      },
      yAxis: {
        type: 'value',
        axisLabel: {
          formatter: `{value}${config.unit}`,
          color: '#8c8c8c'
        },
        splitLine: {
            lineStyle: {
                type: 'dashed',
                color: '#f0f0f0'
            }
        }
      },
      grid: {
        left: '2%',
        right: '4%',
        bottom: '5%',
        top: '10%',
        containLabel: true,
      },
      series: [
        {
          name: config.name,
          type: 'line',
          data: processedData.map((item) => {
            const value = item[config.dataKey as keyof typeof item];
            if (typeof value === 'string') {
              return parseFloat(value) || 0;
            }
            return value || 0;
          }),
          smooth: true,
          showSymbol: false,
          symbolSize: 8,
          itemStyle: {
            color: config.color,
            borderWidth: 2,
            borderColor: '#fff'
          },
          lineStyle: {
            width: 3
          },
          areaStyle: {
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: config.color + '33' }, // 20% opacity
              { offset: 1, color: config.color + '00' }, // 0% opacity
            ]),
          },
          emphasis: {
             focus: 'series'
          }
        },
      ],
    };

    chartInstance.current.setOption(option);

    const handleResize = () => {
      chartInstance.current?.resize();
    };

    window.addEventListener('resize', handleResize);

    return () => {
      window.removeEventListener('resize', handleResize);
    };
  }, [processedData, metricType]);

  useEffect(() => {
    return () => {
      chartInstance.current?.dispose();
      chartInstance.current = null;
    };
  }, []);

  return <div ref={chartRef} style={{ height: 350, width: '100%' }} />;
}