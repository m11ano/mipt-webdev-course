<script setup lang="ts">
definePageMeta({});
</script>

<template>
    <div>
        <template v-if="$authData.isAuth"> <PageMain /></template>
        <template v-else> <PageLogin /></template>
    </div>
</template>
