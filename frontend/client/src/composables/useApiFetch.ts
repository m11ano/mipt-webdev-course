import type { UseFetchOptions } from 'nuxt/app';

export function useAPIFetch<T = unknown>(url: string | (() => string), options?: UseFetchOptions<T>) {
    return useFetch(url, {
        ...options,
        $fetch: useNuxtApp().$apiFetch,
    });
}
