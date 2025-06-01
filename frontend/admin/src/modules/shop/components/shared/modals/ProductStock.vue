<script setup lang="ts">
import type { RadioGroupItem } from '@nuxt/ui';
import { updateProductStock } from '~/modules/shop/domain/api/updateProductStock';
import { StandartErrorList } from '~/shared/errors/errors';

const props = defineProps<{
    productID: number;
}>();

const emit = defineEmits<{ close: [boolean] }>();

const operationItems = ref<RadioGroupItem[]>([
    { label: 'Увеличить остаток на:', value: 'increase' },
    { label: 'Уменьшить остаток на:', value: 'decrease' },
]);
const operation = ref<'increase' | 'decrease'>('increase');

const value = ref(0);

const isLoading = ref(false);
const errors = ref<string[]>([]);

const toast = useToast();

const save = async () => {
    if (isLoading.value) return;

    errors.value = [];

    if (value.value < 1) {
        errors.value.push('Значение не указано');
        return;
    }

    isLoading.value = true;
    try {
        await updateProductStock(props.productID, {
            operation: operation.value,
            value: value.value,
        });

        toast.add({
            title: 'Успех',
            description: 'Остаток изменён',
            color: 'success',
            icon: 'i-lucide-check-circle',
        });

        emit('close', true);
    } catch (e) {
        if (e instanceof StandartErrorList) {
            errors.value = e.details;
        }
    } finally {
        isLoading.value = false;
    }
};
</script>

<template>
    <UModal
        :title="`Изменение свободного остатка на складе`"
        @close="emit('close', false)"
    >
        <template #body>
            <div :class="$style.wrapper">
                <div>
                    <URadioGroup
                        v-model="operation"
                        :items="operationItems"
                    />
                </div>
                <div>
                    <UInputNumber
                        v-model="value"
                        size="xl"
                        orientation="vertical"
                        :min="0"
                    />
                </div>
                <div v-if="errors.length">
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
            </div>
        </template>
        <template #footer>
            <div class="flex justify-end">
                <UButton
                    :label="operation === 'increase' ? 'Добавить' : 'Уменьшить'"
                    color="primary"
                    :loading="isLoading"
                    :disabled="isLoading"
                    @click="save"
                />
            </div>
        </template>
    </UModal>
</template>

<style lang="less" module>
@import '@styles/includes';

.wrapper {
    display: flex;
    flex-direction: column;
    gap: 20px;
}
</style>
