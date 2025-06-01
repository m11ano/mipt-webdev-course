<script setup lang="ts">
useSeoMeta({
    title: 'Вход в панель управления',
});
</script>

<template>
    <div :id="$style.form"><FeatureLogin /></div>
</template>

<style lang="less" module>
@import '@styles/includes';

#form {
    display: flex;
    justify-content: center;
    .std-wrapper(@std-x-padding, @std-x-padding, @std-x-padding-sm, @std-x-padding-sm);
}
</style>
