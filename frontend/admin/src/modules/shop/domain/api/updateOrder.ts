import { tryToCatchApiErrors } from '~/shared/errors/errors';
import type { IOrderItem } from '../model/types/order';

interface Request {
    details: {
        client_name: string;
        client_surname: string;
        client_email: string;
        client_phone: string;
        delivery_address: string;
    };
    products: {
        id: number;
        quantity: number;
        price: number;
    }[];
}

const mapDataToRequest = (data: IOrderItem): Request => {
    const reqData: Request = {
        details: {
            client_name: data.details.client_name,
            client_surname: data.details.client_surname,
            client_email: data.details.client_email,
            client_phone: data.details.client_phone,
            delivery_address: data.details.delivery_address,
        },
        products: data.products.map((item) => ({
            id: item.id,
            quantity: item.quantity,
            price: item.price,
        })),
    };

    return reqData;
};

export async function updateOrder(data: IOrderItem) {
    try {
        await useNuxtApp().$apiFetch(`/orders/${data.id}`, {
            method: 'PUT',
            body: mapDataToRequest(data),
        });
    } catch (e: unknown) {
        throw tryToCatchApiErrors(e);
    }
}
