<script setup lang="ts">
import { setMenu } from '~/plugins/app/model/actions/setMenu';
import { setModuleBreadcrums } from '~/modules/shop/domain/actions/setModuleBreadcrums';
import { ShopModule } from '~/modules/shop/const';
import type { IProductItem } from '~/modules/shop/domain/model/types/product';
import { checkProductModel } from '~/modules/shop/domain/hooks/checkProductModel';
import { StandartErrorList } from '~/shared/errors/errors';
import { fetchProduct } from '~/modules/shop/domain/api/fetchProduct';
import { updateProduct } from '~/modules/shop/domain/api/updateProduct';
import ProductStock from '../../shared/modals/ProductStock.vue';

const props = defineProps<{
    id: number;
}>();

useSeoMeta({
    title: 'Редактирование товара',
});

setMenu(ShopModule.urlName, 'products');

setModuleBreadcrums([
    {
        name: 'Товары',
        to: '/products',
    },
    {
        name: 'Редактирование товара',
    },
]);

const productModel = ref<IProductItem | null>(null);

const isLoading = ref(false);

const errors = ref<string[]>([]);

const toast = useToast();

watch(
    () => props.id,
    async () => {
        isLoading.value = true;
        try {
            const data = await fetchProduct(props.id);
            productModel.value = data;
        } catch (e) {
            if (e instanceof StandartErrorList) {
                if (e.code === 404) {
                    showError({
                        statusCode: e.code,
                        statusMessage: 'Товар не найден',
                    });
                }
                errors.value = e.details;
            }
        } finally {
            isLoading.value = false;
        }
    },
    {
        immediate: true,
    },
);

const save = async () => {
    if (isLoading.value || !productModel.value) return;

    errors.value = checkProductModel(productModel.value);

    if (errors.value.length) return;

    isLoading.value = true;
    try {
        await updateProduct(productModel.value);

        toast.add({
            title: 'Успех',
            description: 'Товар сохранен',
            color: 'success',
            icon: 'i-lucide-check-circle',
        });
    } catch (e) {
        if (e instanceof StandartErrorList) {
            errors.value = e.details;
        }
    } finally {
        isLoading.value = false;
    }
};

const updateStaticValues = async () => {
    if (!productModel.value) return;
    try {
        const data = await fetchProduct(props.id);
        productModel.value.stock_available = data.stock_available;
    } catch (e) {
        //
    }
};

const overlay = useOverlay();
const changeStock = async () => {
    if (!productModel.value) return;

    const modal = overlay.create(ProductStock, {
        props: {
            productID: productModel.value.id,
        },
        destroyOnClose: true,
    });

    const instance = modal.open();

    const shouldUpdate = await instance.result;
    if (shouldUpdate) {
        updateStaticValues();
    }
};
</script>

<template>
    <div>
        <div>
            <div class="form-table">
                <div>
                    <div class="title">ID:</div>
                    <div class="value">
                        {{ productModel?.id }}
                    </div>
                </div>
                <div>
                    <div class="title">URL:</div>
                    <div class="value">
                        <a
                            :href="`/product-${productModel?.id}`"
                            target="_blank"
                            style="text-decoration: underline"
                            >/product-{{ productModel?.id }}</a
                        >
                    </div>
                </div>
                <div>
                    <div class="title">Доступный остаток:</div>
                    <div class="value">
                        <div class="flex items-center gap-4">
                            <div>{{ productModel?.stock_available }} шт.</div>

                            <UButton
                                color="graylight"
                                @click="changeStock"
                                >Изменить остаток</UButton
                            >
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <div class="form_title mt-10">
            <div class="title">Карточка товара</div>
            <div class="buttons">
                <UButton
                    :disabled="isLoading"
                    :loading="isLoading"
                    @click="save"
                    >Сохранить</UButton
                >
            </div>
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
                v-if="productModel"
                v-model="productModel"
                mode="edit"
                :disabled="isLoading"
            />
        </div>
    </div>
</template>

<style lang="less" module>
@import '@styles/includes';
</style>
