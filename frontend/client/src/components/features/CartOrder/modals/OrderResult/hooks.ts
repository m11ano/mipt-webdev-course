import { useModal } from 'vue-final-modal';
import SuccessAddToBasket from './OrderResult.vue';

interface IProps {
    orderID: MaybeRef<number>;
    link: MaybeRef<string>;
    isLoaded: MaybeRef<boolean>;
}

export const useOrderResult = (props: IProps) => {
    const modal = useModal({
        component: SuccessAddToBasket,
        attrs: {
            modalObj: () => modal,
            orderID: props.orderID,
            link: props.link,
            isLoaded: props.isLoaded,
        },
    });

    return modal;
};
