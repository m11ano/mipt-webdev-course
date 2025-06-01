import { tryToCatchApiErrors } from '~/shared/errors/errors';
import type { IOrderListItem } from '../model/types/order';

const DEFAULT_LIMIT = 20;

export const fetchOrdersList = async (page?: number, limit?: number) => {
    if (!page) page = 1;
    if (!limit) limit = DEFAULT_LIMIT;

    try {
        return await useNuxtApp().$apiFetch<{ items: IOrderListItem[]; total: number }>('/orders', {
            params: {
                offset: (page - 1) * limit,
                limit,
            },
        });
    } catch (e: unknown) {
        throw tryToCatchApiErrors(e);
    }
};
