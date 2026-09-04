export interface OrderModel {
    orderID: string;
    buyerID: string;
    buyerEmail: string;
    totalPrice: number;
    location: string;
    status: string;
    createdAt: string;
    orderItems: OrderItemModel[];
}

export interface OrderItemModel {
    orderItemID: string;
    orderID: string;
    productID: string;
    quantity: number;
    price: number;
}

export function mapOrderResponse(raw: Record<string, unknown>): OrderModel {
    return {
        orderID: raw.orderID as string || '',
        buyerID: raw.buyerID as string || '',
        buyerEmail: raw.buyerEmail as string || '',
        totalPrice: raw.totalPrice as number || 0,
        location: raw.location as string || '',
        status: raw.status as string || '',
        createdAt: raw.createdAt as string || '',
        orderItems: (raw.orderItems as Record<string, unknown>[] || []).map(item => ({
            orderItemID: item.orderItemID as string || '',
            orderID: item.orderID as string || '',
            productID: item.productID as string || '',
            quantity: item.quantity as number || 0,
            price: item.price as number || 0
        }))
    }
}