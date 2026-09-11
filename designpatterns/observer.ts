interface Observer {
    update(data:any): void;
}

class NewsAgency {
    private observers: Observer[] = [];
    private news: string[] = [];

    subcribe(observer: Observer){
        this.observers.push(observer);
    }

    unsubscribe(observer: Observer){
        this.observers = this.observers.filter((obs)=> obs !== observer);
    }

    notify(data: any){
        for(const observer of this.observers){
            observer.update(data);
        }
    }

    publishNews(data:any){
        this.news.push(data);
        this.notify(data);
    }
}

class EmailNotification implements Observer {
    update(data: any): void {
        console.log("Email Notification: ", data)
    }
}

class SMSNotification implements Observer {
    update(data: any): void {
        console.log("SMS Notification: ", data);
    }
}
class WhatsappNotification implements Observer {
    update(data: any): void {
        console.log("Whatsapp Notification: ", data);
    }
}

const newsAgency = new NewsAgency();

const smsNoti = new SMSNotification();
const emailNoti = new EmailNotification();
const whatsappNoti = new WhatsappNotification();

newsAgency.subcribe(smsNoti);
newsAgency.subcribe(emailNoti);
newsAgency.subcribe(whatsappNoti);

newsAgency.unsubscribe(smsNoti);

newsAgency.publishNews("Indian wins todays cricket match")