import { Product } from "./Order";

export class Invoice {
    generateInvoice(products: Product[], total: number){
        console.log(`
Invoice Data: ${new Date().toDateString()}
------------------------------
Product Name\tPrice
------------------------------ 
        `);

        products.forEach((product)=>{
            console.log(`${product.name}\t\t${product.price}`);
        })

        console.log(`------------------------------`);
        console.log(`Total: ${total}`)
        console.log(`------------------------------`);
    }
}