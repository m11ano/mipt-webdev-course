import { setAppBreadcrumbs } from '~/plugins/app/model/actions/setAppBreadcrumbs';
import type { IAppDataBreadcrumb } from '~/plugins/app/model/types/types';
import { ShopModule } from '../../const';

export function setModuleBreadcrums(items: IAppDataBreadcrumb[]) {
    setAppBreadcrumbs([
        {
            name: 'Магазин',
            to: `/${ShopModule.urlName}`,
            icon: 'i-lucide-shopping-cart',
        },
        ...items.map((item) => ({
            name: item.name,
            to: item.to ? `/${ShopModule.urlName}${item.to}` : undefined,
            icon: item.icon,
        })),
    ]);
}
