<script setup lang="ts">
import { useModalErrorsList } from '~/components/shared/modals/ErrorsList/useModalErrorsList';
import { StandartErrorList } from '~/shared/errors/errors';
import { useCartStore } from '~/domain/shop';
import type { IOrderDetailsFormData, IOrderFormData } from './model/types/types';
import { fetchCreateOrder } from './api/fetchCreateOrder';
import { useOrderResult } from './modals/OrderResult/hooks';

const cartStore = useCartStore();

const orderFormData = ref<IOrderDetailsFormData>({
    client_name: '',
    client_surname: '',
    client_email: '',
    client_phone: '',
    delivery_address: '',
});

const errors = ref<string[]>([]);
const isPhoneComplete = ref(false);

const orderID = ref(0);
const orderLink = ref('');
const isLoaded = ref(false);

const resultModal = useOrderResult({
    orderID,
    link: orderLink,
    isLoaded: isLoaded,
});

const errorModal = useModalErrorsList({
    errors,
});

const isSending = ref(false);

const sendForm = async (e: Event) => {
    e.preventDefault();
    if (isSending.value) return;

    errors.value = [];
    if (orderFormData.value.client_name.replaceAll(' ', '').length === 0) {
        errors.value.push('Имя не указано');
    }
    if (orderFormData.value.client_surname.replaceAll(' ', '').length === 0) {
        errors.value.push('Фамилия не указана');
    }
    if (orderFormData.value.client_email.replaceAll(' ', '').length === 0) {
        errors.value.push('E-mail не указан');
    }
    if (!isPhoneComplete.value) {
        errors.value.push('Телефон указан некорректно');
    }
    if (orderFormData.value.delivery_address.replaceAll(' ', '').length === 0) {
        errors.value.push('Адрес не указан');
    }

    if (errors.value.length) {
        errorModal.open();
    } else {
        const sendData: IOrderFormData = {
            products: cartStore.items.map((item) => ({
                id: item.id,
                quantity: item.quantity,
            })),
            details: {
                client_email: orderFormData.value.client_email,
                client_name: orderFormData.value.client_name,
                client_phone: orderFormData.value.client_phone,
                client_surname: orderFormData.value.client_surname,
                delivery_address: orderFormData.value.delivery_address,
            },
        };

        isSending.value = true;
        isLoaded.value = false;
        resultModal.open();

        try {
            await setTimeout(() => {}, 3000);
            const result = await fetchCreateOrder(sendData);

            setTimeout(() => {
                orderID.value = result.id;
                orderLink.value = `/order-${result.id}-${result.secret_key}`;
                //cartStore.clear();
                isLoaded.value = true;
            }, 500);
        } catch (e) {
            resultModal.close();
            if (e instanceof StandartErrorList) {
                errors.value = e.details;
                errorModal.open();
            }
        } finally {
            isSending.value = false;
            isLoaded.value = false;
        }
    }
};
</script>

<template>
    <div :class="$style.wrapper">
        <div :class="$style.form">
            <div :class="[$style.line, $style.x2]">
                <div>
                    <div :class="$style.label">Имя:</div>
                    <div :class="$style.input"><SharedUiInput v-model="orderFormData.client_name" /></div>
                </div>
                <div>
                    <div :class="$style.label">Фамилия:</div>
                    <div :class="$style.input"><SharedUiInput v-model="orderFormData.client_surname" /></div>
                </div>
            </div>
            <div :class="[$style.line, $style.x2]">
                <div>
                    <div :class="$style.label">Е-mail:</div>
                    <div :class="$style.input"><SharedUiInput v-model="orderFormData.client_email" /></div>
                </div>
                <div>
                    <div :class="$style.label">Телефон:</div>
                    <div :class="$style.input">
                        <SharedUiInput
                            v-model="orderFormData.client_phone"
                            :mask="{ mask: '+{7} (000) 000-00-00' }"
                            :mask-unlazy-on-focus="true"
                            @mask-accept="isPhoneComplete = false"
                            @mask-complete="isPhoneComplete = true"
                        />
                    </div>
                </div>
            </div>
            <div :class="[$style.line]">
                <div>
                    <div :class="$style.label">Адрес:</div>
                    <div :class="$style.input"><SharedUiInput v-model="orderFormData.delivery_address" /></div>
                </div>
            </div>
        </div>
        <div :class="$style.button">
            <button
                class="button_1 big"
                @click="sendForm"
            >
                Оформить заказ
            </button>
        </div>
    </div>
</template>

<style lang="less" module>
@import '@styles/includes';

.form {
    display: flex;
    flex-direction: column;
    gap: 20px;

    > .line {
        display: flex;
        flex-wrap: wrap;
        gap: 20px;

        > div {
            width: 100%;

            > .label {
                margin-bottom: 5px;

                .width-size-sm-less({
                    font-size: 14px;
                });
            }
        }

        &.x2 {
            > div {
                width: calc(50% - 10px);

                .width-size-sm-less({
                    width: 100%;
                });
            }
        }
    }
}

.button {
    margin-top: 30px;
    text-align: center;
}
</style>
