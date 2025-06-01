<script setup lang="ts">
import { setMenu } from '~/plugins/app/model/actions/setMenu';
import { setModuleBreadcrums } from '~/modules/shop/domain/actions/setModuleBreadcrums';
import { ShopModule } from '~/modules/shop/const';
import type { IProductItem } from '~/modules/shop/domain/model/types/product';
import { checkProductModel } from '~/modules/shop/domain/hooks/checkProductModel';
import { createProduct } from '~/modules/shop/domain/api/createProduct';
import { StandartErrorList } from '~/shared/errors/errors';

useSeoMeta({
    title: 'Создание товара',
});

setMenu(ShopModule.urlName, 'products');

setModuleBreadcrums([
    {
        name: 'Товары',
        to: '/products',
    },
    {
        name: 'Создание товара',
    },
]);

const productModel = ref<IProductItem>({
    id: 0,
    name: '',
    is_published: true,
    full_description: '',
    price: 0,
    stock_available: 0,
    image_preview: null,
    slider: [],
});

const isLoading = ref(false);

const errors = ref<string[]>([]);

const toast = useToast();

const save = async () => {
    if (isLoading.value) return;
    errors.value = checkProductModel(productModel.value);

    if (errors.value.length) return;

    isLoading.value = true;
    try {
        const result = await createProduct(productModel.value);

        toast.add({
            title: 'Успех',
            description: 'Товар создан',
            color: 'success',
            icon: 'i-lucide-check-circle',
        });

        navigateTo(`/${ShopModule.urlName}/products/${result.id}`);
    } catch (e) {
        if (e instanceof StandartErrorList) {
            errors.value = e.details;
        }
    } finally {
        isLoading.value = false;
    }
};
</script>

<template>
    <div>
        <div class="flex justify-end">
            <UButton
                :disabled="isLoading"
                :loading="isLoading"
                @click="save"
                >Создать</UButton
            >
        </div>
        <div
            v-if="errors.length"
            class="mt-4"
        >
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
        <div class="mt-4">
            <ShopWidgetProductForm
                v-model="productModel"
                mode="new"
                :disabled="isLoading"
            />
        </div>
    </div>
</template>

<style lang="less" module>
@import '@styles/includes';
</style>
