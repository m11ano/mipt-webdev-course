export type ImageType = 'preview' | 'slider';

export interface IProductListItem {
    id: number;
    image_preview: string;
    is_published: boolean;
    name: string;
    price: number;
    stock_available: number;
}

export interface IProductItem {
    id: number;
    name: string;
    is_published: boolean;
    full_description: string;
    price: number;
    stock_available: number;
    image_preview: {
        id: string;
        url: string;
    } | null;
    slider: {
        id: string;
        url: string;
    }[];
}
