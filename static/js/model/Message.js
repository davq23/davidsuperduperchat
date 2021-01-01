

export default class Message {
    constructor(senderName, body, sendAt, receivedAt) {
        this.senderName = senderName
        this.body = body
        this.sendAt = sendAt
        this.receivedAt = receivedAt
    }

    asHTMLNode() {
        const messageElement = document.createElement("div");
        messageElement.classList.add("message");

        const sentAt = document.createElement("i");
        sentAt.innerText = `Sent at ${this.sendAt}`;

        const sender = document.createElement("b");
        sender.innerText = this.senderName;

        const body = document.createElement("p");
        body.innerText = this.body;

        const recAt = document.createElement("i");
        recAt.innerText = `Received at ${this.receivedAt}`;

        messageElement.appendChild(sender);
        messageElement.appendChild(document.createElement('br'));
        messageElement.appendChild(sentAt);       
        messageElement.appendChild(body);
        messageElement.appendChild(recAt);               

        return messageElement;
    }
}