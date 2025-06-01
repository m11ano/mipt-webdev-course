import type { IProductListItem } from '~/modules/shop/domain/model/types/product';
import { tryToCatchApiErrors } from '~/shared/errors/errors';

const DEFAULT_LIMIT = 20;

export const fetchProductsList = async (page?: number, limit?: number) => {
    if (!page) page = 1;
    if (!limit) limit = DEFAULT_LIMIT;

    try {
        return await useNuxtApp().$apiFetch<{ items: IProductListItem[]; total: number }>('/products', {
            params: {
                offset: (page - 1) * limit,
                limit,
            },
        });
    } catch (e: unknown) {
        throw tryToCatchApiErrors(e);
    }
};
