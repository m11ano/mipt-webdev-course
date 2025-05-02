<script setup lang="ts">
import { ModalsContainer } from 'vue-final-modal';
import { AppProvider } from './domain/app';
import { ShopProvider } from './domain/shop';
</script>

<template>
    <AppProvider>
        <ShopProvider>
            <NuxtLayout>
                <NuxtPage />
            </NuxtLayout>
            <ModalsContainer />
        </ShopProvider>
    </AppProvider>
</template>
