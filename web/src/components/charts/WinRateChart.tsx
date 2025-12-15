import { useEffect, useRef } from 'react';
import * as echarts from 'echarts';
import type { Statistics } from '@/types/statistics';

interface WinRateChartProps {
  data: Statistics[];
}

export function WinRateChart({ data }: WinRateChartProps) {
  const chartRef = useRef<HTMLDivElement>(null);
  const chartInstance = useRef<echarts.ECharts | null>(null);

  useEffect(() => {
    if (!chartRef.current) return;

    if (!chartInstance.current) {
      chartInstance.current = echarts.init(chartRef.current);
    }

    const option = {
      title: {
        text: '策略胜率对比',
        left: 'center',
      },
      tooltip: {
        trigger: 'axis',
        axisPointer: {
          type: 'shadow',
        },
      },
      legend: {
        data: ['基础胜率', 'K线最高价胜率', 'K线收盘价胜率'],
        bottom: 0,
      },
      grid: {
        left: '3%',
        right: '4%',
        bottom: '10%',
        containLabel: true,
      },
      xAxis: {
        type: 'category',
        data: data.map((d) => d.strategy_name),
        axisLabel: {
          rotate: 30,
        },
      },
      yAxis: {
        type: 'value',
        name: '胜率 (%)',
        axisLabel: {
          formatter: '{value}%',
        },
      },
      series: [
        {
          name: '基础胜率',
          type: 'bar',
          data: data.map((d) => (d.win_rate ? parseFloat(d.win_rate) : 0)),
          itemStyle: {
            color: '#1890ff',
          },
        },
        {
          name: 'K线最高价胜率',
          type: 'bar',
          data: data.map((d) =>
            d.kline_theoretical_win_rate ? parseFloat(d.kline_theoretical_win_rate) : 0
          ),
          itemStyle: {
            color: '#52c41a',
          },
        },
        {
          name: 'K线收盘价胜率',
          type: 'bar',
          data: data.map((d) =>
            d.kline_close_win_rate ? parseFloat(d.kline_close_win_rate) : 0
          ),
          itemStyle: {
            color: '#faad14',
          },
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
      chartInstance.current?.dispose();
      chartInstance.current = null;
    };
  }, [data]);

  return <div ref={chartRef} style={{ height: '100%', minHeight: 400, width: '100%' }} />;
}
