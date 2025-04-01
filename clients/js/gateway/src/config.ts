import type { CreateClientConfig } from './generated/client.gen';

export const createClientConfig: CreateClientConfig = (config) => ({
  ...config,
  baseUrl: 'http://localhost:8080',
});