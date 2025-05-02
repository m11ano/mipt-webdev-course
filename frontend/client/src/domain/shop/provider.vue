<script setup lang="ts">
import { useCartStore } from './model/store/cart';

const cartStore = useCartStore();
cartStore.loadFromLocalStorage();

if (import.meta.client) {
    window.addEventListener('focus', cartStore.loadFromLocalStorage);
}
</script>

<template>
    <slot />
</template>
