import { useModal } from 'vue-final-modal';
import SuccessAddToBasket from './SuccessAddToBasket.vue';
import type { IProductItem } from '../../model/types/product';

interface IProps {
    title: MaybeRef<string>;
    productID: number;
    checkAvailability: (product: IProductItem) => boolean;
}

export const useSuccessAddToBasket = (props: IProps) => {
    const modal = useModal({
        component: SuccessAddToBasket,
        attrs: {
            title: props.title,
            modalObj: () => modal,
            productID: props.productID,
            checkAvailability: props.checkAvailability,
        },
    });

    return modal;
};

export const openSuccessAddToBasket = (productID: number, checkAvailability: (product: IProductItem) => boolean) => {
    const modal = useSuccessAddToBasket({ title: 'Товар добавлен в корзину!', productID: productID, checkAvailability });
    modal.open();
    return modal;
};
