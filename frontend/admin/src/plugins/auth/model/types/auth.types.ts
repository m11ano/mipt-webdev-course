export interface IAuthUserData {
    id: string;
    name: string;
    surname: string;
    email: string;
}

export interface IAuthData {
    isLoading: boolean;
    isAuth: boolean;
    userData: IAuthUserData | null;
}
