import { openSuccessAddToBasket } from '../../modals/SuccessAddToBasket/hooks';
import { useCartStore } from '../store/cart';

export function addProductToCart(id: number) {
    const cartStore = useCartStore();
    cartStore.add(id);

    openSuccessAddToBasket(id);
}
