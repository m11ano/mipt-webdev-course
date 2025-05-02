<script setup lang="ts"></script>

<template>
    <div>
        <div :id="$style.block_desc">
            <div :class="$style.title">Магазин лампочек СВЕТЗАВОД</div>
            <div :class="$style.text">
                Добро пожаловать в СВЕТЗАВОД — ваш надежный партнер в мире освещения! У нас вы найдете лампочки всех типов: от классических
                и светодиодных до умных решений для «умного дома».
                <br /><br />
                Мы предлагаем качественную продукцию от ведущих производителей по выгодным ценам, с доставкой по всей России.
            </div>
        </div>
        <div :id="$style.block_catalog">
            <div :class="$style.title">Каталог</div>
            <WidgetProductsList :class="$style.content" />
        </div>
    </div>
</template>

<style lang="less" module>
@import '@styles/includes';

#block_desc {
    .std-wrapper();

    > .title {
        font-size: 32px;
        text-align: center;
        font-family: 'Strong';
        color: var(--font-color-3);

        .width-size-sm-less({
            font-size:24px;
        });
    }

    > .text {
        margin-top: 35px;
        line-height: 1.5;
        font-size: 14px;
        text-align: justify;

        .width-size-sm-less({
            font-size:12px;
            margin-top: 20px;
        });
    }
}

#block_catalog {
    .std-wrapper(40px, 0, 20px, 0);

    > .title {
        font-size: 24px;
        text-align: center;
        font-family: 'Strong';
        color: var(--font-color-3);

        .width-size-sm-less({
            font-size:20px;
        });
    }

    > .content {
        margin-top: 40px;

        .width-size-sm-less({
            margin-top: 20px;
        });
    }
}
</style>
