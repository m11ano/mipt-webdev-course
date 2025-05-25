<script setup lang="ts">
import { OrderStatusText, type IOrder } from '~/domain/shop/model/types/order';
import { coolNumber } from '~/shared/helpers/functions';

const props = defineProps<{ orderID: number; secretKey: string }>();

const orderApiUrl = computed(() => `/orders/${props.orderID}/${props.secretKey}`);

const { data: order, error } = await useAPIFetch<IOrder>(orderApiUrl, {
    lazy: true,
});

watch(
    error,
    () => {
        if (error.value) {
            throw createError({ statusCode: 404, statusMessage: 'Заказ не найден' });
        }
    },
    {
        immediate: true,
    },
);

watch(order, () => {
    if (order.value) {
        useHead({
            title: `Заказ №${order.value.id}`,
        });
    }
});
</script>

<template>
    <div>
        <div :id="$style.block_order">
            <div>
                <div :class="$style.title">Проверка заказа</div>
                <div :class="$style.data">
                    <div>
                        <div :class="$style.title">Заказ:</div>
                        <div :class="$style.value">№{{ orderID }}</div>
                    </div>
                    <div>
                        <div :class="$style.title">Статус:</div>
                        <div :class="$style.value">
                            <template v-if="order?.status">{{ OrderStatusText[order.status] }}</template
                            ><template v-else>Загрузка</template>
                        </div>
                    </div>
                    <div>
                        <div :class="$style.title">Сумма:</div>
                        <div :class="$style.value">
                            <template v-if="order?.status">{{ coolNumber(order.order_sum) }} руб.</template
                            ><template v-else>Загрузка</template>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<style lang="less" module>
@import '@styles/includes';

#block_order {
    .std-wrapper();

    > div {
        max-width: 540px;
        margin: 0 auto;

        > .title {
            font-size: 32px;
            text-align: center;
            font-family: 'Strong';
            color: var(--font-color-3);

            .width-size-sm-less({
                font-size:24px;
            });
        }

        > .data {
            margin-top: 35px;
            display: flex;
            flex-direction: column;
            gap: 20px;

            > div {
                display: flex;
                align-items: center;
                gap: 20px;

                .width-size-sm-less({
                margin-top: 20px;
            });

                > .title {
                    width: calc(50% - 10px);
                    font-size: 20px;
                }

                > .value {
                    width: calc(50% - 10px);
                    text-align: right;
                    font-size: 20px;
                }
            }
        }
    }
}
</style>
