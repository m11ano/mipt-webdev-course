import { tryToCatchApiErrors } from '~/shared/errors/errors';

export async function deleteProduct(id: number) {
    try {
        return await useNuxtApp().$apiFetch<Response>(`/products/${id}`, {
            method: 'DELETE',
        });
    } catch (e: unknown) {
        throw tryToCatchApiErrors(e);
    }
}
