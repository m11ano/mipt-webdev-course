import { tryToCatchApiErrors } from '~/shared/errors/errors';
import type { OrderStatus } from '../model/types/order';

export async function setOrderStatus(id: number, status: OrderStatus) {
    try {
        return await useNuxtApp().$apiFetch<Response>(`/orders/${id}/status`, {
            method: 'PUT',
            body: {
                status,
            },
        });
    } catch (e: unknown) {
        throw tryToCatchApiErrors(e);
    }
}
