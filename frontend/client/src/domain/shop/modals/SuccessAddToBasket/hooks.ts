import { useModal } from 'vue-final-modal';
import SuccessAddToBasket from './SuccessAddToBasket.vue';

interface IProps {
    title: MaybeRef<string>;
    productID: number;
}

export const useSuccessAddToBasket = (props: IProps) => {
    const modal = useModal({
        component: SuccessAddToBasket,
        attrs: {
            title: props.title,
            modalObj: () => modal,
            productID: props.productID,
        },
    });

    return modal;
};

export const openSuccessAddToBasket = (productID: number) => {
    const modal = useSuccessAddToBasket({ title: 'Товар добавлен в корзину!', productID: productID });
    modal.open();
    return modal;
};
