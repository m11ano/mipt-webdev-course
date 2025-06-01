import { th } from '@nuxt/ui/runtime/locale/index.js';
import { doLogout } from './auth/model';
import { getAuthToken } from './auth/model/actions/authToken';

export default defineNuxtPlugin({
    name: 'api-fetch',
    enforce: 'pre',
    async setup(nuxtApp) {
        const apiFetch = $fetch.create({
            baseURL: nuxtApp.$config.public.apiBase,
            retry: 3,
            retryStatusCodes: [500, 503],
            retryDelay: 500,

            onRequest({ options }) {
                options.headers.set('Authorization', 'Bearer ' + getAuthToken());
            },

            async onResponseError({ response, error }) {
                if (response.status === 401) {
                    doLogout();
                }
            },
        });

        return {
            provide: {
                apiFetch,
            },
        };
    },
});
