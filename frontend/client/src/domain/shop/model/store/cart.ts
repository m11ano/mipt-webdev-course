import { defineStore } from 'pinia';
import type { ICartItem } from '../types/cart_item';

const CART_LOCAL_STORAGE_KEY = 'shop_cart';

export const useCartStore = defineStore('shop_cart', {
    state: () => ({
        items: [] as ICartItem[],
    }),
    getters: {
        count(): number {
            return this.items.length;
        },
        total(): number {
            return this.items.reduce((acc, i) => acc + i.quantity, 0);
        },
        productIDs(): number[] {
            return this.items.map((i) => i.id);
        },
    },
    actions: {
        getProductQuantity(id: number): number {
            const item = this.items.find((i) => i.id === id);
            return item ? item.quantity : 0;
        },
        add(id: number, quantity: number = 1) {
            const existing = this.items.find((i) => i.id === id);
            if (existing) {
                existing.quantity += quantity;
                if (existing.quantity < 1) {
                    existing.quantity = 1;
                }
            } else if (quantity > 0) {
                this.items.push({
                    id,
                    quantity,
                });
            }
            this.saveToLocalStorage();
        },
        setQuantity(id: number, quantity: number) {
            const i = this.items.findIndex((i) => i.id === id);
            if (i !== undefined) {
                this.items[i].quantity = quantity;
            }
            this.saveToLocalStorage();
        },
        filterByIDs(ids: number[]) {
            this.items = this.items.filter((i) => ids.includes(i.id));
            this.saveToLocalStorage();
        },
        remove(id: number) {
            this.items = this.items.filter((i) => i.id !== id);
            this.saveToLocalStorage();
        },
        clear() {
            this.items = [];
            this.saveToLocalStorage();
        },
        saveToLocalStorage() {
            if (import.meta.server) {
                return;
            }
            localStorage.setItem(CART_LOCAL_STORAGE_KEY, JSON.stringify(this.items));
        },
        loadFromLocalStorage() {
            if (import.meta.server) {
                return;
            }
            const data = localStorage.getItem(CART_LOCAL_STORAGE_KEY);
            if (data) {
                this.items = JSON.parse(data);
            }
        },
    },
});
