import { tryToCatchApiErrors } from '~/shared/errors/errors';
import type { IOrderItem } from '../model/types/order';

export const fetchOrder = async (id: number) => {
    try {
        return await useNuxtApp().$apiFetch<IOrderItem>(`/orders/${id}`);
    } catch (e: unknown) {
        throw tryToCatchApiErrors(e);
    }
};
