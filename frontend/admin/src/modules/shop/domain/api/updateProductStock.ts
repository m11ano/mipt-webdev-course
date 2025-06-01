import { tryToCatchApiErrors } from '~/shared/errors/errors';

interface Request {
    operation: 'increase' | 'decrease';
    value: number;
}

export async function updateProductStock(id: number, data: Request) {
    try {
        await useNuxtApp().$apiFetch(`/products/${id}/stock`, {
            method: 'POST',
            body: data,
        });
    } catch (e: unknown) {
        throw tryToCatchApiErrors(e);
    }
}
