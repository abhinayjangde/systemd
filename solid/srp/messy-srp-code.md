
Order.ts

```js
// Single Responsibility Principle (SRP)

export class Product {
    id: string;
    name: string;
    price: number;

    constructor(id: string, name: string, price: number){
        this.id = id;
        this.name = name;
        this.price = price;
    }
}

export class Order {
    products: Product[] = [];

    addProduct(product: Product){
        this.products.push(product);
    }

    getProducts(){
        return this.products;
    }

    removeProduct(productId: string){
        this.products = this.products.filter(
            (prd) => productId !== prd.id
        )
    }

    calculatePricing(){
        return this.products.reduce((total, product)=> total + product.price, 0)
    }

    generateInvoice(){
        console.log(`
Invoice Data: ${new Date().toDateString()}
------------------------------
Product Name\tPrice
------------------------------ 
        `);

        this.products.forEach((product)=>{
            console.log(`${product.name}\t\t${product.price}`);
        })

        console.log(`------------------------------`);
        console.log(`Total: ${this.calculatePricing()}`)
        console.log(`------------------------------`);
    }

    processPayment(){
        console.log('Processing payment...');
        console.log('Payment processed successfully');
        console.log('Sending notification...');
    }
}
```

srp.ts
```js
import {Product, Order} from "./Order.js"

const p1 = new Product('1', 'Laptop', 100);
const p2 = new Product('2', 'iPhone', 500);

const order = new Order();

order.addProduct(p1);
order.addProduct(p2);

order.generateInvoice();

order.processPayment();
```