<script setup lang="ts">
import { useModalConfirm } from '~/components/shared/modals/Confirm/useModalConfirm';
import { useCartStore } from '~/domain/shop';
import type { IProductListItem } from '~/domain/shop/model/types/product';
import { coolNumber } from '~/shared/helpers/functions';

const props = defineProps<{ productsItems: IProductListItem[] }>();

const cartStore = useCartStore();

interface ProductInCart {
    productID: number;
    quantity: number;
}

const productsMap = computed(() => {
    const result: Record<number, IProductListItem> = {};

    if (props.productsItems) {
        props.productsItems.forEach((item) => {
            result[item.id] = item;
        });
    }

    return result;
});

const getProduct = (id: number) => {
    if (typeof productsMap.value[id] === 'undefined') {
        return null;
    }
    return productsMap.value[id];
};

const products = computed<ProductInCart[]>(() => {
    const result: ProductInCart[] = [];
    cartStore.items.forEach((item) => {
        result.push({
            productID: item.id,
            quantity: item.quantity,
        });
    });

    return result;
});

const setQuantity = (id: number, quantity: number) => {
    if (isNaN(quantity) || quantity < 1) {
        quantity = 1;
    }

    if (!getProduct(id)) return;

    quantity = Math.min(quantity, getProduct(id)!.stock_available);
    cartStore.setQuantity(id, quantity);
};

const total = computed(() => {
    let total = 0;
    products.value.forEach((item) => {
        const product = getProduct(item.productID);
        if (product) {
            total += product.price * item.quantity;
        }
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
                :key="item.productID"
            >
                <div>
                    <template v-if="getProduct(item.productID)">
                        <NuxtLink
                            :class="$style.img"
                            :to="`/product-${item.productID}`"
                        >
                            <template v-if="getProduct(item.productID)">
                                <img
                                    v-if="getProduct(item.productID)?.image_preview"
                                    :src="getProduct(item.productID)?.image_preview"
                                    :alt="getProduct(item.productID)?.name"
                                    :title="getProduct(item.productID)?.name"
                                    loading="lazy"
                                />
                            </template>

                            <span></span>
                        </NuxtLink>
                    </template>
                    <template v-else>
                        <div :class="$style.img">
                            <div :class="['skeleton', 'bg', 'square']"><div></div></div>
                        </div>
                    </template>
                    <div :class="$style.info">
                        <div :class="$style.main_info">
                            <template v-if="getProduct(item.productID)">
                                <div :class="$style.name">
                                    <NuxtLink :to="`/product-${item.productID}`">
                                        {{ getProduct(item.productID)?.name }}
                                    </NuxtLink>
                                </div>
                                <div :class="$style.price">{{ coolNumber(getProduct(item.productID)?.price || 0) }} ₽</div>
                                <div :class="$style.stock_available">
                                    <span :class="$style.title">Остаток на складе:</span>
                                    <template v-if="getProduct(item.productID)?.is_published">
                                        <template v-if="getProduct(item.productID)?.stock_available">
                                            <span :class="$style.available">{{ getProduct(item.productID)?.stock_available }} шт.</span>
                                        </template>
                                        <template v-else>
                                            <span :class="$style.not_available">нет в наличии</span>
                                        </template>
                                    </template>
                                    <template v-else>
                                        <span :class="$style.not_available">снято с продажи</span>
                                    </template>
                                </div>
                            </template>
                            <template v-else>
                                <div
                                    :class="['skeleton', 'bg']"
                                    style="height: 40px"
                                ></div>
                            </template>
                        </div>
                        <div :class="$style.quantity">
                            <template v-if="getProduct(item.productID)">
                                <SharedUiInput
                                    :model-value="item.quantity.toString()"
                                    type="number"
                                    :min="getProduct(item.productID)?.stock_available && getProduct(item.productID)?.is_published ? 1 : 0"
                                    :max="getProduct(item.productID)?.is_published ? getProduct(item.productID)?.stock_available : 0"
                                    :read-value-from-v-model-on-blur="true"
                                    @update:model-value="setQuantity(getProduct(item.productID)?.id || 0, Number($event))"
                                />
                            </template>
                            <template v-else>
                                <div
                                    :class="['skeleton', 'bg']"
                                    style="height: 40px"
                                ></div>
                            </template>
                        </div>
                        <div :class="$style.remove">
                            <template v-if="getProduct(item.productID)">
                                <button
                                    title="Удалить"
                                    @click="removeProduct(getProduct(item.productID)?.id || 0)"
                                ></button>
                            </template>
                            <template v-else>
                                <div
                                    :class="['skeleton', 'bg']"
                                    style="height: 40px"
                                ></div>
                            </template>
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
