import { tryToCatchApiErrors } from '~/shared/errors/errors';
import type { IAuthUserData } from '../types/auth.types';
import { setAuthUserData } from './setAuthData';
import { getAuthToken } from './authToken';

export const doAuth = async (): Promise<void> => {
    try {
        const result = await $fetch<IAuthUserData>('auth', {
            method: 'POST',
            baseURL: useNuxtApp().$config.public.apiBase,
            onRequest({ options }) {
                const token = getAuthToken();
                if (token) {
                    options.headers.set('Authorization', 'Bearer ' + token);
                }
            },
        });

        setAuthUserData(result);
    } catch (e: unknown) {
        throw tryToCatchApiErrors(e);
    }
};
