<script setup lang="ts">
import type { UseModalReturnType } from 'vue-final-modal';

const props = defineProps<{
    orderID: MaybeRef<number>;
    link: MaybeRef<string>;
    isLoaded: MaybeRef<boolean>;
    modalObj: () => UseModalReturnType<any>;
}>();
</script>

<template>
    <SharedModalsDefaultModal :modal-obj="props.modalObj()">
        <template v-if="!isLoaded"><SharedModalsDefaultLoader /></template>
        <template v-else>
            <SharedModalsDefaultWrapper
                :modal-obj="props.modalObj()"
                :class-name="$style.box_wrapper"
            >
                <template #title> Заказ создан </template>
                <template #content>
                    <div :class="$style.wrapper">
                        Заказу присвоен №{{ orderID }}.<br />
                        <br />
                        Наши специалисты свяжутся с Вами в ближайшее время.<br />
                        <br />
                        Следить за заказом можно
                        <NuxtLink :to="toValue(link)">по ссылке</NuxtLink>.
                    </div>
                </template>
                <template #buttons>
                    <NuxtLink
                        :to="toValue(link)"
                        class="button_1"
                        @click="modalObj().close()"
                    >
                        Перейти к заказу
                    </NuxtLink>
                </template>
            </SharedModalsDefaultWrapper>
        </template>
    </SharedModalsDefaultModal>
</template>

<style lang="less" module>
@import '@styles/includes';

.box_wrapper {
    width: 100%;
    max-width: 750px;
}

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
