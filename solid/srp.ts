import { Invoice } from "./Invoice.js";
import {Product, Order} from "./Order.js"
import { PaymentProcessor } from "./PaymentProcessor.js";
import { PricingCalculator } from "./PricingCalculator.js";

const p1 = new Product('1', 'Laptop', 100);
const p2 = new Product('2', 'iPhone', 500);

const order = new Order();

order.addProduct(p1);
order.addProduct(p2);

const pricingCal = new PricingCalculator();
const total = pricingCal.calculatePricing(order.getProducts())

const invoice = new Invoice();
invoice.generateInvoice(order.getProducts(), total)

const payment = new PaymentProcessor();
payment.processPayment(order)