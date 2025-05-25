<script setup lang="ts">
import { addProductToCart } from '~/domain/shop';
import type { IProductItem } from '~/domain/shop/model/types/product';
import { coolNumber } from '~/shared/helpers/functions';

const props = defineProps<{
    id: number;
}>();

const productApiUrl = computed(() => `/products/${props.id}`);

const {
    data: product,
    error,
    execute,
} = await useAPIFetch<IProductItem>(productApiUrl, {
    lazy: true,
    immediate: import.meta.server,
});

watch(
    productApiUrl,
    () => {
        // setTimeout(() => {
        //     execute();
        // }, 2000);
        execute();
    },
    {
        immediate: true,
    },
);

const isLoading = computed(() => product.value === null);

watch(
    error,
    () => {
        if (error.value && error.value.statusCode) {
            showError({
                statusCode: error.value.statusCode,
                statusMessage: 'Товар не найден',
            });
        }
    },
    { immediate: true },
);

watch(
    product,
    () => {
        if (product.value) {
            useSeoMeta({
                title: product.value.name,
            });
        }
    },
    { immediate: true },
);
</script>

<template>
    <div>
        <div :id="$style.bread_crumbs">
            <NuxtLink to="/">Каталог товаров</NuxtLink>
            <template v-if="!isLoading">
                <span :class="$style.delim">/</span>
                <span :class="$style.title">{{ product?.name }}</span>
            </template>
        </div>
        <div :id="$style.product">
            <div>
                <div :class="$style.slider">
                    <template v-if="!isLoading">
                        <EntityProductSlider :slider="product?.slider || []" />
                    </template>
                    <template v-else>
                        <div :class="['skeleton', 'bg', 'square']">
                            <div></div>
                        </div>
                    </template>
                </div>
                <div :class="$style.content">
                    <template v-if="!isLoading">
                        <div :class="$style.name">{{ product?.name }}</div>
                        <div :class="$style.price">{{ coolNumber(product?.price || 0) }} ₽</div>
                        <div :class="$style.stock">
                            <div :class="$style.stock_available">
                                <span :class="$style.title">Остаток на складе:</span>
                                <template v-if="product?.stock_available">
                                    <span :class="$style.available">{{ product?.stock_available }} шт.</span>
                                </template>
                                <template v-else>
                                    <span :class="$style.not_available">нет в наличии</span>
                                </template>
                            </div>
                            <div
                                v-if="product?.stock_available"
                                :class="$style.add_to_cart"
                            >
                                <button
                                    class="button_1"
                                    @click="addProductToCart(product?.id || 0)"
                                >
                                    Добавить в корзину
                                </button>
                            </div>
                        </div>
                        <div
                            v-if="product?.full_description"
                            :class="$style.description"
                            v-html="product.full_description"
                        ></div>
                    </template>
                    <template v-else>
                        <div
                            :id="$style.skeleton_name"
                            :class="['skeleton', 'bg']"
                        ></div>
                    </template>
                </div>
            </div>
        </div>
    </div>
</template>

<style lang="less" module>
@import '@styles/includes';

#bread_crumbs {
    .std-wrapper();
    line-height: 1.5;
    font-size: 14px;

    > a {
        text-decoration: underline;
        color: var(--font-color-3);

        .hover({
            text-decoration: none;
        });
    }

    > .delim {
        margin: 0 5px;
    }

    > .title {
        color: var(--font-color-4);
    }

    .width-size-sm-less({
        font-size: 12px;
    });
}

#product {
    .std-wrapper(20px, 0, 20px, 0);

    > div {
        display: flex;
        gap: 40px;

        > .slider {
            width: calc(50% - 20px);
            min-width: 400px;
            max-width: 540px;
            flex-shrink: 0;
        }

        > .content {
            width: 100%;
        }

        .width-size-less(850px, {
            flex-wrap: wrap;
            justify-content: center;

            > .content {
                width: 100%;
            }
        });

        .width-size-sm-less({
            gap:20px;

            > .slider {
                width: 100%;
                min-width: 0;
            }
        });
    }
}

.name {
    font-size: 32px;
    font-family: 'Strong';
    color: var(--font-color-3);

    .width-size-sm-less({
        font-size: 20px;
    });
}

.price {
    font-size: 24px;
    font-weight: 600;
    margin-top: 20px;
    color: var(--font-color-3);

    .width-size-sm-less({
        font-size: 20px;
    });
}

.stock {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    justify-content: space-between;
    gap: 15px;
    margin-top: 30px;
    padding: 20px;
    border-radius: 10px;
    background-color: var(--bg-color-2);

    .width-size-sm-less({
        margin-top: 20px;
        padding:15px;
    });

    > .stock_available {
        > .title {
            .width-size-sm-less({
                font-size: 12px;
            });
        }

        > .available {
            margin-left: 10px;
            font-weight: 600;
            color: var(--color-4);
        }

        > .not_available {
            margin-left: 10px;
            font-weight: 600;
            color: var(--color-1);
        }
    }

    > .add_to_cart {
        flex-shrink: 0;
    }
}

.description {
    margin-top: 40px;
    line-height: 1.5;

    .width-size-sm-less({
        margin-top: 30px;
        font-size: 12px;
    });
}

#skeleton_name {
    height: 60px;
}
</style>
