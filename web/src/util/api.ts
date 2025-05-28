import axios, { HttpStatusCode, type AxiosResponse } from "axios"
import type { Alert, Product } from "./types"

export async function getProducts(): Promise<Product[]> {
    let products: Product[] = []
    await axios.get("products").then((response: AxiosResponse<Product[]>) => {
        products = response.data
    }).catch((err) => {
        console.log(err)
    })
    return products
}

export async function getPartsForProduct(id: string): Promise<Product[]> {
    let parts: Product[] = []
    await axios.get(`parts/${id}`).then((response: AxiosResponse<Product[]>) => {
        parts = response.data
    }).catch((err) => {
        console.log(err)
    })
    return parts
}

export async function getPartDetails(p: Product): Promise<Product | null> {
    let details: Product | null = null
    await axios.get(`partinfo/${p.id}`).then((response: AxiosResponse<Product>) => {
        details = response.data
    }).catch((err) => {
        console.log(err)
    })
    return details
}

export async function putAlert(email: string, id: string): Promise<HttpStatusCode> {
    const response = await axios.put('alerts/', { email, partID: id })
    return response.status
}

export async function getAlertsByEmail(email: string): Promise<Alert[]> {
    let alerts: Alert[] = [];
    await axios.get(`alerts/${email}`).then((response: AxiosResponse<Alert[]>) => {
        alerts = response.data
    }).catch((err) => {
        console.log(err)
    })
    return alerts
}