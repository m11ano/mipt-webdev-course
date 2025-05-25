<script setup lang="ts">
import { useCartStore } from '~/domain/shop';
import type { IProductListItem } from '~/domain/shop/model/types/product';

useHead({
    title: 'Корзина',
});

const cartStore = useCartStore();

const productsIDs = computed(() => cartStore.items.map((i) => i.id));

const productApiUrl = computed(() => `/products?ids=${productsIDs.value.join(',')}`);

const { data: productsItems, execute } = await useAPIFetch<{ items: IProductListItem[] }>(productApiUrl, {
    lazy: true,
    immediate: false,
});

// setTimeout(() => {
//     execute();
// }, 3000);

execute();
</script>

<template>
    <div>
        <div :id="$style.block_basket">
            <div :class="$style.title">Корзина</div>
            <template v-if="cartStore.items.length == 0">
                <div :class="$style.empty_title">Корзина пуста</div>
            </template>
            <template v-else>
                <div :class="$style.basket">
                    <WidgetCartBasket :products-items="productsItems?.items || []" />
                </div>
            </template>
        </div>
        <div
            :id="$style.block_order"
            :style="{ display: cartStore.items.length == 0 ? 'none' : 'block' }"
        >
            <div>
                <div :class="$style.title">Оформление заказа</div>
                <div :class="$style.order">
                    <FeatureCartOrder />
                </div>
            </div>
        </div>
    </div>
</template>

<style lang="less" module>
@import '@styles/includes';

#block_basket {
    .std-wrapper();

    > .title {
        font-size: 32px;
        text-align: center;
        font-family: 'Strong';
        color: var(--font-color-3);

        .width-size-sm-less({
            font-size:24px;
        });
    }

    > .empty_title {
        margin-top: 35px;
    }

    > .basket {
        margin-top: 35px;

        .width-size-sm-less({
            margin-top: 20px;
        });
    }
}

#block_order {
    .std-wrapper(40px, 0, 30px, 0);

    > div {
        max-width: 540px;
        margin: 0 auto;

        > .title {
            font-size: 32px;
            text-align: center;
            font-family: 'Strong';
            color: var(--font-color-3);

            .width-size-sm-less({
                font-size:24px;
            });
        }

        > .order {
            margin-top: 35px;

            .width-size-sm-less({
                margin-top: 20px;
            });
        }
    }
}
</style>
