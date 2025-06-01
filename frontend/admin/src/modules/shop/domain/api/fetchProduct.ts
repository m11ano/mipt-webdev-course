import type { IProductItem } from '~/modules/shop/domain/model/types/product';
import { tryToCatchApiErrors } from '~/shared/errors/errors';

export const fetchProduct = async (id: number) => {
    try {
        return await useNuxtApp().$apiFetch<IProductItem>(`/products/${id}`);
    } catch (e: unknown) {
        throw tryToCatchApiErrors(e);
    }
};
