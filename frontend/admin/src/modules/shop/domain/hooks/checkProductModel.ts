import type { IProductItem } from '../model/types/product';

export function checkProductModel(data: IProductItem): string[] {
    const result: string[] = [];

    if (data.name.length === 0) {
        result.push('Название не указано');
    }

    if (data.image_preview === null) {
        result.push('Загрузите основную фотографию');
    }

    if (data.slider.length === 0) {
        result.push('Загрузите фотографии для слайдера');
    }

    return result;
}
