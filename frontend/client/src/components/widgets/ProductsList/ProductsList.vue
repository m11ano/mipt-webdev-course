<script setup lang="ts">
import { mockProductsList } from '~/mocks/products';
</script>

<template>
    <div>
        <div class="grid_4_in_row_style_1">
            <template
                v-for="item in mockProductsList"
                :key="item.id"
            >
                <EntityProductCard :product="item" />
            </template>
        </div>
    </div>
</template>

<style lang="less" module>
@import '@styles/includes';
</style>
