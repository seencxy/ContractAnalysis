import { useEffect, useRef } from 'react';
import * as echarts from 'echarts';
import { Radio } from 'antd';
import type { SignalType } from '@/types/signal';
import type { KlineData } from '@/types/kline';
import dayjs from 'dayjs';

interface CandlestickChartProps {
  klines: KlineData[];
  signalPrice: string;
  signalType?: SignalType;
  signalTime?: string;
  confirmedAt?: string; // New prop for confirmed time
  closedAt?: string;    // New prop for closed time
  interval?: string;
  onIntervalChange?: (interval: string) => void;
}

export function CandlestickChart({ klines, signalPrice, signalType = 'LONG', signalTime, confirmedAt, closedAt, interval = '15m', onIntervalChange }: CandlestickChartProps) {
  const chartRef = useRef<HTMLDivElement>(null);
  const chartInstance = useRef<echarts.ECharts | null>(null);

  useEffect(() => {
    if (!chartRef.current) return;

    // Dispose existing instance if it exists (handles React StrictMode)
    if (chartInstance.current) {
      chartInstance.current.dispose();
    }
    chartInstance.current = echarts.init(chartRef.current);

    // Sort klines by time
    const sortedKlines = [...klines].sort((a, b) => a.time - b.time);

    if (sortedKlines.length === 0) {
      chartInstance.current.clear();
      return;
    }

    const dates = sortedKlines.map((k) => dayjs(k.time).format('MM-DD HH:mm'));
    
    // [open, close, low, high, volume, sign]
    const data = sortedKlines.map((k) => [
      k.open,
      k.close,
      k.low,
      k.high,
      k.volume,
      k.close > k.open ? 1 : -1
    ]);

    const price = parseFloat(signalPrice);
    
    // Use a neutral/distinct color for the signal price line (Purple) to avoid confusion with Long/Short candle colors
    const signalLineColor = '#722ed1'; 

    // Calculate Y-axis range to ensure signal price is visible
    const allLow = sortedKlines.map(k => k.low);
    const allHigh = sortedKlines.map(k => k.high);
    let minPrice = Math.min(...allLow);
    let maxPrice = Math.max(...allHigh);

    // Expand range to include signal price
    if (!isNaN(price)) {
      minPrice = Math.min(minPrice, price);
      maxPrice = Math.max(maxPrice, price);
    }
    
    // Add some padding (5%)
    const range = maxPrice - minPrice;
    minPrice = minPrice - range * 0.05;
    maxPrice = maxPrice + range * 0.05;

    const profitColor = '#00da3c'; // Green
    const lossColor = '#ec0000';   // Red

    // Find closest time for signal markLine
    let signalDateStr = '';
    if (signalTime) {
      const signalTs = dayjs(signalTime).valueOf();
      const closestKline = sortedKlines.reduce((prev, curr) => {
        return (Math.abs(curr.time - signalTs) < Math.abs(prev.time - signalTs) ? curr : prev);
      });
      
      if (closestKline) {
         signalDateStr = dayjs(closestKline.time).format('MM-DD HH:mm');
      }
    }

    let confirmedDateStr = '';
    if (confirmedAt) {
      const confirmedTs = dayjs(confirmedAt).valueOf();
      const closestKline = sortedKlines.reduce((prev, curr) => {
        return (Math.abs(curr.time - confirmedTs) < Math.abs(prev.time - confirmedTs) ? curr : prev);
      });
      if (closestKline) {
        confirmedDateStr = dayjs(closestKline.time).format('MM-DD HH:mm');
      }
    }

    let closedDateStr = '';
    if (closedAt) {
      const closedTs = dayjs(closedAt).valueOf();
      const closestKline = sortedKlines.reduce((prev, curr) => {
        return (Math.abs(curr.time - closedTs) < Math.abs(prev.time - closedTs) ? curr : prev);
      });
      if (closestKline) {
        closedDateStr = dayjs(closestKline.time).format('MM-DD HH:mm');
      }
    }

    type MarkLineDataItem = echarts.SeriesOption['markLine']['data'][number];
    const markLineData: MarkLineDataItem[] = [];
    
    // Y-Axis MarkLine (Price)
    if (!isNaN(price)) {
        markLineData.push({
            yAxis: price,
            name: 'Signal Price',
            lineStyle: {
                color: signalLineColor, // Use determined color
                type: 'solid', // Make it solid for prominence
                width: 2 // Thicker line
            },
            label: {
                formatter: `Signal Price: ${price}`,
                position: 'end',
                backgroundColor: 'rgba(255, 255, 255, 0.8)', // Add background for readability
                padding: [4, 8],
                borderRadius: 4,
                color: signalLineColor,
                fontWeight: 'bold'
            }
        });
    }

    // X-Axis MarkLine (Signal Time)
    if (signalDateStr) {
        markLineData.push({
            xAxis: signalDateStr,
            name: 'Signal Time',
            lineStyle: {
                color: '#1677ff', // Blue for time
                type: 'dashed',
                width: 1
            },
            label: {
                formatter: 'Signal Time',
                position: 'start',
                backgroundColor: 'rgba(255, 255, 255, 0.8)',
                padding: [4, 8],
                borderRadius: 4,
                color: '#1677ff'
            }
        });
    }

    // X-Axis MarkLine (Confirmed Time)
    if (confirmedDateStr) {
      markLineData.push({
        xAxis: confirmedDateStr,
        name: 'Confirmed Time',
        lineStyle: {
          color: '#52c41a', // Green for confirmed
          type: 'dashed',
          width: 1
        },
        label: {
          formatter: 'Confirmed Time',
          position: 'middle',
          backgroundColor: 'rgba(255, 255, 255, 0.8)',
          padding: [4, 8],
          borderRadius: 4,
          color: '#52c41a'
        }
      });
    }

    // X-Axis MarkLine (Closed Time)
    if (closedDateStr) {
      markLineData.push({
        xAxis: closedDateStr,
        name: 'Closed Time',
        lineStyle: {
          color: '#f5222d', // Red for closed
          type: 'dashed',
          width: 1
        },
        label: {
          formatter: 'Closed Time',
          position: 'end',
          backgroundColor: 'rgba(255, 255, 255, 0.8)',
          padding: [4, 8],
          borderRadius: 4,
          color: '#f5222d'
        }
      });
    }

    const option: echarts.EChartsOption = {
      title: {
        text: 'K线走势',
        left: 0,
        textStyle: {
            fontSize: 16,
            fontWeight: 'normal',
            color: '#333'
        }
      },
      tooltip: {
        trigger: 'axis',
        axisPointer: {
          type: 'cross'
        },
        backgroundColor: 'rgba(255, 255, 255, 0.95)', // Transparent background
        borderColor: '#ccc',
        borderWidth: 1,
        textStyle: {
            color: '#333'
        },
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        formatter: (params: any) => { // Using 'any' due to complex ECharts type definitions
          const klineParams = params.find((p: any) => p.seriesName === 'KLine');
          const volParams = params.find((p: any) => p.seriesName === 'Volume');
          
          if (!klineParams || !klineParams.value) return '';
          
          const [o, c, l, h] = klineParams.value.slice(1);
          // volume is at index 5 in the original array, but ECharts passes the full array in .value
          // Our data structure: [open, close, low, high, volume, sign]
          const v = volParams ? volParams.value : klineParams.value[5];
          
          const date = klineParams.name;
          return `
            <div style="font-size: 12px; line-height: 1.5;">
              <div style="margin-bottom: 4px; color: #666;">${date}</div>
              <span style="color: #666;">开盘:</span> <span style="font-weight: bold;">${o}</span><br/>
              <span style="color: #666;">收盘:</span> <span style="font-weight: bold;">${c}</span><br/>
              <span style="color: #666;">最低:</span> <span style="font-weight: bold;">${l}</span><br/>
              <span style="color: #666;">最高:</span> <span style="font-weight: bold;">${h}</span><br/>
              <span style="color: #666;">成交量:</span> <span style="font-weight: bold;">${v}</span>
            </div>
          `;
        },
      },
      legend: {
        data: ['KLine', 'Volume'],
        bottom: 0
      },
      grid: [
        {
          left: '5%',
          right: '5%',
          height: '60%', // Top 60%
          top: '10%',
          containLabel: true
        },
        {
          left: '5%',
          right: '5%',
          top: '75%', // Start at 75%
          height: '15%', // 15% height
          containLabel: true
        }
      ],
      xAxis: [
        {
          type: 'category',
          data: dates,
          boundaryGap: false,
          axisLine: { onZero: false, lineStyle: { color: '#e0e0e0' } },
          splitLine: { show: false },
          axisLabel: { show: false }, // Hide labels for top axis
          gridIndex: 0
        },
        {
          type: 'category',
          data: dates,
          boundaryGap: false,
          axisLine: { onZero: false, lineStyle: { color: '#e0e0e0' } },
          splitLine: { show: false },
          axisLabel: { color: '#888', rotate: 30 },
          gridIndex: 1
        }
      ],
      yAxis: [
        {
          scale: true,
          splitArea: {
            show: true,
            areaStyle: {
                color: ['rgba(250,250,250,0.3)','rgba(210,210,210,0.3)']
            }
          },
          splitLine: {
              lineStyle: {
                  color: '#f0f0f0'
              }
          },
          axisLabel: {
              color: '#888'
          },
          min: minPrice,
          max: maxPrice,
          gridIndex: 0
        },
        {
          scale: true,
          splitLine: { show: false },
          axisLabel: { show: false },
          axisTick: { show: false },
          gridIndex: 1
        }
      ],
      dataZoom: [
        {
          type: 'inside',
          xAxisIndex: [0, 1],
          start: 0,
          end: 100
        },
        {
          show: true,
          xAxisIndex: [0, 1],
          type: 'slider',
          height: 20,
          bottom: 10,
          textStyle: {
            color: '#333'
          },
          fillerColor: 'rgba(22, 119, 255, 0.2)'
        }
      ],
      series: [
        {
          name: 'KLine',
          type: 'candlestick',
          data: data,
          itemStyle: {
            color: profitColor,
            color0: lossColor,
            borderColor: profitColor,
            borderColor0: lossColor
          },
          markLine: {
            symbol: ['none', 'none'],
            data: markLineData,
            animation: false 
          },
          xAxisIndex: 0,
          yAxisIndex: 0
        },
        {
          name: 'Volume',
          type: 'bar',
          xAxisIndex: 1,
          yAxisIndex: 1,
          data: data.map((item) => {
              return {
                  value: item[4], // Volume
                  itemStyle: {
                      color: item[5] > 0 ? profitColor : lossColor // Use sign to color
                  }
              }
          })
        }
      ]
    };

    chartInstance.current.setOption(option, true);

    const resizeObserver = new ResizeObserver(() => {
      chartInstance.current?.resize();
    });
    
    if (chartRef.current) {
      resizeObserver.observe(chartRef.current);
    }

    return () => {
      resizeObserver.disconnect();
      if (chartInstance.current && !chartInstance.current.isDisposed()) {
        chartInstance.current.dispose();
        chartInstance.current = null;
      }
    };
  }, [klines, signalPrice, signalType, signalTime, confirmedAt, closedAt]);

  return (
    <div style={{ position: 'relative' }}>
      {onIntervalChange && (
        <div style={{ position: 'absolute', right: 10, top: 0, zIndex: 10 }}>
          <Radio.Group 
            value={interval} 
            onChange={(e) => onIntervalChange(e.target.value)}
            size="small"
            buttonStyle="solid"
          >
            <Radio.Button value="15m">15m</Radio.Button>
            <Radio.Button value="1h">1h</Radio.Button>
            <Radio.Button value="4h">4h</Radio.Button>
            <Radio.Button value="1d">1d</Radio.Button>
          </Radio.Group>
        </div>
      )}
      <div ref={chartRef} style={{ height: '450px', width: '100%' }} />
    </div>
  );
}