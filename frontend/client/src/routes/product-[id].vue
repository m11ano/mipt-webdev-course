<script setup lang="ts">
definePageMeta({
    validate: async (route) => {
        return typeof route.params.id === 'string' && /^\d+$/.test(route.params.id);
    },
});

const route = useRoute();
</script>

<template>
    <LazyPageProduct :id="Number(route.params.id)" />
</template>

<style lang="less"></style>
