export interface Product {
    productID: string;
    productName: string;
    productDescription: string;
    productPrice: number;
    productImage: string;
    productCategory: string;
    productRating: number;
    productStock: number;
    createdAt: string;
    updatedAt: string;
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
        productStock: raw.productStock as number || 0,
        createdAt: raw.createdAt as string || '',
        updatedAt: raw.updatedAt as string || ''    
    }
}