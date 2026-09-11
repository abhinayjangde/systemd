#include<iostream>
#include<vector>
#include<algorithm>

using namespace std;

class ISubscriber {
    public:
        virtual void update() = 0;
        virtual ~ISubscriber() {} // virtual destructor for interface
};

class IChannel {
    public:
        virtual void subscribe(ISubscriber* subscriber) = 0;
        virtual void unsubscribe(ISubscriber* subscriber) = 0;
        virtual void notifySubscribers() = 0;
        virtual ~IChannel() {}
};

class Channel : public IChannel {
    private:
        vector<ISubscriber*> subscribers;
        string name;
        string latestVideo;
    public:
        Channel(const string& name){
            this->name = name;
        }

        void subscribe(ISubscriber* subscriber) override {
            if(find(subscribers.begin(), subscribers.end(), subscriber) == subscribers.end()){
                subscribers.push_back(subscriber);
            }
        }

        void unsubscribe(ISubscriber* subscriber) override {
            auto it = find(subscribers.begin(), subscribers.end(), subscriber);
            if(it != subscribers.end()){
                subscribers.erase(it);
            }
        }

        void notifySubscribers() override {
            for(ISubscriber* sub : subscribers){
                sub->update();
            }
        }

        void uploadVideo(const string& title){
            latestVideo = title;
            cout << "\n[" << name << " uploaded \"" << title << "\"]\n";
            notifySubscribers();
        }

        string getVideoData(){
            return "\nCheckout our new video : " + latestVideo + "\n";
        }
};

class Subscriber : public ISubscriber {
    private:
        string name;
        Channel* channel;
    public:
        Subscriber(const string& name, Channel* channel) {
            this->name = name;
            this->channel = channel;
        }

        void update() override {
            cout << "hey "<< name << "," << this->channel->getVideoData();
        }
};
int main() {
    cout << "Observer Design Pattern Example" << endl;
    
    Channel * channel = new Channel("Codebhaiya");

    Subscriber* subs1 = new Subscriber("Aditi", channel);
    Subscriber* subs2 = new Subscriber("Arya", channel);
    Subscriber* subs3 = new Subscriber("Preeti", channel);

    channel->subscribe(subs1);
    channel->subscribe(subs2);
    channel->subscribe(subs3);

    channel->unsubscribe(subs1);
    
    channel->uploadVideo("Observer pattern tutorial");

    return 0;
}