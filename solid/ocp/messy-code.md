```js
// Open-Close Principle

class PaymentProcessor {
    processPayment(amount: number, paymentType: string){

        if (paymentType === 'paypal'){
            console.log(`processing payment of ${amount} using Paypal`)
        }
        else if (paymentType === 'razorpay'){
            console.log(`processing payment of ${amount} using Razorpay`)
        }
        else if(paymentType === 'creditcard'){
            console.log(`processing payment of ${amount} using Credi Card`)
        }
        else {
            throw new Error("unknown payment type")
        }
    }
}

const processor = new PaymentProcessor()

processor.processPayment(100, 'paypal');


```