export default defineNuxtPlugin((nuxtApp) => {
    const apiFetch = $fetch.create({
        baseURL: nuxtApp.$config.public.apiBase,
        retry: 3,
        retryStatusCodes: [500, 503],
        retryDelay: 500,
    });

    return {
        provide: {
            apiFetch,
        },
    };
});
