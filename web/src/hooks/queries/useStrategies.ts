import { useQuery } from '@tanstack/react-query';
import { strategiesApi } from '@/api/endpoints/strategies';

export const useStrategies = () => {
  return useQuery({
    queryKey: ['strategies'],
    queryFn: async () => {
      const response = await strategiesApi.getStrategies();
      return response.data;
    },
    staleTime: 1000 * 60 * 60, // Cache for 1 hour (strategies rarely change)
  });
};
