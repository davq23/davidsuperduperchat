import Message from "../model/Message.js";
import View from "./View.js";



export default class ChatView extends View {
    constructor(anchor) {
        super(anchor);

        this.ws = null
        
        this.formHTML();

        const self = this;

        this.form.onsubmit = function(event) {
            event.preventDefault();
        
            const message = {
                body: self.messageInput.value
            }; 

            if (message.body !== "") {
                try {
                    self.ws.send(JSON.stringify(message));
            
                    self.messageInput.value = "";
                        
                } catch(err) {
                    self.chat.innerHTML += `<h3>An error occured!</h3>`;
                    console.log(err);
                }
            }
        }
    }

    applyListeners() {
        
    }

    formHTML() {
        this.chatDiv = document.createElement("div");
        this.chat = document.createElement("div");
        this.chat.id = "chat";
        
        this.form = document.createElement("form"); 
        this.form.id = "message-form"
        this.messageInput = document.createElement("input"); 
        this.messageInput.id = "message"

        this.submitButton = document.createElement('button');
        this.submitButton.innerText = "Send"
        this.submitButton.type = 'submit';

        this.closeConMessage = document.createElement('h3');
        this.submitButton.innerText = ""

        this.form.appendChild(this.messageInput);
        this.form.appendChild(this.submitButton);
        this.chatDiv.appendChild(this.chat);
        this.chatDiv.appendChild(this.form);
        this.chatDiv.appendChild(this.closeConMessage);
        this.fragment.appendChild(this.chatDiv);

        return this.chatDiv;
    }

    initConnection() {
        const self = this;
        return new Promise(function(resolve, reject) {
            self.ws = new WebSocket("wss://infinite-ocean-38389.herokuapp.com/chat");
            

            self.ws.onopen = function() {
                this.send(JSON.stringify({body: "Hello people!!"}))
                this.closeConMessage.innerText = ""
                resolve()
            }

            self.ws.onclose = function(event) {
                this.closeConMessage.innerText = "Connection closed, please refresh to reconnect"
            }
            
            self.ws.onmessage = function(event) {
                const message = JSON.parse(event.data);

                const messageElement = new Message(message.sendername, message.body, message.sendat, message.receivedat);           

                const node = messageElement.asHTMLNode();

                if (self.chat.scrollTop + self.chat.clientHeight >= self.chat.scrollHeight) {
                    self.chat.appendChild(node);        
                    self.chat.scrollTop = node.offsetHeight + node.offsetTop;
                } else {
                    self.chat.appendChild(node);        
                }
            }

            self.ws.onerror = function(err) {
                self.chat.innerHTML += `<h3 class="message">An error occured!</h3>`;
                reject(err)
            }
            
            
        })
        
        
    }

    render() {
        this.anchor.innerHTML = '';
        this.anchor.appendChild(this.fragment);
    }
}