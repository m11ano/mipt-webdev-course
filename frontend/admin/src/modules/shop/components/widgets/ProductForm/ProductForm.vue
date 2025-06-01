<script setup lang="ts">
import type { IProductItem } from '~/modules/shop/domain/model/types/product';

const props = defineProps<{ disabled?: boolean; mode: 'new' | 'edit' }>();

const dataModel = defineModel<IProductItem>({ required: true });
</script>

<template>
    <div class="form-table">
        <div>
            <div class="title">Название:</div>
            <div class="value">
                <UInput
                    v-model="dataModel.name"
                    size="xl"
                    class="w-full"
                    :disabled="disabled"
                />
            </div>
        </div>
        <div>
            <div class="title">Опубликовано?</div>
            <div class="value">
                <USwitch
                    v-model="dataModel.is_published"
                    size="xl"
                    color="success"
                    :disabled="disabled"
                />
            </div>
        </div>
        <div>
            <div class="title">Цена:</div>
            <div class="value">
                <UInputNumber
                    v-model="dataModel.price"
                    size="xl"
                    :disabled="disabled"
                    :format-options="{
                        style: 'currency',
                        currency: 'RUB',
                        currencyDisplay: 'code',
                        currencySign: 'accounting',
                    }"
                    :min="0"
                    orientation="vertical"
                />
            </div>
        </div>
        <div v-if="mode === 'new'">
            <div class="title">Доступный остаток:</div>
            <div class="value">
                <UInputNumber
                    v-model="dataModel.stock_available"
                    size="xl"
                    :disabled="disabled"
                    orientation="vertical"
                    :min="0"
                />
            </div>
        </div>
        <div>
            <div class="title">Основная фотография:</div>
            <div class="value">
                <SharedFileUploader
                    mode="solo"
                    upload-url="products/image?image_type=preview"
                    accept-types="image/*"
                    :model-value="dataModel.image_preview ? [dataModel.image_preview] : []"
                    @update:model-value="dataModel.image_preview = $event.length ? $event[0] : null"
                />
            </div>
        </div>
        <div>
            <div class="title">Слайдер фотографий:</div>
            <div class="value">
                <SharedFileUploader
                    v-model="dataModel.slider"
                    mode="multi"
                    upload-url="products/image?image_type=slider"
                    accept-types="image/*"
                />
            </div>
        </div>
        <div>
            <div class="title">Описание товара:</div>
            <div class="value">
                <UTextarea
                    v-model="dataModel.full_description"
                    size="xl"
                    class="w-full"
                    :rows="8"
                    :disabled="disabled"
                />
            </div>
        </div>
    </div>
</template>

<style lang="less" module>
@import '@styles/includes';
</style>
