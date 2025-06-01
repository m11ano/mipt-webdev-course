<script setup lang="ts">
const route = useRoute();
const rawId = route.params.id as string;
const id = Number(rawId);
</script>
<template><ShopPageProductsEdit :id="id" /></template>
