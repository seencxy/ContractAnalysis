import client from '../client';
import type { Strategy } from '@/types/strategy';

export const strategiesApi = {
  getStrategies: () => client.get<{ data: Strategy[] }>('/strategies'),
};
