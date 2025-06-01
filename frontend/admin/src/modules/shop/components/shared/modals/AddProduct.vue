<script setup lang="ts">
import { fetchProduct } from '~/modules/shop/domain/api/fetchProduct';
import { StandartErrorList } from '~/shared/errors/errors';
import type { IProductItem } from '~/modules/shop/domain/model/types/product';

const props = defineProps<{
    orderID: number;
}>();

const emit = defineEmits<{ close: [IProductItem | null] }>();

const productID = ref('');

const errors = ref<string[]>([]);

const isLoading = ref(false);

const save = async () => {
    if (isLoading.value) return;

    errors.value = [];

    if (productID.value.length < 1) {
        errors.value.push('ID не указан');
        return;
    }

    isLoading.value = true;
    try {
        const product = await fetchProduct(Number(productID.value));
        if (product.stock_available < 1) {
            errors.value.push('Товар закончился');
            return;
        }
        if (!product.is_published) {
            errors.value.push('Товар не опубликован');
            return;
        }
        emit('close', product);
    } catch (e) {
        if (e instanceof StandartErrorList) {
            if (e.code === 404) {
                errors.value = ['Товар не найден'];
            } else {
                errors.value = e.details;
            }
        }
    } finally {
        isLoading.value = false;
    }
};
</script>

<template>
    <UModal
        :title="`Добавление товара в заказ`"
        @close="emit('close', null)"
    >
        <template #body>
            <div :class="$style.wrapper">
                <div>Укажите ID товара:</div>
                <div>
                    <UInput
                        v-model="productID"
                        size="xl"
                    />
                </div>
                <div v-if="errors.length">
                    <UAlert
                        title="Возникли ошибки!"
                        icon="i-lucide-ban"
                    >
                        <template #description>
                            <template
                                v-for="error in errors"
                                :key="error"
                            >
                                <div>– {{ error }}</div>
                            </template>
                        </template>
                    </UAlert>
                </div>
            </div>
        </template>
        <template #footer>
            <div class="flex justify-end">
                <UButton
                    label="Добавить"
                    color="primary"
                    :loading="isLoading"
                    :disabled="isLoading"
                    @click="save"
                />
            </div>
        </template>
    </UModal>
</template>

<style lang="less" module>
@import '@styles/includes';

.wrapper {
    display: flex;
    flex-direction: column;
    gap: 20px;
}
</style>
