<script setup lang="ts">
import { doLogout } from '~/plugins/auth/model';

const { smLess } = useScreenSizes();

const { $authData } = useNuxtApp();
</script>

<template>
    <header id="auth_header">
        <div>
            <div class="left">
                <div class="logo"><span>&nbsp;</span></div>
                <div class="menu">
                    <UDrawer direction="left">
                        <UButton
                            label="Меню"
                            :size="smLess ? 'sm' : 'md'"
                            icon="i-lucide-menu"
                        />

                        <template #content>
                            <div :id="$style.draw_menu"><WidgetMenu /></div>
                        </template>
                    </UDrawer>
                </div>
            </div>
            <div class="right">
                <div class="greet">Привет, {{ $authData.userData?.name }}!</div>
                <div class="logout">
                    <UButton
                        :size="smLess ? 'sm' : 'md'"
                        icon="i-lucide-log-out"
                        variant="subtle"
                        @click="doLogout"
                    >
                        <template v-if="!smLess">Выход</template>
                    </UButton>
                </div>
                <!-- <AButton @click="$toogleAppTheme">Смена темы</AButton> -->
            </div>
        </div>
    </header>
</template>

<style lang="less" module>
@import '@styles/includes';

#draw_menu {
    padding: 15px 0 15px 15px;
    width: 250px;
}
</style>
