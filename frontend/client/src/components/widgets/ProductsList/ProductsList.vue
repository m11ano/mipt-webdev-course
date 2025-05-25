<script setup lang="ts">
import type { IProductListItem } from '~/domain/shop/model/types/product';
import { useFetchProductsList, useLoadProductsList } from './api/useFetchProductsList';

const props = withDefaults(
    defineProps<{
        skeletonCount?: number;
        limit?: number;
        gridClassName?: string;
    }>(),
    {
        skeletonCount: 4,
        limit: 20,
        gridClassName: 'grid_4_in_row_style_1',
    },
);

const page = ref(1);
const isLoading = ref(true);
const hasNextPage = ref(false);

const newsList = ref<IProductListItem[]>([]);

const { data, execute } = await useFetchProductsList(1, props.limit, true, import.meta.server);

watch(
    data,
    () => {
        if (data.value) {
            isLoading.value = false;
            hasNextPage.value = data.value.total > page.value * props.limit;
            newsList.value = data.value.items;
        }
    },
    { immediate: true },
);

const loadNextPage = async () => {
    if (isLoading.value || !hasNextPage.value) return;

    isLoading.value = true;
    try {
        const data = await useLoadProductsList(page.value + 1, props.limit);
        if (data.items) {
            page.value++;
            newsList.value = [...newsList.value, ...data.items];
            hasNextPage.value = data.total > page.value * props.limit;
        }
    } catch (e: unknown) {
        //
    } finally {
        isLoading.value = false;
    }
};

onMounted(() => {
    // setTimeout(() => {
    //     execute();
    // }, 2000);
    execute();
});
</script>

<template>
    <div>
        <div :class="gridClassName">
            <template
                v-for="item in newsList"
                :key="item.id"
            >
                <EntityProductCard :product="item" />
            </template>
            <template v-if="isLoading">
                <slot
                    v-for="i in skeletonCount"
                    :key="i"
                    name="skeleton"
                >
                    <div :class="['skeleton', 'bg', 'square']"><div></div></div>
                </slot>
            </template>
        </div>
        <template v-if="hasNextPage">
            <div :class="$style.nextPage">
                <button
                    class="button_1"
                    @click="loadNextPage"
                >
                    Показать еще
                </button>
            </div>
        </template>
    </div>
</template>

<style lang="less" module>
@import '@styles/includes';

.nextPage {
    margin-top: 40px;
    text-align: center;
    .width-size-sm-less({
        margin-top: 20px;
    });
}
</style>
