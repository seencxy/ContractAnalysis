import numeral from 'numeral';
import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';
import 'dayjs/locale/zh-cn';

dayjs.extend(relativeTime);
dayjs.locale('zh-cn');

// 格式化数字
export function formatNumber(value: number | string, format = '0,0.00'): string {
  const num = typeof value === 'string' ? parseFloat(value) : value;
  return numeral(num).format(format);
}

// 格式化百分比
export function formatPercent(value: number | string): string {
  const num = typeof value === 'string' ? parseFloat(value) : value;
  return numeral(num / 100).format('0.00%');
}

// 直接格式化百分比字符串(后端已经是百分比)
export function formatPercentString(value: string | undefined): string {
  if (!value) return '--';
  const num = parseFloat(value);
  return numeral(num).format('0.00') + '%';
}

// 格式化时间
export function formatTime(time: string, format = 'YYYY-MM-DD HH:mm:ss'): string {
  return dayjs(time).format(format);
}

// 格式化相对时间
export function formatRelativeTime(time: string): string {
  return dayjs(time).fromNow();
}

// 格式化价格 (5位小数)
export function formatPrice(price: number | string): string {
  const num = typeof price === 'string' ? parseFloat(price) : price;
  return numeral(num).format('0,0.00000');
}
