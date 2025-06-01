const AUTH_TOKEN_LOCAL_STORAGE_KEY = 'auth_token';

export const setAuthToken = (token: string) => {
    localStorage.setItem(AUTH_TOKEN_LOCAL_STORAGE_KEY, token);
};

export const clearAuthToken = () => {
    localStorage.removeItem(AUTH_TOKEN_LOCAL_STORAGE_KEY);
};

export const getAuthToken = () => {
    return localStorage.getItem(AUTH_TOKEN_LOCAL_STORAGE_KEY);
};
