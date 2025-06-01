import type { IAuthData, IAuthUserData } from '../types/auth.types';

export const setAuthUserData = async (userData: IAuthUserData) => {
    const authData = useState<IAuthData>('authData');
    authData.value.isAuth = true;
    authData.value.userData = userData;
};
