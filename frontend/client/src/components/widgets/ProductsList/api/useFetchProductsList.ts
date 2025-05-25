import type { IProductListItem } from '~/domain/shop/model/types/product';

const DEFAULT_LIMIT = 20;

export const useLoadProductsList = (page?: number, limit?: number) => {
    if (!page) page = 1;
    if (!limit) limit = DEFAULT_LIMIT;

    return useNuxtApp().$apiFetch<{ items: IProductListItem[]; total: number }>('/products', {
        params: {
            offset: (page - 1) * limit,
            limit,
        },
    });
};

export const useFetchProductsList = (page?: number, limit?: number, lazy: boolean = true, immediate: boolean = true) => {
    if (!page) page = 1;
    if (!limit) limit = DEFAULT_LIMIT;

    return useAPIFetch<{ items: IProductListItem[]; total: number }>('/products', {
        params: {
            offset: (page - 1) * limit,
            limit: limit,
        },
        lazy,
        immediate,
    });
};
