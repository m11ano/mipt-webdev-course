<script setup lang="ts">
import { useModalConfirm } from '~/components/shared/modals/Confirm/useModalConfirm';
import { useCartStore } from '~/domain/shop';
import type { IProductListItem } from '~/domain/shop/model/types/product';
import { coolNumber } from '~/shared/helpers/functions';

const props = defineProps<{ productsData: IProductListItem[] }>();

const cartStore = useCartStore();

interface ProductInCart {
    product: IProductListItem;
    quantity: number;
}

const products = computed<ProductInCart[]>(() => {
    const result: ProductInCart[] = [];
    cartStore.items.forEach((item) => {
        const product = props.productsData.find((i) => i.id === item.id);
        if (!product) return;

        result.push({
            product,
            quantity: item.quantity,
        });
    });

    return result;
});

const setQuantity = (id: number, quantity: number) => {
    if (isNaN(quantity) || quantity < 1) {
        quantity = 1;
    }
    const productsDataItem = props.productsData.find((i) => i.id === id);
    if (!productsDataItem) return;

    quantity = Math.min(quantity, productsDataItem.stock_available);
    cartStore.setQuantity(id, quantity);
};

const total = computed(() => {
    let total = 0;
    products.value.forEach((item) => {
        total += item.product.price * item.quantity;
    });
    return total;
});

const removeProduct = (id: number) => {
    const confirmModal = useModalConfirm({
        slot: 'Вы действительно хотите удалить товар из корзины?',
        onConfirm: () => {
            cartStore.remove(id);
        },
    });
    confirmModal.open();
};
</script>

<template>
    <div>
        <div :class="$style.products">
            <template
                v-for="item in products"
                :key="item.product.id"
            >
                <div>
                    <NuxtLink
                        :class="$style.img"
                        :to="`/product-${item.product.id}`"
                    >
                        <img
                            v-if="item.product.image_preview"
                            :src="item.product.image_preview"
                            :alt="item.product.name"
                            :title="item.product.name"
                            loading="lazy"
                        />
                        <span></span>
                    </NuxtLink>
                    <div :class="$style.info">
                        <div :class="$style.main_info">
                            <div :class="$style.name">
                                <NuxtLink :to="`/product-${item.product.id}`">
                                    {{ item.product.name }}
                                </NuxtLink>
                            </div>
                            <div :class="$style.price">{{ coolNumber(item.product.price) }} ₽</div>
                            <div :class="$style.stock_available">
                                <span :class="$style.title">Остаток на складе:</span>
                                <template v-if="item.product?.stock_available">
                                    <span :class="$style.available">{{ item.product?.stock_available }} шт.</span>
                                </template>
                                <template v-else>
                                    <span :class="$style.not_available">нет в наличии</span>
                                </template>
                            </div>
                        </div>
                        <div :class="$style.quantity">
                            <SharedUiInput
                                :model-value="item.quantity.toString()"
                                type="number"
                                :min="item.product.stock_available ? 1 : 0"
                                :max="item.product.stock_available"
                                :read-value-from-v-model-on-blur="true"
                                @update:model-value="setQuantity(item.product.id, Number($event))"
                            />
                        </div>
                        <div :class="$style.remove">
                            <button
                                title="Удалить"
                                @click="removeProduct(item.product.id)"
                            ></button>
                        </div>
                    </div>
                </div>
            </template>
        </div>
        <div :class="$style.products_total">
            Итого: <span>{{ coolNumber(total) }} ₽</span>
        </div>
    </div>
</template>

<style lang="less" module>
@import '@styles/includes';

.products {
    display: flex;
    flex-direction: column;
    gap: 40px;

    .width-size-sm-less({
        gap: 30px;
    });

    > div {
        display: flex;
        align-items: center;
        gap: 40px;

        .width-size-less(850px,{
            align-items: flex-start;
        });

        .width-size-sm-less({
            gap:20px;
        });

        > .img {
            max-width: 125px;
            width: 100%;
            flex-shrink: 0;
            display: block;
            position: relative;
            background-color: var(--ui-bg-color);
            border-radius: 10px;
            overflow: hidden;
            box-shadow: 0px 0px 10px rgba(0, 0, 0, 0.01);
            transition: box-shadow 0.3s ease;

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

            .hover({
                box-shadow: 0px 0px 10px rgba(0, 0, 0, 0.1);
            });

            .width-size-sm-less({
                max-width: 100px;
            });
        }

        > .info {
            width: 100%;
            gap: 40px;
            display: flex;
            align-items: center;

            .width-size-less(850px,{
                flex-wrap: wrap
            });

            .width-size-sm-less({
                gap:20px;
            });

            > .main_info {
                width: 100%;
                display: flex;
                flex-direction: column;
                gap: 5px;

                > .name {
                    font-family: 'Strong';
                    font-size: 20px;
                    color: var(--font-color-3);
                    word-break: break-word;

                    .width-size-sm-less({
                        font-size: 16px;
                    });
                }

                > .price {
                    font-weight: 600;
                    color: var(--font-color-3);
                }

                > .stock_available {
                    > .title {
                        font-size: 14px;
                        color: var(--font-color-4);
                        margin-right: 10px;

                        .width-size-sm-less({
                            font-size: 12px;
                        });
                    }

                    > .available {
                        font-weight: 600;
                        color: var(--color-4);

                        .width-size-less(850px,{
                            font-size: 14px;
                        });
                    }

                    > .not_available {
                        font-weight: 600;
                        color: var(--color-1);

                        .width-size-less(850px,{
                            font-size: 14px;
                        });
                    }
                }
            }

            > .quantity {
                width: 100px;
                flex-shrink: 0;

                input {
                    text-align: center;
                }
            }

            > .remove {
                flex-shrink: 0;

                > button {
                    display: block;
                    position: relative;
                    width: 24px;
                    height: 24px;
                    background-color: var(--color-2);
                    border-radius: 5px;

                    &::after {
                        display: block;
                        content: '';
                        position: absolute;
                        left: 0;
                        top: 0;
                        width: 100%;
                        height: 100%;
                        background-color: var(--stable-white-color);
                        mask: url('@/assets/icons/remove.svg') no-repeat center center;
                        mask-size: 16px auto;
                    }
                }
            }
        }
    }
}

.products_total {
    margin-top: 40px;
    font-size: 20px;
    text-align: right;

    > span {
        font-weight: 600;
    }

    .width-size-sm-less({
        font-size: 18px;
        margin-top: 30px;
    });
}
</style>
