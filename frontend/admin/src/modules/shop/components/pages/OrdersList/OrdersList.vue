<script setup lang="ts">
import { setMenu } from '~/plugins/app/model/actions/setMenu';
import { setModuleBreadcrums } from '~/modules/shop/domain/actions/setModuleBreadcrums';
import { ShopModule } from '~/modules/shop/const';
import type { TableColumn, DropdownMenuItem } from '@nuxt/ui';
import { OrderStatusParams, type IOrderListItem, type OrderStatus } from '~/modules/shop/domain/model/types/order';
import { fetchOrdersList } from '~/modules/shop/domain/api/fetchOrdersList';
import { coolNumber } from '~/shared/helpers/functions';

useSeoMeta({
    title: 'Список заказов',
});

setMenu(ShopModule.urlName, 'orders');

setModuleBreadcrums([
    {
        name: 'Заказы',
        to: '/orders',
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

const ordersList = ref<IOrderListItem[]>([]);

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
        const data = await fetchOrdersList(page.value, limit.value < 1 ? 1 : limit.value);
        if (data.items) {
            ordersList.value = data.items;
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

interface Order {
    id: number;
    order_sum: number;
    status: OrderStatus;
}

const prepariedItems = computed<Order[]>(() => {
    return ordersList.value.map((item) => {
        return {
            id: item.id,
            order_sum: item.order_sum,
            status: item.status,
        };
    });
});

const columns: TableColumn<Order>[] = [
    {
        accessorKey: 'id',
        header: 'ID',
    },
    {
        accessorKey: 'order_sum',
        header: 'Сумма заказа',
    },
    {
        accessorKey: 'status',
        header: 'Статус',
    },
    {
        id: 'action',
    },
];

function getDropdownActions(product: Order): DropdownMenuItem[][] {
    return [
        [
            {
                label: 'Просмотреть',
                icon: 'i-lucide-edit',
                to: `/${ShopModule.urlName}/orders/${product.id}`,
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
        <UTable
            v-model:column-pinning="columnPinning"
            :data="prepariedItems"
            :columns="columns"
            :loading="isLoading"
        >
            <template #id-cell="{ row }"> №{{ row.original.id }} </template>
            <template #order_sum-cell="{ row }"> {{ coolNumber(row.original.order_sum) }} руб. </template>
            <template #status-cell="{ row }">
                <UBadge
                    :variant="OrderStatusParams[row.original.status].variant"
                    :color="OrderStatusParams[row.original.status].color"
                    >{{ OrderStatusParams[row.original.status].title }}</UBadge
                >
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
