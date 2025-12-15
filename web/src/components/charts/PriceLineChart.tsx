import { useEffect, useRef } from 'react';
import * as echarts from 'echarts';
import type { SignalTracking, SignalType } from '@/types/signal';
import dayjs from 'dayjs';

interface PriceLineChartProps {
  trackings: SignalTracking[];
  signalPrice: string;
  signalType?: SignalType;
}

export function PriceLineChart({ trackings, signalPrice, signalType = 'LONG' }: PriceLineChartProps) {
  const chartRef = useRef<HTMLDivElement>(null);
  const chartInstance = useRef<echarts.ECharts | null>(null);

  useEffect(() => {
    if (!chartRef.current) return;

    // Dispose existing instance if it exists (handles React StrictMode)
    if (chartInstance.current) {
      chartInstance.current.dispose();
    }
    chartInstance.current = echarts.init(chartRef.current);

    const price = parseFloat(signalPrice);

    // Filter and sort trackings first to ensure data alignment
    const validTrackings = trackings
      .filter((t) => t.current_price && !isNaN(parseFloat(t.current_price)))
      .sort((a, b) => new Date(a.tracked_at).getTime() - new Date(b.tracked_at).getTime());
      
    if (isNaN(price) || validTrackings.length === 0) {
      chartInstance.current?.clear();
      return;
    }

    const prices = validTrackings.map((t) => parseFloat(t.current_price));
    const minData = Math.min(...prices, price);
    const maxData = Math.max(...prices, price);
    
    // Ensure we have a valid range for Y-axis
    const range = maxData - minData;
    let padding = 0;
    
    if (range === 0) {
      // If flat, add +/- 0.5% buffer (min 0.01)
      padding = Math.max(Math.abs(price) * 0.005, 0.01);
    } else {
      padding = range * 0.1;
    }

    // Define colors based on profit/loss status
    const profitColor = '#52c41a'; // Green
    const lossColor = '#ff4d4f';   // Red
    
    // For LONG: Price > Signal => Profit (Green), Price < Signal => Loss (Red)
    // For SHORT: Price < Signal => Profit (Green), Price > Signal => Loss (Red)
    const isLong = signalType === 'LONG';

    // Determine line color based on last price point
    const lastPrice = prices[prices.length - 1];
    const isProfit = isLong ? lastPrice > price : lastPrice < price;
    const lineColor = isProfit ? profitColor : lossColor;

    const option: echarts.EChartsOption = {
      title: {
        text: '价格追踪走势',
        left: 'left',
        textStyle: {
          fontSize: 14,
          fontWeight: 'normal',
          color: '#666'
        }
      },
      tooltip: {
        trigger: 'axis',
        axisPointer: {
          type: 'cross',
          label: {
            backgroundColor: '#6a7985'
          }
        },
        formatter: (params: any) => {
          const data = params[0];
          if (!data) return '';
          
          const tracking = validTrackings[data.dataIndex];
          if (!tracking) return '';

          const currentPrice = parseFloat(tracking.current_price);
          const changePct = parseFloat(tracking.price_change_pct);
          const isPointProfit = isLong ? currentPrice > price : currentPrice < price;
          const color = isPointProfit ? profitColor : lossColor;
          
          return `
            <div style="font-size: 12px;">
              <div style="margin-bottom: 4px; color: #999;">
                ${dayjs(tracking.tracked_at).format('YYYY-MM-DD HH:mm')}
              </div>
              <div style="font-weight: bold; margin-bottom: 4px;">
                价格: ${currentPrice}
              </div>
              <div style="color: ${color}">
                幅度: ${changePct > 0 ? '+' : ''}${changePct}%
              </div>
            </div>
          `;
        },
      },
      grid: {
        left: '3%',
        right: '4%',
        bottom: '15%',
        top: '15%',
        containLabel: true,
      },
      xAxis: {
        type: 'category',
        boundaryGap: false,
        data: validTrackings.map((t) => 
          dayjs(t.tracked_at).format('MM-DD HH:mm')
        ),
        axisLabel: {
          rotate: 0,
          formatter: (value: string) => value.split(' ').join('\n'),
          color: '#999',
          fontSize: 10
        },
        axisLine: {
          lineStyle: {
            color: '#eee'
          }
        }
      },
      yAxis: {
        type: 'value',
        scale: true,
        min: minData - padding,
        max: maxData + padding,
        splitLine: {
          lineStyle: {
            color: '#f5f5f5'
          }
        },
        axisLabel: {
          color: '#999',
          formatter: (value: number) => value.toFixed(price < 1 ? 4 : 2)
        }
      },
      dataZoom: [
        {
          type: 'inside',
          start: 0,
          end: 100
        },
        {
          start: 0,
          end: 100,
          height: 20,
          bottom: 0
        }
      ],
      series: [
        {
          name: '价格',
          type: 'line',
          data: prices,
          smooth: true,
          symbol: prices.length === 1 ? 'circle' : 'none',
          symbolSize: 6,
          lineStyle: {
            width: 2,
            color: lineColor
          },
          itemStyle: {
            color: lineColor
          },
          areaStyle: {
            opacity: 0.1,
            color: lineColor
          },
          markLine: {
            symbol: ['none', 'none'],
            label: {
              position: 'insideEndTop',
              formatter: '{b}: {c}'
            },
            data: [
              {
                yAxis: price,
                name: '开仓价',
                lineStyle: {
                  color: '#666',
                  type: 'dashed',
                  width: 1
                },
                label: {
                  formatter: `开仓: ${price}`
                }
              }
            ]
          },
          markPoint: {
            symbol: 'circle',
            symbolSize: 6,
            label: {
              show: true,
              position: 'top',
              formatter: '{c}',
              fontSize: 10,
              color: '#666'
            },
            data: [
              { type: 'max', name: '最高' },
              { type: 'min', name: '最低' }
            ]
          }
        },
      ],
    };

    chartInstance.current.setOption(option, true);

    const handleResize = () => {
      chartInstance.current?.resize();
    };

    window.addEventListener('resize', handleResize);

    // Delay resize to ensure container is properly sized (e.g., after Drawer animation)
    const resizeTimer = setTimeout(() => {
      chartInstance.current?.resize();
    }, 350);

    return () => {
      clearTimeout(resizeTimer);
      window.removeEventListener('resize', handleResize);
      if (chartInstance.current && !chartInstance.current.isDisposed()) {
        chartInstance.current.dispose();
        chartInstance.current = null;
      }
    };
  }, [trackings, signalPrice, signalType]);

  return <div ref={chartRef} style={{ height: '400px', width: '100%' }} />;
}
