import { createApp } from '@grpc-kit/adm';
import { businessRoutes } from './routes';

createApp({
  title: 'Your Service Admin',
  routes: businessRoutes,
});
