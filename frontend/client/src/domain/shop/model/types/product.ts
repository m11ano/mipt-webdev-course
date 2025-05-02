export interface IProductListItem {
    id: number;
    name: string;
    price: number;
    stock_available: number;
    image_preview: string;
}

export interface IProductItem {
    id: number;
    name: string;
    is_published: boolean;
    full_description: string;
    price: number;
    stock_available: number;
    image_preview: string;
    slider: string[];
}
