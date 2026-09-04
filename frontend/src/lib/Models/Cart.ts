export interface CartModel {
    cartID: string;
    userID: string;
    cartItems: CartItemModel[];
}

export interface CartItemModel {
    cartItemID: string;
    cartID: string;
    productID: string;
    quantity: number;
}

export function mapCartResponse(raw: Record<string, unknown>): CartModel {
    return {
        cartID: raw.cartID as string || '',
        userID: raw.userID as string || '',
        cartItems: (raw.cartItems as Record<string, unknown>[] || []).map(item => ({
            cartItemID: item.cartItemID as string || '',
            cartID: item.cartID as string || '',
            productID: item.productID as string || '',
            quantity: item.quantity as number || 0
        }))
    }
}