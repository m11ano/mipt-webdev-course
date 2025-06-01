<script setup lang="ts">
import { setAppBreadcrumbs } from '~/plugins/app/model/actions/setAppBreadcrumbs';
import { setMenu } from '~/plugins/app/model/actions/setMenu';

useSeoMeta({
    title: 'Добро пожаловать',
});

setMenu('', '');

setAppBreadcrumbs([]);
</script>

<template>
    <div>
        <div :id="$style.title">Добро пожаловать в панель управления</div>
    </div>
</template>

<style lang="less" module>
@import '@styles/includes';

#title {
    text-align: center;
    font-family: 'Strong';
    font-size: 32px;

    .width-size-sm-less({
        font-size: 24px;
    });
}
</style>
