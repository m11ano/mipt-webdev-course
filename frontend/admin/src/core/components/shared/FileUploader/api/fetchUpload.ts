import { tryToCatchApiErrors } from '~/shared/errors/errors';

interface Response {
    id: string;
    url: string;
}

export async function fetchUpload(url: string, file: File) {
    const formData = new FormData();
    formData.append('file', file);

    try {
        return await useNuxtApp().$apiFetch<Response>(url, {
            method: 'POST',
            body: formData,
        });
    } catch (e: unknown) {
        throw tryToCatchApiErrors(e);
    }
}
