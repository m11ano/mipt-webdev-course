import { openSuccessAddToBasket } from '../../modals/SuccessAddToBasket/hooks';
import { useCartStore } from '../store/cart';
import type { IProductItem } from '../types/product';

export function addProductToCart(id: number) {
    openSuccessAddToBasket(id, (product: IProductItem) => {
        const cartStore = useCartStore();

        if (cartStore.getProductQuantity(id) >= product.stock_available) {
            return false;
        }

        cartStore.add(id);
        return true;
    });
}
