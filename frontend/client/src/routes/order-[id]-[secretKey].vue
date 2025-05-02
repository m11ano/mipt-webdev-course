<script setup lang="ts">
definePageMeta({
    validate: async (route) => {
        return typeof route.params.id === 'string' && /^\d+$/.test(route.params.id) && typeof route.params.secretKey === 'string';
    },
});

const route = useRoute();
</script>

<template>
    <LazyPageOrder
        :order-i-d="Number(route.params.id)"
        :secret-key="String(route.params.secretKey)"
    />
</template>

<style lang="less"></style>
