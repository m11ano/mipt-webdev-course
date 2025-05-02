<script setup lang="ts">
const props = defineProps<{ orderID: number; secretKey: string }>();

if (props.secretKey.length == 0) {
    throw createError({ statusCode: 404, statusMessage: 'Order not found' });
}

useHead({
    title: `Заказ №${props.orderID}`,
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
                        <div :class="$style.value">Ожидает обработки</div>
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
