export interface IOrderDetailsFormData {
    client_name: string;
    client_surname: string;
    client_email: string;
    client_phone: string;
    delivery_address: string;
}

export interface IOrderFormData {
    products: {
        id: number;
        quantity: number;
    }[];
    details: IOrderDetailsFormData;
}
