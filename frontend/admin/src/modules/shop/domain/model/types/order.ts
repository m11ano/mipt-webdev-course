export enum OrderStatus {
    New = 'new',
    Created = 'created',
    InWork = 'in_work',
    Finished = 'finished',
    Canceled = 'canceled',
}

type OrderStatusParams = {
    title: string;
    variant?: 'solid' | 'outline' | 'soft' | 'subtle';
    color?: 'primary' | 'secondary' | 'graylight' | 'info' | 'success' | 'warning' | 'error' | 'neutral';
};

export const OrderStatusParams: Record<OrderStatus, OrderStatusParams> = {
    [OrderStatus.New]: {
        title: 'Новый',
        variant: 'subtle',
        color: 'success',
    },
    [OrderStatus.Created]: {
        title: 'Создан',
        variant: 'subtle',
        color: 'success',
    },
    [OrderStatus.InWork]: {
        title: 'В работе',
        variant: 'subtle',
        color: 'info',
    },
    [OrderStatus.Finished]: {
        title: 'Выполнен',
        variant: 'subtle',
        color: 'neutral',
    },
    [OrderStatus.Canceled]: {
        title: 'Отменен',
        variant: 'subtle',
        color: 'graylight',
    },
};

export interface IOrderListItem {
    details: {
        client_name: string;
        client_surname: string;
        client_phone: string;
        client_email: string;
        delivery_address: string;
    };
    id: number;
    order_sum: number;
    secret_key: string;
    status: OrderStatus;
}

export interface IOrderItem {
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
