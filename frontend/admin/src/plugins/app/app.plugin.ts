import type { IAppData } from './model/types/types';

export default defineNuxtPlugin({
    name: 'appProvider',
    enforce: 'pre',
    async setup() {
        const appProvider = useState<IAppData>('app', () =>
            reactive({
                title: 'Панель управления',
                breadcrumbs: [],
                menuSel: '',
                subMenuSel: '',
            }),
        );

        return {
            provide: {
                appProvider: appProvider.value,
            },
        };
    },
    env: {
        // Set this value to `false` if you don't want the plugin to run when rendering server-only or island components.
        islands: true,
    },
});
