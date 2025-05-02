<script setup lang="ts">
import type { UseModalReturnType } from 'vue-final-modal';
import type { IProductItem } from '../../model/types/product';
import { mockProductsList } from '~/mocks/products';
import { coolNumber } from '~/shared/helpers/functions';

const props = defineProps<{
    title: MaybeRef<string>;
    modalObj: () => UseModalReturnType<any>;
    productID: number;
}>();

const isLoading = ref(true);

const product = ref<IProductItem | null>(null);

onMounted(() => {
    setTimeout(() => {
        product.value = mockProductsList.find((item) => item.id === props.productID) || null;
        isLoading.value = false;
    }, 1000);
});
</script>

<template>
    <SharedModalsDefaultModal :modal-obj="props.modalObj()">
        <template v-if="isLoading"><SharedModalsDefaultLoader /></template>
        <template v-else>
            <SharedModalsDefaultWrapper
                :modal-obj="props.modalObj()"
                :class-name="$style.box_wrapper"
            >
                <template #title>{{ title }}</template>
                <template #content>
                    <div :class="$style.content">
                        <div :class="$style.img">
                            <img
                                v-if="product?.image_preview"
                                :src="product?.image_preview"
                                :alt="product?.name"
                                :title="product?.name"
                            />
                            <span></span>
                        </div>
                        <div :class="$style.text">
                            <div :class="$style.name">
                                {{ product?.name }}
                            </div>
                            <div :class="$style.price">{{ coolNumber(product?.price || 0) }} ₽</div>
                        </div>
                    </div>
                </template>
                <template #buttons>
                    <button
                        class="button_1 gray"
                        @click="modalObj().close()"
                    >
                        Продолжить покупки
                    </button>
                    <NuxtLink
                        to="/cart"
                        class="button_1"
                        @click="modalObj().close()"
                    >
                        Оформить заказ
                    </NuxtLink>
                </template>
            </SharedModalsDefaultWrapper>
        </template>
    </SharedModalsDefaultModal>
</template>

<style lang="less" module>
@import '@styles/includes';

.box_wrapper {
    width: 100%;
    max-width: 750px;
}

.content {
    display: flex;
    gap: 40px;

    .width-size-sm-less({
        gap: 20px;
        flex-wrap: wrap;
    });
}

.img {
    flex-shrink: 0;
    width: 50%;
    max-width: 250px;
    position: relative;
    background-color: var(--ui-bg-color);
    border-radius: 10px;
    overflow: hidden;
    box-shadow: 0px 0px 10px rgba(0, 0, 0, 0.1);
    transition: box-shadow 0.3s ease;

    .width-size-sm-less({
        width: 100%;
    });

    > img {
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        object-fit: cover;
    }

    > span {
        display: block;
        padding-top: 100%;
    }
}

.text {
    word-break: break-word;

    > .name {
        font-size: 24px;
        font-family: 'Strong';
    }

    > .price {
        font-weight: 600;
        font-size: 18px;
        margin-top: 20px;
    }
}
</style>
