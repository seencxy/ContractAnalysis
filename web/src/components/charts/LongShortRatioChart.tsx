import { useEffect, useRef, useMemo } from 'react';
import * as echarts from 'echarts';
import type { LSRatioData } from '@/hooks/queries/useBinanceKlines';
import dayjs from 'dayjs';

interface LongShortRatioChartProps {
  topTraderData: LSRatioData[];
  globalData: LSRatioData[];
  symbol: string;
}

export function LongShortRatioChart({ topTraderData, globalData, symbol }: LongShortRatioChartProps) {
  const chartRef = useRef<HTMLDivElement>(null);
  const chartInstance = useRef<echarts.ECharts | null>(null);

  // Merge and sort data
  const sortedSeries = useMemo(() => {
     if ((!topTraderData || topTraderData.length === 0) && (!globalData || globalData.length === 0)) return { times: [], top: [], global: [] };
     
     const timestamps = new Set([...(topTraderData?.map(d => d.timestamp) || []), ...(globalData?.map(d => d.timestamp) || [])]);
     const sortedTimestamps = Array.from(timestamps).sort((a, b) => a - b);
     
     const times: string[] = [];
     const top: (number | null)[] = [];
     const global: (number | null)[] = [];

     const topMap = new Map(topTraderData?.map(d => [d.timestamp, parseFloat(d.longShortRatio)]));
     const globalMap = new Map(globalData?.map(d => [d.timestamp, parseFloat(d.longShortRatio)]));

     sortedTimestamps.forEach(ts => {
         times.push(dayjs(ts).format('MM-DD HH:mm'));
         top.push(topMap.get(ts) ?? null);
         global.push(globalMap.get(ts) ?? null);
     });

     return { times, top, global };
  }, [topTraderData, globalData]);


  useEffect(() => {
    if (!chartRef.current || sortedSeries.times.length === 0) {
        if(chartInstance.current) chartInstance.current.clear();
        return;
    }

    if (!chartInstance.current) {
      chartInstance.current = echarts.init(chartRef.current);
    }

    const { times, top, global } = sortedSeries;
    
    const allValues = [...top, ...global].filter(v => v !== null) as number[];
    const maxVal = Math.max(...allValues);
    const minVal = Math.min(...allValues);
    const range = maxVal - minVal;
    
    const yMin = Math.max(0, minVal - range * 0.1); 
    const yMax = maxVal + range * 0.1;

    const option: echarts.EChartsOption = {
      title: {
        text: '多空比分析 (大户 vs 散户)',
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
        backgroundColor: 'rgba(255, 255, 255, 0.95)',
        borderColor: '#f0f0f0',
        borderWidth: 1,
        padding: [8, 12],
        textStyle: { color: '#333', fontSize: 12 },
        formatter: (params: any) => {
          let html = `<div style="color: #8c8c8c; margin-bottom: 6px; font-size: 11px;">${params[0].name}</div>`;
          
          params.forEach((param: any) => {
             const isTop = param.seriesName === '大户持仓比(Smart Money)';
             const val = parseFloat(param.value).toFixed(2);
             const badgeColor = val >= "1.00" ? '#52c41a' : '#ff4d4f';
             const badgeBg = val >= "1.00" ? 'rgba(82, 196, 26, 0.1)' : 'rgba(255, 77, 79, 0.1)';
             
             html += `
                <div style="display: flex; align-items: center; justify-content: space-between; gap: 12px; margin-bottom: 4px;">
                    <div style="display: flex; align-items: center; gap: 6px;">
                        <span style="display:inline-block;width:8px;height:8px;border-radius:50%;background-color:${param.color};"></span>
                        <span style="color: #595959;">${isTop ? '大户(资金)' : '全局(人数)'}</span>
                    </div>
                    <span style="font-weight: 600; color: ${badgeColor}; background: ${badgeBg}; padding: 1px 6px; border-radius: 4px;">
                        ${val}
                    </span>
                </div>
             `;
          });
          return html;
        }
      },
      grid: {
        left: '1%',
        right: '1%',
        bottom: '12%',
        top: '18%',
        containLabel: true,
      },
      legend: {
        data: ['大户持仓比(Smart Money)', '全局账户比(Sentiment)'],
        top: 0,
        right: 0,
        icon: 'circle',
        textStyle: { fontSize: 11, color: '#8c8c8c' },
        itemWidth: 8,
        itemHeight: 8
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
        min: yMin,
        max: yMax,
        splitNumber: 4,
        axisLabel: {
            color: '#bfbfbf',
            fontSize: 10,
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
          name: '大户持仓比(Smart Money)',
          type: 'line',
          data: top,
          smooth: true,
          showSymbol: false,
          symbolSize: 2,
          lineStyle: {
              width: 2.5,
              color: '#722ed1' // Purple for Smart Money
          },
          itemStyle: { color: '#722ed1' },
          markLine: {
              symbol: ['none', 'none'],
              label: { show: false },
              lineStyle: {
                  color: '#bfbfbf',
                  width: 1,
                  type: 'dashed'
              },
              data: [
                  { yAxis: 1.0 }
              ],
              silent: true
          }
        },
        {
          name: '全局账户比(Sentiment)',
          type: 'line',
          data: global,
          smooth: true,
          showSymbol: false,
          symbolSize: 2,
          lineStyle: {
              width: 1.5,
              color: '#13c2c2' // Cyan/Teal for Sentiment
          },
          itemStyle: { color: '#13c2c2' }
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

    // Add ResizeObserver
    const resizeObserver = new ResizeObserver(() => {
        chartInstance.current?.resize();
    });

    if (chartRef.current) {
        resizeObserver.observe(chartRef.current);
    }

    return () => {
        resizeObserver.disconnect();
    };
  }, [sortedSeries, symbol]);

  useEffect(() => {
    return () => {
      chartInstance.current?.dispose();
      chartInstance.current = null;
    };
  }, []);

  return <div ref={chartRef} style={{ height: 260, width: '100%' }} />;
}