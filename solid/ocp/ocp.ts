// Open-Close Principle

interface IPaymentProcessor {
    processPayment(amount: number): void
}

class PaymentProcessor implements IPaymentProcessor {

    private processor: IPaymentProcessor;
    constructor(paymentProcessor: IPaymentProcessor){
        this.processor = paymentProcessor
    }

    processPayment(amount: number): void{
        this.processor.processPayment(amount);
    }
}

class PaypalGateway implements IPaymentProcessor {
    processPayment(amount: number): void {
        console.log(`processing payment of ${amount} using Paypal`);
    }
}

class RazorpayGateway implements IPaymentProcessor {
    processPayment(amount: number): void {
        console.log(`processing payment of ${amount} using Razorpay`);
    }
}

const processor = new PaymentProcessor(new RazorpayGateway());

processor.processPayment(100);

