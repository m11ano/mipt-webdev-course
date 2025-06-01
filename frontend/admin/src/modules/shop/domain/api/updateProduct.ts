import { tryToCatchApiErrors } from '~/shared/errors/errors';
import type { IProductItem } from '../model/types/product';

interface Request {
    name: string;
    is_published: boolean;
    full_description: string;
    price: number;
    image_preview_file_id: string;
    slider_files_ids: string[];
}

const mapDataToRequest = (data: IProductItem): Request => {
    const reqData: Request = {
        name: data.name,
        is_published: data.is_published,
        full_description: data.full_description,
        price: data.price,
        image_preview_file_id: data.image_preview ? data.image_preview.id : '',
        slider_files_ids: data.slider.map((item) => item.id),
    };

    return reqData;
};

export async function updateProduct(data: IProductItem) {
    try {
        await useNuxtApp().$apiFetch(`/products/${data.id}`, {
            method: 'PUT',
            body: mapDataToRequest(data),
        });
    } catch (e: unknown) {
        throw tryToCatchApiErrors(e);
    }
}
