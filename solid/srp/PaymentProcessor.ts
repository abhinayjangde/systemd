import { Order } from "./Order";

export class PaymentProcessor {
     processPayment(order: Order){
        console.log('Processing payment...');
        console.log('Payment processed successfully');
        console.log('Sending notification...');
    }
}