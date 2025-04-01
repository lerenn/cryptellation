import { defineConfig } from '@hey-api/openapi-ts';

export default defineConfig({
  input: '../../../api/gateway/v1.yaml',
  output: 'src/generated',
  plugins: [{
    name: '@hey-api/client-fetch',
    runtimeConfigPath: './src/config.ts',
    }],
});
