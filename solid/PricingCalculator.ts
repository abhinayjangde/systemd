import { Product } from "./Order";

export class PricingCalculator {
    calculatePricing(products: Product[]): number {
        return products.reduce((total: number, product: Product)=> total + product.price, 0)
    }
}