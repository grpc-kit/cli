import type { RouteItem } from '@grpc-kit/adm';
import ExampleDemo from './pages/example/demo';

export const businessRoutes: RouteItem[] = [
  { path: '/example/demo', element: <ExampleDemo /> },
];
