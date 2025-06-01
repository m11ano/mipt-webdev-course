import type { IMenu } from '~/core/model/types/type';
import { ShopModule } from '../const';

export function getMenu(): IMenu[] {
    return [
        {
            name: 'Магазин',
            icon: ShopModule.icon,
            to: `/${ShopModule.urlName}`,
            menuSel: ShopModule.urlName,
            defaultOpen: true,
            kids: [
                {
                    name: 'Товары',
                    to: `/${ShopModule.urlName}/products`,
                    subMenuSel: 'products',
                },
                {
                    name: 'Заказы',
                    to: `/${ShopModule.urlName}/orders`,
                    subMenuSel: 'orders',
                },
            ],
        },
    ];
}
