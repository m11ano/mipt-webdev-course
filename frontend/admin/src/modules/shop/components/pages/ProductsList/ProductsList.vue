<script setup lang="ts">
import { setMenu } from '~/plugins/app/model/actions/setMenu';
import { setModuleBreadcrums } from '~/modules/shop/domain/actions/setModuleBreadcrums';
import { ShopModule } from '~/modules/shop/const';
import type { TableColumn, DropdownMenuItem } from '@nuxt/ui';
import type { IProductListItem } from '~/modules/shop/domain/model/types/product';
import { fetchProductsList } from '../../../domain/api/fetchProductsList';
import { removeProduct } from '~/modules/shop/domain/actions/removeProduct';

useSeoMeta({
    title: 'Список товаров',
});

setMenu(ShopModule.urlName, 'products');

setModuleBreadcrums([
    {
        name: 'Товары',
        to: '/products',
    },
]);

const route = useRoute();
const router = useRouter();

let routePage = Number(route.query.page);
if (!routePage || isNaN(routePage)) {
    routePage = 1;
}

const page = ref(routePage);

const isLoading = ref(true);

const productsList = ref<IProductListItem[]>([]);

const defaultLimit = 20;

let routeLimit = Number(route.query.limit);
if (!routeLimit || isNaN(routeLimit)) {
    routeLimit = defaultLimit;
}

const limit = ref(routeLimit);

const total = ref(0);

const fetchData = async () => {
    isLoading.value = true;

    try {
        const data = await fetchProductsList(page.value, limit.value < 1 ? 1 : limit.value);
        if (data.items) {
            productsList.value = data.items;
            total.value = data.total;

            router.push({
                query: {
                    ...route.query,
                    page: page.value > 1 ? page.value : undefined,
                    limit: limit.value !== defaultLimit ? limit.value : undefined,
                },
            });
        }
    } catch (e: unknown) {
        //
    } finally {
        isLoading.value = false;
    }
};

const onPageUpdate = (p: number) => {
    if (isLoading.value) return;
    page.value = p;
    fetchData();
};

watch(limit, () => {
    page.value = 1;
    fetchData();
});

interface Product {
    id: number;
    image_preview: string;
    is_published: boolean;
    info: {
        name: string;
        price: number;
        stock_available: number;
    };
}

const prepariedItems = computed(() => {
    return productsList.value.map((item) => {
        return {
            id: item.id,
            image_preview: item.image_preview,
            is_published: item.is_published,
            info: {
                name: item.name,
                price: item.price,
                stock_available: item.stock_available,
            },
        };
    });
});

const columns: TableColumn<Product>[] = [
    {
        accessorKey: 'id',
        header: 'ID',
    },
    {
        accessorKey: 'image_preview',
        header: 'Фото',
    },
    {
        accessorKey: 'is_published',
        header: 'Статус',
    },
    {
        accessorKey: 'info',
        header: 'Товар',
        cell(props) {
            console.log(props);
            return {};
        },
    },
    {
        id: 'action',
    },
];

function getDropdownActions(product: Product): DropdownMenuItem[][] {
    return [
        [
            {
                label: 'Редактировать',
                icon: 'i-lucide-edit',
                to: `/${ShopModule.urlName}/products/${product.id}`,
            },
            {
                label: 'Удалить',
                icon: 'i-lucide-trash',
                color: 'error',
                onSelect: async () => {
                    const result = await removeProduct(product.id);
                    if (result) {
                        productsList.value = productsList.value.filter((p) => p.id !== product.id);
                        total.value -= 1;
                    }
                },
            },
        ],
    ];
}

const columnPinning = ref({ left: [], right: ['action'] });

onMounted(() => {
    fetchData();
});
</script>

<template>
    <div>
        <div class="flex justify-end">
            <UButton :to="`/${ShopModule.urlName}/products/new`">Создать товар</UButton>
        </div>
        <UTable
            v-model:column-pinning="columnPinning"
            :data="prepariedItems"
            :columns="columns"
            :loading="isLoading"
        >
            <template #image_preview-cell="{ row }">
                <img
                    :src="row.original.image_preview"
                    alt=""
                    style="width: 100px; min-width: 100px"
                />
            </template>
            <template #is_published-cell="{ row }">
                <template v-if="row.original.is_published">
                    <UBadge
                        variant="subtle"
                        color="success"
                        >Опубликовано</UBadge
                    >
                </template>
                <template v-else>
                    <UBadge variant="subtle">Не опубликовано</UBadge>
                </template>
            </template>
            <template #info-cell="{ row }">
                <div style="max-width: 300px; min-width: 200px; overflow: hidden; white-space: normal">
                    <div>
                        <b>{{ row.original.info.name }}</b>
                    </div>
                    <div>Цена: {{ row.original.info.price }} руб.</div>
                    <div>Свободный остаток: {{ row.original.info.stock_available }} шт.</div>
                </div>
            </template>
            <template #action-cell="{ row }">
                <UDropdownMenu :items="getDropdownActions(row.original)">
                    <UButton
                        icon="i-lucide-ellipsis-vertical"
                        color="neutral"
                        variant="ghost"
                        aria-label="Actions"
                    />
                </UDropdownMenu>
            </template>
        </UTable>
        <SharedPaginator
            v-model="limit"
            :disabled="isLoading"
        >
            <UPagination
                :page="page"
                :items-per-page="limit"
                :total="total"
                @update:page="onPageUpdate"
            />
        </SharedPaginator>
    </div>
</template>

<style lang="less" module>
@import '@styles/includes';
</style>
