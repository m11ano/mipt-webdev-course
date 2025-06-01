import { clearAuthUserData } from './clearAuthData';
import { clearAuthToken } from './authToken';

export const doLogout = async (): Promise<void> => {
    clearAuthUserData();
    clearAuthToken();
};
