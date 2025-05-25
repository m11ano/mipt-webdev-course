enum OrderStatus {
    New = 'new',
    Created = 'created',
    InWork = 'in_work',
    Finished = 'finished',
    Canceled = 'canceled',
}

export const OrderStatusText = {
    [OrderStatus.New]: 'Новый',
    [OrderStatus.Created]: 'Создан',
    [OrderStatus.InWork]: 'В работе',
    [OrderStatus.Finished]: 'Выполнен',
    [OrderStatus.Canceled]: 'Отменен',
};

export interface IOrder {
    id: number;
    secret_key: string;
    order_sum: number;
    status: OrderStatus;
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
