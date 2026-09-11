class Singleton{
    static instance = null;
    constructor(){
        console.log("singleton constructor called");
    }

    static getInstance(){
        if(this.instance === null){
            this.instance = new Singleton();
        }
        return this.instance;
    }
}

s1 = Singleton.getInstance()
s2 = Singleton.getInstance()

console.log(s1===s2)
