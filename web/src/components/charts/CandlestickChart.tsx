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
  interval?: string;
  onIntervalChange?: (interval: string) => void;
}

export function CandlestickChart({ klines, signalPrice, signalType = 'LONG', signalTime, interval = '15m', onIntervalChange }: CandlestickChartProps) {
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
    const data = sortedKlines.map((k) => [
      k.open,
      k.close,
      k.low,
      k.high,
      k.volume
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

    const markLineData: any[] = [];
    
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

    // X-Axis MarkLine (Time)
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
        formatter: (params: any) => {
          const klineParams = params[0];
          if (!klineParams || !klineParams.value) return '';
          
          const [o, c, l, h, v] = klineParams.value.slice(1);
          
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
        }
      },
      legend: {
        data: ['KLine'],
        bottom: 0
      },
      grid: {
        left: '5%',
        right: '5%',
        bottom: '10%',
        top: '15%',
        containLabel: true
      },
      xAxis: {
        type: 'category',
        data: dates,
        boundaryGap: false,
        axisLine: { onZero: false, lineStyle: { color: '#e0e0e0' } },
        splitLine: { show: false },
        axisLabel: {
            color: '#888',
            rotate: 30
        },
        min: 'dataMin',
        max: 'dataMax'
      },
      yAxis: {
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
        max: maxPrice
      },
      dataZoom: [
        {
          type: 'inside',
          start: 0,
          end: 100
        },
        {
          show: true,
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
            animation: false // Disable animation for markLines
          }
        }
      ]
    };

    chartInstance.current.setOption(option, true);

    const handleResize = () => {
      chartInstance.current?.resize();
    };

    window.addEventListener('resize', handleResize);

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
  }, [klines, signalPrice, signalType, signalTime]);

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
