export interface Product {
    id: string;
    label: string;
    productCode: string;
    url: string;
    price?: string;
    inStock?: boolean;
    description?: string;
}

export interface Alert {
    id: string;
    email: string;
    productCode: string;
    display: string;
}