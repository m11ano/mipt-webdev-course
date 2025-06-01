import Confirm from '~/core/components/shared/Confirm/modals/Confirm.vue';
import { deleteProduct } from '../api/deleteProduct';

export async function removeProduct(id: number): Promise<boolean> {
    const modal = useOverlay().create(Confirm, {
        props: {
            text: 'Вы действительно хотите удалить товар?',
        },
        destroyOnClose: true,
    });

    const instance = modal.open();

    const shouldDelete = await instance.result;
    if (shouldDelete) {
        try {
            await deleteProduct(id);

            useToast().add({
                title: 'Успех',
                description: 'Товар удален',
                color: 'success',
                icon: 'i-lucide-check-circle',
            });

            return true;
        } catch (e) {
            useToast().add({
                title: 'Ошибка',
                description: 'Ошибка при удалении товара',
                color: 'error',
                icon: 'i-lucide-ban',
            });
        }
    }

    return false;
}
