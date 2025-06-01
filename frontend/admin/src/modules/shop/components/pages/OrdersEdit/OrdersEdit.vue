<script setup lang="ts">
import { setMenu } from '~/plugins/app/model/actions/setMenu';
import { setModuleBreadcrums } from '~/modules/shop/domain/actions/setModuleBreadcrums';
import { ShopModule } from '~/modules/shop/const';
import { StandartErrorList } from '~/shared/errors/errors';
import ProductStock from '../../shared/modals/ProductStock.vue';
import { OrderStatus, OrderStatusParams, type IOrderItem } from '~/modules/shop/domain/model/types/order';
import { fetchOrder } from '~/modules/shop/domain/api/fetchOrder';
import Confirm from '~/core/components/shared/Confirm/modals/Confirm.vue';
import { setOrderStatus } from '~/modules/shop/domain/api/setOrderStatus';
import type { IProductListItem } from '~/modules/shop/domain/model/types/product';
import { fetchProductsListByIDs } from '~/modules/shop/domain/api/fetchProductsListByIDs';
import { coolNumber } from '~/shared/helpers/functions';
import { updateProduct } from '~/modules/shop/domain/api/updateProduct';
import { updateOrder } from '~/modules/shop/domain/api/updateOrder';
import AddProduct from '../../shared/modals/AddProduct.vue';

const props = defineProps<{
    id: number;
}>();

useSeoMeta({
    title: 'Просмотр заказа',
});

setMenu(ShopModule.urlName, 'orders');

setModuleBreadcrums([
    {
        name: 'Заказы',
        to: '/orders',
    },
    {
        name: 'Просмотр заказа',
    },
]);

const orderModel = ref<IOrderItem | null>(null);

const isOrderLoading = ref(false);
const isProductsLoading = ref(false);

const isLoading = computed(() => isOrderLoading.value || isProductsLoading.value);

const errors = ref<string[]>([]);

const toast = useToast();

watch(
    () => props.id,
    async () => {
        isOrderLoading.value = true;
        try {
            const data = await fetchOrder(props.id);
            orderModel.value = data;
        } catch (e) {
            if (e instanceof StandartErrorList) {
                if (e.code === 404) {
                    showError({
                        statusCode: e.code,
                        statusMessage: 'Заказ не найден',
                    });
                }
                errors.value = e.details;
            }
        } finally {
            isOrderLoading.value = false;
        }
    },
    {
        immediate: true,
    },
);

const isEditable = computed(() => orderModel.value?.status !== OrderStatus.Canceled && orderModel.value?.status !== OrderStatus.Finished);

const orderProductsIDs = computed(() => {
    return orderModel.value?.products.map((item) => item.id) || [];
});

const orderProducts = ref<Record<number, IProductListItem>>({});

watch(
    orderProductsIDs,
    async () => {
        isProductsLoading.value = true;
        try {
            const data = await fetchProductsListByIDs(orderProductsIDs.value);
            orderProducts.value = {};
            data.items.forEach((item) => {
                orderProducts.value[item.id] = item;
            });
        } catch (e) {
            if (e instanceof StandartErrorList) {
                console.log(e);
            }
        } finally {
            isProductsLoading.value = false;
        }
    },
    {
        immediate: true,
    },
);

const orderTotalSum = computed(() => {
    return orderModel.value?.products.reduce((acc, item) => acc + item.price * item.quantity, 0) || 0;
});

const save = async () => {
    if (isLoading.value || !orderModel.value) return;

    errors.value = [];

    if (orderModel.value.details.client_name.length < 1) {
        errors.value.push('Имя клиента не указано');
    }

    if (orderModel.value.details.client_surname.length < 1) {
        errors.value.push('Фамилия клиента не указана');
    }

    if (orderModel.value.details.client_email.length < 1) {
        errors.value.push('E-mail клиента не указан');
    }

    if (orderModel.value.details.client_phone.length < 1) {
        errors.value.push('Телефон клиента не указан');
    }

    if (orderModel.value.details.delivery_address.length < 1) {
        errors.value.push('Адрес доставки не указан');
    }

    if (orderModel.value.products.length < 1) {
        errors.value.push('Список товаров не указан');
    }

    if (errors.value.length) return;

    isOrderLoading.value = true;
    try {
        await updateOrder(orderModel.value);

        toast.add({
            title: 'Успех',
            description: 'Заказ сохранен',
            color: 'success',
            icon: 'i-lucide-check-circle',
        });
    } catch (e) {
        if (e instanceof StandartErrorList) {
            errors.value = e.details;
        }
    } finally {
        isOrderLoading.value = false;
    }
};

const overlay = useOverlay();
const addProduct = async () => {
    if (!orderModel.value) return;

    const modal = overlay.create(AddProduct, {
        props: {
            orderID: orderModel.value.id,
        },
        destroyOnClose: true,
    });

    const instance = modal.open();

    const newProduct = await instance.result;
    if (newProduct) {
        const check = orderModel.value.products.findIndex((item) => item.id === newProduct.id);
        if (check < 0) {
            orderModel.value.products = [
                ...orderModel.value.products,
                {
                    id: newProduct.id,
                    quantity: 1,
                    price: newProduct.price,
                },
            ];
        } else {
            orderModel.value.products[check].quantity++;
            orderModel.value.products = [...orderModel.value.products];
        }
    }
};

const setStatus = async (status: OrderStatus) => {
    if (!orderModel.value) return false;

    let text = '';

    switch (status) {
        case OrderStatus.Finished:
            text = 'Вы действительно хотите завершить заказ?';
            break;
        case OrderStatus.Canceled:
            text = 'Вы действительно хотите отменить заказ?';
            break;
        case OrderStatus.InWork:
            text = 'Вы действительно хотите взять заказ в работу?';
            break;
    }

    const modal = useOverlay().create(Confirm, {
        props: {
            text,
        },
        destroyOnClose: true,
    });

    const instance = modal.open();

    const shouldDo = await instance.result;
    if (shouldDo) {
        try {
            await setOrderStatus(orderModel.value.id, status);
            orderModel.value.status = status;

            useToast().add({
                title: 'Успех',
                description: 'Статус изменен',
                color: 'success',
                icon: 'i-lucide-check-circle',
            });
        } catch (e) {
            useToast().add({
                title: 'Ошибка',
                description: 'Ошибка при обновлении статуса',
                color: 'error',
                icon: 'i-lucide-ban',
            });
        }
    }
};

const removeProduct = async (productID: number) => {
    if (!orderModel.value) return false;

    const modal = useOverlay().create(Confirm, {
        props: {
            text: 'Вы действительно хотите удалить товар из заказа?',
        },
        destroyOnClose: true,
    });

    const instance = modal.open();

    const shouldDo = await instance.result;
    if (shouldDo) {
        orderModel.value.products = orderModel.value.products.filter((item) => item.id !== productID);
    }
};
</script>

<template>
    <div>
        <div
            v-if="isEditable"
            class="flex justify-end gap-4 mb-4"
        >
            <UButton
                color="graylight"
                variant="subtle"
                @click="setStatus(OrderStatus.Canceled)"
            >
                Отменить
            </UButton>
            <UButton
                v-if="orderModel?.status === OrderStatus.Created"
                color="info"
                variant="subtle"
                @click="setStatus(OrderStatus.InWork)"
            >
                Взять в работу
            </UButton>
            <UButton
                v-if="orderModel?.status === OrderStatus.InWork"
                color="success"
                variant="subtle"
                @click="setStatus(OrderStatus.Finished)"
            >
                Завершить
            </UButton>
        </div>
        <div>
            <div class="form-table">
                <div>
                    <div class="title">ID:</div>
                    <div class="value">№{{ orderModel?.id }}</div>
                </div>
                <div>
                    <div class="title">URL для клиента:</div>
                    <div class="value">
                        <a
                            :href="`/order-${orderModel?.id}-${orderModel?.secret_key}`"
                            target="_blank"
                            style="text-decoration: underline"
                            >/order-{{ orderModel?.id }}-{{ orderModel?.secret_key }}</a
                        >
                    </div>
                </div>
                <div>
                    <div class="title">Статус:</div>
                    <div class="value">
                        <div class="flex items-center gap-4">
                            <div>
                                <UBadge
                                    v-if="orderModel"
                                    :variant="OrderStatusParams[orderModel.status].variant"
                                    :color="OrderStatusParams[orderModel.status].color"
                                    >{{ OrderStatusParams[orderModel.status].title }}</UBadge
                                >
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <div class="form_title mt-10">
            <div class="title">Карточка заказа</div>
            <div
                v-if="isEditable"
                class="buttons"
            >
                <UButton
                    :loading="isLoading"
                    :disabled="isLoading"
                    @click="save"
                    >Сохранить изменения</UButton
                >
            </div>
        </div>
        <div
            v-if="errors.length"
            class="mt-4"
        >
            <UAlert
                title="Возникли ошибки!"
                icon="i-lucide-ban"
            >
                <template #description>
                    <template
                        v-for="error in errors"
                        :key="error"
                    >
                        <div>– {{ error }}</div>
                    </template>
                </template>
            </UAlert>
        </div>
        <div class="mt-4">
            <div class="form-table">
                <div>
                    <div class="title">Имя клиента:</div>
                    <div class="value">
                        <template v-if="isEditable">
                            <UInput
                                v-if="orderModel"
                                v-model="orderModel.details.client_name"
                                size="xl"
                                class="w-full"
                                :disabled="isLoading"
                            />
                        </template>
                        <template v-else>
                            {{ orderModel?.details.client_name }}
                        </template>
                    </div>
                </div>
                <div>
                    <div class="title">Фамилия клиента:</div>
                    <div class="value">
                        <template v-if="isEditable">
                            <UInput
                                v-if="orderModel"
                                v-model="orderModel.details.client_surname"
                                size="xl"
                                class="w-full"
                                :disabled="isLoading"
                            />
                        </template>
                        <template v-else>
                            {{ orderModel?.details.client_surname }}
                        </template>
                    </div>
                </div>
                <div>
                    <div class="title">E-mail клиента:</div>
                    <div class="value">
                        <template v-if="isEditable">
                            <UInput
                                v-if="orderModel"
                                v-model="orderModel.details.client_email"
                                size="xl"
                                class="w-full"
                                :disabled="isLoading"
                            />
                        </template>
                        <template v-else>
                            {{ orderModel?.details.client_email }}
                        </template>
                    </div>
                </div>
                <div>
                    <div class="title">Телефон клиента:</div>
                    <div class="value">
                        <template v-if="isEditable">
                            <UInput
                                v-if="orderModel"
                                v-model="orderModel.details.client_phone"
                                size="xl"
                                class="w-full"
                                :disabled="isLoading"
                            />
                        </template>
                        <template v-else>
                            {{ orderModel?.details.client_phone }}
                        </template>
                    </div>
                </div>
                <div>
                    <div class="title">Адрес доставки:</div>
                    <div class="value">
                        <template v-if="isEditable">
                            <UInput
                                v-if="orderModel"
                                v-model="orderModel.details.delivery_address"
                                size="xl"
                                class="w-full"
                                :disabled="isLoading"
                            />
                        </template>
                        <template v-else>
                            {{ orderModel?.details.delivery_address }}
                        </template>
                    </div>
                </div>
            </div>
        </div>
        <div class="form_title sub mt-4">
            <div class="title">Товары в заказе</div>
            <div
                v-if="isEditable"
                class="buttons"
            >
                <UButton
                    variant="subtle"
                    :disabled="isLoading"
                    @click="addProduct"
                    >Добавить товар</UButton
                >
            </div>
        </div>
        <div
            v-if="orderModel"
            :class="['mt-4', $style.productsBlock]"
        >
            <template
                v-for="(product, index) in orderModel.products"
                :key="product.id"
            >
                <div v-if="orderProducts[product.id]">
                    <div :class="$style.image">
                        <a
                            :href="`/product-${product.id}`"
                            target="_blank"
                        >
                            <template v-if="orderProducts[product.id].image_preview">
                                <img
                                    :src="orderProducts[product.id].image_preview"
                                    alt=""
                                />
                            </template>
                        </a>
                    </div>
                    <div :class="$style.info">
                        <div :class="$style.name">
                            <a
                                :href="`/product-${product.id}`"
                                target="_blank"
                            >
                                {{ orderProducts[product.id].name }}
                            </a>
                        </div>
                        <div :class="$style.price"><span>Цена в заказе:</span> {{ coolNumber(product.price) }} руб.</div>
                    </div>
                    <div :class="$style.quantity">
                        <div :class="$style.title">Количество (шт.):</div>
                        <div :class="$style.value">
                            <template v-if="isEditable">
                                <UInputNumber
                                    v-model="orderModel.products[index].quantity"
                                    size="sm"
                                    :disabled="isLoading"
                                    orientation="vertical"
                                    :min="1"
                                />
                            </template>
                            <template v-else>
                                {{ orderModel.products[index].quantity }}
                            </template>
                        </div>
                    </div>
                    <div
                        v-if="isEditable"
                        :class="$style.actions"
                    >
                        <UButton
                            icon="i-lucide-trash"
                            size="xs"
                            color="graylight"
                            :disabled="isLoading"
                            @click="removeProduct(product.id)"
                        />
                    </div>
                </div>
            </template>
        </div>
        <div class="form_title sub mt-4">
            <div class="title">Сумма заказа</div>
            <div class="buttons">
                <div style="font-size: 24px">{{ coolNumber(orderTotalSum) }} руб.</div>
            </div>
        </div>
    </div>
</template>

<style lang="less" module>
@import '@styles/includes';

.productsBlock {
    display: flex;
    flex-direction: column;
    gap: 15px;

    > div {
        background-color: var(--ui-color-neutral-100);
        padding: 15px;
        display: flex;
        gap: 15px;
        align-items: center;

        .width-size-less(700px, {
            flex-wrap: wrap;
            justify-content: space-between
        });

        > .image {
            width: 50px;
            flex-shrink: 0;

            > a {
                > img {
                    width: 50px;
                }
            }
        }

        > .info {
            width: 50%;

            .width-size-less(700px, {
                width: calc(100% - 65px);
            });

            > .name {
                > a {
                    color: var(--ui-color-primary-600);
                    text-decoration: underline;
                }
            }

            > .price {
                margin-top: 5px;

                > span {
                    color: var(--ui-text-muted);
                }
            }
        }

        > .quantity {
            width: 50%;
            min-width: 120px;

            .width-size-less(700px, {
                width: 120px;
            });

            > .title {
                margin-bottom: 5px;
                color: var(--ui-text-muted);
            }
        }

        > .actions {
            flex-shrink: 0;
            width: 32px;
        }
    }
}
</style>
