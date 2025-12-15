// 根据收益率返回颜色
export function getProfitColor(changePercent: number | string): string {
  const value = typeof changePercent === 'string'
    ? parseFloat(changePercent)
    : changePercent;

  if (value > 0) return '#52c41a'; // 绿色
  if (value < 0) return '#f5222d'; // 红色
  return '#8c8c8c'; // 灰色
}

// 状态标签颜色
export function getStatusColor(status: string): string {
  const colors: Record<string, string> = {
    PENDING: 'orange',
    CONFIRMED: 'blue',
    TRACKING: 'cyan',
    CLOSED: 'green',
    INVALIDATED: 'red',
  };
  return colors[status] || 'default';
}

// 信号类型颜色
export function getSignalTypeColor(type: string): string {
  return type === 'LONG' ? '#52c41a' : '#f5222d';
}
