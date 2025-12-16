import React from 'react';
import { Tag, Tooltip } from 'antd';
import { 
  ThunderboltOutlined, 
  UsergroupDeleteOutlined, 
  RiseOutlined 
} from '@ant-design/icons';
import type { Signal } from '../../types/signal';

interface StrategyBadgeProps {
  signal: Signal;
  showDetails?: boolean;
}

const STRATEGY_CONFIG: Record<string, { color: string; icon: React.ReactNode; label: string }> = {
  'Minority Strategy': {
    color: 'orange',
    icon: <UsergroupDeleteOutlined />,
    label: 'Minority',
  },
  'Whale Strategy': {
    color: 'blue',
    icon: <RiseOutlined />,
    label: 'Whale',
  },
  'Smart Money Strategy': {
    color: 'purple',
    icon: <ThunderboltOutlined />,
    label: 'Smart Money',
  },
};

const StrategyBadge: React.FC<StrategyBadgeProps> = ({ signal, showDetails = false }) => {
  const config = STRATEGY_CONFIG[signal.strategy_name] || {
    color: 'default',
    icon: null,
    label: signal.strategy_name,
  };

  const getTooltipContent = () => {
    if (!signal.strategy_context) return signal.reason || 'No details available';

    const ctx = signal.strategy_context;
    
    switch (signal.strategy_name) {
      case 'Minority Strategy':
        return (
          <div>
            <div><b>Trigger:</b> Extreme Ratio</div>
            <div>Long/Short: {(ctx.generate_short_when_long_ratio_above as number)?.toFixed(0)}% / {(ctx.generate_long_when_short_ratio_above as number)?.toFixed(0)}%</div>
          </div>
        );
      case 'Whale Strategy':
        return (
          <div>
            <div><b>Trigger:</b> Whale Divergence</div>
            <div>Min Divergence: {(ctx.min_divergence as number)?.toFixed(2)}%</div>
            <div>Whale Position: &gt;{(ctx.whale_position_threshold as number)?.toFixed(0)}%</div>
          </div>
        );
      case 'Smart Money Strategy':
        return (
          <div>
            <div><b>Trigger:</b> Confluence</div>
            <div>Pattern: {ctx.setup_type || 'Unknown'}</div>
            <div>Funding Rate: Positive</div>
          </div>
        );
      default:
        return signal.reason;
    }
  };

  const badge = (
    <Tag color={config.color} icon={config.icon} style={{ marginRight: 0 }}>
      {config.label}
    </Tag>
  );

  if (showDetails) {
    return (
      <Tooltip title={getTooltipContent()}>
        {badge}
      </Tooltip>
    );
  }

  return badge;
};

export default StrategyBadge;
