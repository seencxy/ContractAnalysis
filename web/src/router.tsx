import { createBrowserRouter } from 'react-router-dom';
import { MainLayout } from '@/components/layout/MainLayout';
import Dashboard from '@/pages/Dashboard';
import Signals from '@/pages/Signals';
import Analysis from '@/pages/Analysis';
import Ranking from '@/pages/Ranking';
import History from '@/pages/History';
import StatisticsHistory from '@/pages/StatisticsHistory';

export const router = createBrowserRouter([
  {
    path: '/',
    element: <MainLayout />,
    children: [
      { index: true, element: <Dashboard /> },
      { path: 'signals', element: <Signals /> },
      { path: 'analysis', element: <Analysis /> },
      { path: 'ranking', element: <Ranking /> },
      { path: 'history', element: <History /> },
      { path: 'statistics-history', element: <StatisticsHistory /> },
    ],
  },
], {
  future: {
    v7_startTransition: true,
    v7_relativeSplatPath: true,
    v7_fetcherPersist: true,
    v7_normalizeFormMethod: true,
    v7_partialHydration: true,
    v7_skipActionErrorRevalidation: true,
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
  } as any
});
