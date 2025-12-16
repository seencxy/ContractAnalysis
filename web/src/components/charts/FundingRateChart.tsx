import { useEffect, useRef, useMemo } from 'react';
import * as echarts from 'echarts';
import type { FundingRateData } from '@/hooks/queries/useBinanceKlines';
import dayjs from 'dayjs';

interface FundingRateChartProps {
  data: FundingRateData[];
  symbol: string;
}

export function FundingRateChart({ data, symbol }: FundingRateChartProps) {
  const chartRef = useRef<HTMLDivElement>(null);
  const chartInstance = useRef<echarts.ECharts | null>(null);

  const sortedData = useMemo(() => {
    if (!data || data.length === 0) return [];
    // Sort by time ascending
    return [...data].sort((a, b) => a.fundingTime - b.fundingTime);
  }, [data]);

  useEffect(() => {
    if (!chartRef.current || sortedData.length === 0) {
        if(chartInstance.current) chartInstance.current.clear();
        return;
    }

    if (!chartInstance.current) {
      chartInstance.current = echarts.init(chartRef.current);
    }

    const times = sortedData.map(item => dayjs(item.fundingTime).format('MM-DD HH:mm'));
    const rates = sortedData.map(item => parseFloat(item.fundingRate) * 100); // Convert to percentage

    const maxVal = Math.max(...rates);
    const minVal = Math.min(...rates);
    // Ensure we have a valid range even if all rates are 0
    const maxAbs = Math.max(Math.abs(maxVal), Math.abs(minVal)) || 0.01;
    
    // Add some padding to Y axis (10%)
    const yAxisLimit = maxAbs * 1.1;

    // Define colors
    const upColor = '#52c41a';
    const downColor = '#ff4d4f';

    const option: echarts.EChartsOption = {
      title: {
        text: '资金费率趋势',
        left: 0,
        top: 0,
        textStyle: {
            fontSize: 14,
            fontWeight: 500,
            color: '#666'
        }
      },
      tooltip: {
        trigger: 'axis',
        axisPointer: {
            type: 'line',
            lineStyle: {
                color: '#bfbfbf',
                width: 1,
                type: 'dashed'
            }
        },
        backgroundColor: 'rgba(255, 255, 255, 0.95)',
        borderColor: '#f0f0f0',
        borderWidth: 1,
        padding: [8, 12],
        textStyle: { color: '#333', fontSize: 12 },
        formatter: (params: any) => {
          const param = params[0];
          const rate = parseFloat(param.value);
          const color = rate > 0 ? upColor : downColor;
          const bg = rate > 0 ? 'rgba(82, 196, 26, 0.1)' : 'rgba(255, 77, 79, 0.1)';
          
          return `
            <div style="font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial;">
                <div style="color: #8c8c8c; margin-bottom: 6px; font-size: 11px;">${param.name}</div>
                <div style="display: flex; align-items: center; justify-content: space-between; gap: 12px;">
                    <span style="color: #595959;">费率</span>
                    <span style="font-weight: 600; color: ${color}; background: ${bg}; padding: 1px 6px; border-radius: 4px;">
                        ${rate.toFixed(4)}%
                    </span>
                </div>
            </div>
          `;
        },
      },
      grid: {
        left: '1%',
        right: '1%',
        bottom: '12%',
        top: '18%',
        containLabel: true,
      },
      xAxis: {
        type: 'category',
        data: times,
        axisLine: { show: false },
        axisTick: { show: false },
        axisLabel: { 
            color: '#bfbfbf', 
            fontSize: 10,
        },
        boundaryGap: false, 
      },
      yAxis: {
        type: 'value',
        scale: true,
        max: yAxisLimit,
        min: -yAxisLimit,
        splitNumber: 4,
        axisLabel: {
            formatter: (val: number) => `${val.toFixed(3)}%`,
            color: '#bfbfbf',
            fontSize: 10,
            showMinLabel: false,
            showMaxLabel: false
        },
        splitLine: {
            lineStyle: {
                type: 'dashed',
                color: '#f5f5f5'
            }
        }
      },
      series: [
        {
          name: 'Funding Rate',
          type: 'line',
          data: rates,
          smooth: 0.3,
          symbol: 'circle',
          symbolSize: 4,
          showSymbol: false,
          lineStyle: {
              width: 2,
              color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
                { offset: 0, color: upColor },      // Top (Positive)
                { offset: 0.5, color: upColor },    // Middle (Zero)
                { offset: 0.5, color: downColor },  // Middle (Zero)
                { offset: 1, color: downColor }     // Bottom (Negative)
              ])
          },
          areaStyle: {
              opacity: 0.2,
              color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
                { offset: 0, color: upColor },
                { offset: 0.5, color: upColor },
                { offset: 0.5, color: downColor },
                { offset: 1, color: downColor }
              ])
          },
          markLine: {
              symbol: ['none', 'none'],
              label: { show: false },
              lineStyle: {
                  color: '#e8e8e8',
                  width: 1,
                  type: 'solid'
              },
              data: [
                  { yAxis: 0 }
              ],
              animation: false,
              silent: true
          }
        }
      ],
      dataZoom: [
        {
            type: 'inside',
            start: 0,
            end: 100
        },
         {
          show: true,
          type: 'slider',
          height: 12,
          bottom: 0,
          borderColor: 'transparent',
          backgroundColor: '#f5f5f5',
          fillerColor: 'rgba(0, 0, 0, 0.05)',
          handleSize: '100%',
          handleStyle: {
             color: '#fff',
             borderColor: '#d9d9d9',
             shadowBlur: 2,
             shadowColor: 'rgba(0,0,0,0.1)'
          },
          textStyle: { color: 'transparent' }, 
          showDataShadow: false 
        }
      ]
    };

    chartInstance.current.setOption(option);

    // Use ResizeObserver to handle container resize
    const resizeObserver = new ResizeObserver(() => {
      chartInstance.current?.resize();
    });

    if (chartRef.current) {
      resizeObserver.observe(chartRef.current);
    }

    return () => {
      resizeObserver.disconnect();
    };
  }, [sortedData, symbol]);

  useEffect(() => {
    return () => {
      chartInstance.current?.dispose();
      chartInstance.current = null;
    };
  }, []);

  return <div ref={chartRef} style={{ height: 260, width: '100%' }} />;
}