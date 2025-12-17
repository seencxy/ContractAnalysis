import { useQuery } from '@tanstack/react-query';
import { strategiesApi } from '@/api/endpoints/strategies';
import type { Strategy } from '@/types/strategy';

export const useStrategies = () => {
  return useQuery({
    queryKey: ['strategies'],
    queryFn: async () => {
      const response = await strategiesApi.getStrategies();
      // Interceptor returns response.data (the body), but TS thinks it's AxiosResponse.
      // At runtime, response is { data: Strategy[] }. We want to return Strategy[].
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      return (response as any).data as Strategy[];
    },
    staleTime: 1000 * 60 * 60, // Cache for 1 hour (strategies rarely change)
  });
};
