import { useEffect, useRef } from 'react';
import * as echarts from 'echarts';
import type { Signal } from '@/types/signal';

interface SignalStatusPieChartProps {
  signals?: Signal[];
  distribution?: Record<string, number>;
}

export function SignalStatusPieChart({ signals = [], distribution }: SignalStatusPieChartProps) {
  const chartRef = useRef<HTMLDivElement>(null);
  const chartInstance = useRef<echarts.ECharts | null>(null);

  useEffect(() => {
    if (!chartRef.current) return;

    if (!chartInstance.current) {
      chartInstance.current = echarts.init(chartRef.current);
    }

    let statusCount: Record<string, number> = {};

    if (distribution) {
      // Normalize keys to uppercase to match getStatusLabel expectations
      Object.entries(distribution).forEach(([key, value]) => {
        statusCount[key.toUpperCase()] = value;
      });
    } else {
      statusCount = signals.reduce((acc, signal) => {
        acc[signal.status] = (acc[signal.status] || 0) + 1;
        return acc;
      }, {} as Record<string, number>);
    }

    const data = Object.entries(statusCount).map(([name, value]) => ({
      name: getStatusLabel(name),
      value,
    }));

    const option = {
      color: ['#1677ff', '#52c41a', '#faad14', '#ff4d4f', '#bfbfbf'],
      tooltip: {
        trigger: 'item',
        formatter: '{b}: {c} ({d}%)',
        backgroundColor: 'rgba(255, 255, 255, 0.9)',
        borderColor: '#f0f0f0',
        borderWidth: 1,
        textStyle: {
            color: '#333'
        }
      },
      legend: {
        bottom: '0%',
        left: 'center',
        itemWidth: 10,
        itemHeight: 10,
        textStyle: {
            fontSize: 12
        }
      },
      series: [
        {
          name: '信号状态',
          type: 'pie',
          radius: ['45%', '70%'],
          center: ['50%', '45%'],
          itemStyle: {
            borderRadius: 6,
            borderColor: '#fff',
            borderWidth: 2,
          },
          label: {
            show: false,
            position: 'center'
          },
          emphasis: {
            label: {
              show: true,
              fontSize: 16,
              fontWeight: 'bold',
            },
            scale: true,
            scaleSize: 10
          },
          data,
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
  }, [signals, distribution]);

  return <div ref={chartRef} style={{ height: '100%', minHeight: 300, width: '100%' }} />;
}

function getStatusLabel(status: string): string {
  const labels: Record<string, string> = {
    PENDING: '待确认',
    CONFIRMED: '已确认',
    TRACKING: '追踪中',
    CLOSED: '已关闭',
    INVALIDATED: '已失效',
  };
  return labels[status] || status;
}