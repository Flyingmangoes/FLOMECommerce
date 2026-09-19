export interface Product {
    productID: string;
    productName: string;
    productDescription: string;
    productPrice: number;
    productImage: string;
    productCategory: string;
    productRating: number;
    productTotalStock: number;
    productExtra: ProductExtra[];
    createdAt: string;
    updatedAt: string;
}

export interface ProductExtra {
    color: string;
    size: string;
    stock: number;
}

export function mapProductResponse(raw: Record<string, unknown>): Product {
    return {
        productID: raw.productID as string || '',   
        productName: raw.productName as string || '',
        productDescription: raw.productDescription as string || '',
        productPrice: raw.productPrice as number || 0,
        productImage: raw.productImage as string || '',
        productCategory: raw.productCategory as string || '',
        productRating: raw.productRating as number || 0.0,
        productTotalStock: raw.productTotalStock as number || 0,
        productExtra: (raw.productExtra as Record<string, unknown>[] || []).map(item => ({
            color: item.Color as string || '',
            size: item.Color as string || '',
            stock: item.Stock as number || 0,
        })),
        createdAt: raw.createdAt as string || '',
        updatedAt: raw.updatedAt as string || ''    
    }
}