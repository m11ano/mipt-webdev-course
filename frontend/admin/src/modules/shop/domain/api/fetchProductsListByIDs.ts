import type { IProductListItem } from '~/modules/shop/domain/model/types/product';
import { tryToCatchApiErrors } from '~/shared/errors/errors';

export const fetchProductsListByIDs = async (ids: number[]) => {
    try {
        return await useNuxtApp().$apiFetch<{ items: IProductListItem[]; total: number }>('/products', {
            params: {
                ids: ids.join(','),
            },
        });
    } catch (e: unknown) {
        throw tryToCatchApiErrors(e);
    }
};
