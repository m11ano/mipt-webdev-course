<script setup lang="ts">
defineProps<{ orderID: number; link: string }>();
</script>

<template>
    <div :class="$style.wrapper">
        Заказу присвоен №{{ orderID }}.<br />
        <br />
        Наши специалисты свяжутся с Вами в ближайшее время.<br />
        <br />
        Следить за заказом можно <NuxtLink :to="link">по ссылке</NuxtLink>.
    </div>
</template>

<style lang="less" module>
@import '@styles/includes';

.wrapper {
    a {
        color: var(--color-1);
        text-decoration: underline;
        transition: color 0.25s ease;

        .hover({
        color:var(--color-1-hover);
    });
    }
}
</style>
