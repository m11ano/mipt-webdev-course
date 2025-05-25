import { tryToThrowApiErrors } from '~/shared/errors/errors';
import type { IOrderFormData } from '../model/types/types';

interface Response {
    id: number;
    secret_key: string;
}

export async function fetchCreateOrder(data: IOrderFormData): Promise<Response> {
    try {
        return await useNuxtApp().$apiFetch<Response>('/orders', {
            method: 'POST',
            body: data,
        });
    } catch (e: unknown) {
        throw tryToThrowApiErrors(e);
    }
}
