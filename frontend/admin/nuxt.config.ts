// https://nuxt.com/docs/api/configuration/nuxt-config

import type { NuxtConfig } from 'nuxt/schema';
import type { ConfigLayerMeta, InputConfig } from 'c12';
import adminModules from './admin-modules.config';

const adminModulesConfigComponents: InputConfig<NuxtConfig, ConfigLayerMeta>['components'] = [];

adminModules.forEach((module) => {
    adminModulesConfigComponents.push(
        {
            path: `~/modules/${module.name}/components/shared`,
            extensions: ['.vue'],
            prefix: `${module.prefixName}Shared`,
        },
        {
            path: `~/modules/${module.name}/components/features`,
            extensions: ['.vue'],
            prefix: `${module.prefixName}Feature`,
        },
        {
            path: `~/modules/${module.name}/components/entities`,
            extensions: ['.vue'],
            prefix: `${module.prefixName}Entity`,
        },
        {
            path: `~/modules/${module.name}/components/widgets`,
            extensions: ['.vue'],
            prefix: `${module.prefixName}Widget`,
        },
        {
            path: `~/modules/${module.name}/components/pages`,
            extensions: ['.vue'],
            prefix: `${module.prefixName}Page`,
        },
    );
});

const config: InputConfig<NuxtConfig, ConfigLayerMeta> = {
    compatibilityDate: '2024-10-27',
    devtools: { enabled: false },
    $production: {
        nitro: {
            compressPublicAssets: true,
        },
    },
    devServer: {
        port: Number(process.env.NUXT_PORT) || 3001,
    },
    ssr: false,
    srcDir: 'src',
    alias: {
        '@styles': '/assets/styles',
    },
    dir: {
        pages: 'routes',
    },
    imports: {
        dirs: ['shared/hooks/**/*.ts'],
    },
    pinia: {
        storesDirs: ['shared/stores/**'],
    },
    modules: ['@nuxt/eslint', '@pinia/nuxt', '@nuxt/ui', '@vueuse/nuxt'],
    plugins: ['~/plugins/api-fetch.plugin.ts', '~/plugins/auth/auth.plugin.ts', '~/plugins/theme/theme.plugin.ts', '~/plugins/app/app.plugin.ts'],
    components: [
        {
            path: '~/core/components/shared',
            extensions: ['.vue'],
            prefix: 'Shared',
        },
        {
            path: '~/core/components/features',
            extensions: ['.vue'],
            prefix: 'Feature',
        },
        {
            path: '~/core/components/entities',
            extensions: ['.vue'],
            prefix: 'Entity',
        },
        {
            path: '~/core/components/widgets',
            extensions: ['.vue'],
            prefix: 'Widget',
        },
        {
            path: '~/core/components/pages',
            extensions: ['.vue'],
            prefix: 'Page',
        },
        ...adminModulesConfigComponents,
    ],
    css: ['@/assets/styles/index.less', '@/assets/styles/nuxt-ui.css'],
    app: {
        baseURL: '/admin/',
        head: {
            viewport: 'width=device-width, initial-scale=1, maximum-scale=1, minimum-scale=1, user-scalable=no',
        },
    },
    runtimeConfig: {
        public: {
            apiBase: '',
        },
    },
    ui: {
        theme: {
            colors: ['primary', 'secondary', 'graylight', 'info', 'success', 'warning', 'error'],
        },
    },
};

export default defineNuxtConfig(config);
