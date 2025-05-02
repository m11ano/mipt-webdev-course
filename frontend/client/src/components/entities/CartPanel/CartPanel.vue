<script setup lang="ts">
import { useCartStore } from '~/domain/shop';

const cartStore = useCartStore();
</script>

<template>
    <component
        :is="cartStore.count > 0 ? 'router-link' : 'span'"
        :to="cartStore.count > 0 ? '/cart' : undefined"
        :class="[$style.cart, cartStore.count > 0 ? $style.active : null]"
    >
        <span :class="$style.counter">{{ cartStore.count }}</span>
        <span :class="$style.icon"></span>
        <template v-if="cartStore.count > 0">
            <span :class="$style.button">Оформить заказ</span>
        </template>
        <template v-else>
            <span :class="$style.empy_cart">Корзина пуста</span>
        </template>
    </component>
</template>

<style lang="less" module>
@import '@styles/includes';

.cart {
    display: flex;
    align-items: center;
    gap: 8px;
}

.counter {
    min-width: 22px;
    padding: 0 5px;
    height: 22px;
    display: flex;
    align-items: center;
    justify-content: center;
    background-color: var(--color-2);
    color: var(--stable-white-color);
    border-radius: 5px;
    line-height: 1;
    transition: background-color 0.25s ease;
}

.icon {
    display: block;
    width: 27px;
    height: 27px;
    mask: url('@/assets/icons/basket.svg') no-repeat center center;
    mask-size: 100% auto;
    background-color: var(--stable-black-color);

    .width-size-sm-less({
        width: 24px;
    });
}

.button {
    margin-left: 10px;
    display: flex;
    font-size: 16px;
    letter-spacing: 1px;
    text-decoration: none !important;
    background-color: var(--color-1);
    color: var(--stable-white-color);
    min-height: 40px;
    padding: 10px 15px;
    border-radius: 5px;
    text-align: center;
    align-items: center;
    justify-content: center;
    line-height: 1.3;
    transition:
        background-color 0.3s,
        border 0.3s;
    -webkit-background-clip: padding-box;
    background-clip: padding-box;

    .width-size-sm-less({
        display: none;
    });
}

.dark-theme({
    .icon {
        background-color: var(--stable-white-color);
    }
});

.empy_cart {
    font-size: 14px;
    color: var(--font-color-2);

    .width-size-sm-less({
        display: none;
    });
}

.cart.active {
    > .counter {
        background-color: var(--color-1);
    }

    .hover({
        > .counter {
            background-color: var(--color-1-hover);
        }

        > .button {
            background-color: var(--color-1-hover);
        }
    });
}
</style>
