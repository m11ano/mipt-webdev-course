<script setup lang="ts">
import { StandartErrorList } from '~/shared/errors/errors';
import { fetchUpload } from './api/fetchUpload';
import Error from './modals/Error.vue';

interface UploadedFile {
    id: string;
    url: string;
}

const props = defineProps<{
    mode: 'solo' | 'multi';
    uploadUrl: string;
    acceptTypes?: string;
    modelValue?: UploadedFile[];
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', payload: UploadedFile[]): void;
}>();

const { mode, uploadUrl, acceptTypes = '', modelValue = [] } = props;

const isLoading = ref(false);

const files = ref<UploadedFile[]>([...modelValue]);

watch(
    () => modelValue,
    (newVal) => {
        files.value = newVal ? [...newVal] : [];
    },
    { deep: true },
);

const updateModelValue = () => {
    emit('update:modelValue', [...files.value]);
};

const extractBaseName = (url: string): string => {
    try {
        const parts = url.split('/');
        return parts[parts.length - 1] || url;
    } catch {
        return url;
    }
};

const isImage = (file: UploadedFile) => {
    return /\.(jpe?g|png|gif|webp)$/i.test(file.url);
};

const getFileExtension = (url: string): string => {
    const clean = url.split('?')[0].split('#')[0];
    const lastDot = clean.lastIndexOf('.');
    if (lastDot === -1) return '';
    const lastSlash = clean.lastIndexOf('/');
    if (lastDot < lastSlash) return '';
    return clean.slice(lastDot + 1);
};

const overlay = useOverlay();
const showErrors = (name: string, errors: string[]) => {
    const modal = overlay.create(Error, {
        props: {
            errors,
            name,
        },
        destroyOnClose: true,
    });
    modal.open();
};

const handleFiles = async (event: Event) => {
    const input = event.target as HTMLInputElement;
    if (!input.files || input.files.length === 0) return;
    if (isLoading.value) {
        input.value = '';
        return;
    }

    const selectedFiles = Array.from(input.files);
    isLoading.value = true;

    for (const file of selectedFiles) {
        const formData = new FormData();
        formData.append('file', file);

        try {
            const uploaded = await fetchUpload(uploadUrl, file);

            if (mode === 'solo') {
                files.value = [];
            }

            files.value.push({
                id: uploaded.id,
                url: uploaded.url,
            });
            updateModelValue();
        } catch (e) {
            if (e instanceof StandartErrorList) {
                showErrors(file.name, e.details);
            }
        }
    }

    isLoading.value = false;
    input.value = '';
};

const deleteFile = (index: number) => {
    files.value.splice(index, 1);
    updateModelValue();
};

const moveUp = (index: number) => {
    if (index <= 0) return;
    const temp = files.value[index - 1];
    files.value[index - 1] = files.value[index];
    files.value[index] = temp;
    updateModelValue();
};

const moveDown = (index: number) => {
    if (index >= files.value.length - 1) return;
    const temp = files.value[index + 1];
    files.value[index + 1] = files.value[index];
    files.value[index] = temp;
    updateModelValue();
};
</script>

<template>
    <div>
        <div>
            <div
                v-for="(file, index) in files"
                :key="file.id"
                :class="$style.file_item"
            >
                <div :class="$style.icon">
                    <a
                        :href="file.url"
                        target="_blank"
                    >
                        <template v-if="isImage(file)">
                            <img
                                v-if="isImage(file)"
                                :src="file.url"
                                alt="preview"
                            />
                        </template>
                        <template v-else> {{ getFileExtension(file.url) ? `.${getFileExtension(file.url)}` : '.unknown' }} </template>
                    </a>
                </div>

                <div :class="$style.file_id">
                    <a
                        :href="file.url"
                        target="_blank"
                        >{{ file.id }}</a
                    >
                </div>

                <div :class="$style.actions">
                    <template v-if="mode === 'multi'">
                        <UButton
                            icon="i-lucide-chevron-up"
                            size="xs"
                            color="graylight"
                            :disabled="index === 0"
                            @click="moveUp(index)"
                        />
                        <UButton
                            icon="i-lucide-chevron-down"
                            size="xs"
                            color="graylight"
                            :disabled="index === files.length - 1"
                            @click="moveDown(index)"
                        />
                    </template>
                    <UButton
                        icon="i-lucide-trash"
                        size="xs"
                        color="info"
                        @click="deleteFile(index)"
                    />
                </div>
            </div>
        </div>

        <div>
            <UButton
                as="label"
                class="bg-blue-500 text-white px-4 py-2 rounded cursor-pointer hover:bg-blue-600"
                :loading="isLoading"
                :disabled="isLoading"
            >
                {{ mode === 'multi' ? 'Выбрать файлы' : 'Выбрать файл' }}
                <input
                    type="file"
                    :multiple="mode === 'multi'"
                    :accept="acceptTypes"
                    class="hidden"
                    @change="handleFiles"
                />
            </UButton>
        </div>
    </div>
</template>

<style lang="less" module>
@import '@styles/includes';

.file_item {
    padding: 5px;
    border-radius: 10px;
    background-color: var(--color-neutral-200);
    display: flex;
    gap: 10px;
    align-items: center;
    margin-bottom: 10px;

    > .icon {
        flex-shrink: 0;
        > a {
            display: block;
            width: 50px;
            height: 50px;
            background: #fff;
            display: flex;
            align-items: center;
            justify-content: center;
            overflow: hidden;

            > img {
                max-width: 50px;
                max-height: 50px;
            }
        }
    }

    > .file_id {
        font-size: 11px;

        > a {
            text-decoration: underline;
            color: #253aff;
        }

        .width-size-sm-less({
            display: none
        });
    }

    > .actions {
        flex-shrink: 0;
        margin-left: auto;
        display: flex;
        gap: 5px;
    }
}
</style>
