<script setup lang="ts">
const props = defineProps<{
    slider: string[];
}>();

const nowSlideSelected = ref(0);

const nowSlideItem = computed(() => {
    return props.slider.find((_, i) => i == nowSlideSelected.value) || null;
});

const next = () => {
    nowSlideSelected.value = (nowSlideSelected.value + 1) % props.slider.length;
    //showNextTimeout();
};

let timeout: ReturnType<typeof setTimeout> | null = null;

const showNextTimeout = () => {
    if (timeout) {
        clearTimeout(timeout);
    }
    timeout = setTimeout(() => {
        next();
    }, 4000);
};

const clearNextTimeout = () => {
    if (timeout) {
        clearTimeout(timeout);
    }
};

onMounted(() => {
    //showNextTimeout();
});

onUnmounted(() => {
    clearNextTimeout();
});
</script>

<template>
    <div :class="$style.wrapper">
        <div
            :class="$style.img"
            @click="next"
        >
            <img
                v-if="nowSlideItem"
                :src="nowSlideItem"
                alt=""
            />
            <span></span>
        </div>
        <div :class="$style.items">
            <template
                v-for="(item, i) in slider"
                :key="item"
            >
                <button
                    :class="i == nowSlideSelected ? $style.active : false"
                    @click="nowSlideSelected = i"
                >
                    <img
                        v-if="item"
                        :src="item"
                        alt=""
                        loading="lazy"
                    />
                    <span></span>
                </button>
            </template>
        </div>
    </div>
</template>

<style lang="less" module>
@import '@styles/includes';

.img {
    position: relative;
    background-color: var(--ui-bg-color);
    border-radius: 10px;
    overflow: hidden;
    box-shadow: 0px 0px 10px rgba(0, 0, 0, 0.01);
    cursor: pointer;
    -webkit-tap-highlight-color: transparent;

    > img {
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        object-fit: cover;
    }

    > span {
        display: block;
        padding-top: 100%;
    }
}

.items {
    display: flex;
    flex-wrap: wrap;
    gap: 20px;
    margin-top: 20px;

    > button {
        display: block;
        width: calc((100% - 80px) / 5);
        position: relative;
        background-color: var(--ui-bg-color);
        border-radius: 5px;
        overflow: hidden;
        box-shadow: 0px 0px 10px rgba(0, 0, 0, 0.01);
        transition: box-shadow 0.3s ease;
        -webkit-tap-highlight-color: transparent;

        > img {
            position: absolute;
            top: 0;
            left: 0;
            width: 100%;
            height: 100%;
            object-fit: cover;
        }

        > span {
            display: block;
            padding-top: 100%;
        }

        .hover({
            box-shadow: 0px 0px 10px rgba(0, 0, 0, 0.1);
        });

        &.active {
            outline: rgba(21, 93, 253, 0.2) solid 2px;
            outline-offset: 5px;
        }
    }
}
</style>
