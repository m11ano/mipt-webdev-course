import { tryToCatchApiErrors } from '~/shared/errors/errors';
import type { IAuthUserData } from '../types/auth.types';
import { setAuthToken } from './authToken';
import { setAuthUserData } from './setAuthData';

interface ILoginResponse {
    token: string;
    auth_user_data: IAuthUserData;
}

export const doLogin = async (email: string, password: string): Promise<void> => {
    try {
        const result = await $fetch<ILoginResponse>('auth/login', {
            baseURL: useNuxtApp().$config.public.apiBase,
            method: 'POST',
            body: { email, password },
        });

        setAuthToken(result.token);
        setAuthUserData(result.auth_user_data);
    } catch (e: unknown) {
        throw tryToCatchApiErrors(e);
    }
};
