document.onreadystatechange = function () {
    switch (document.readyState) {
        case 'complete':
            // App anchor
            var AppDiv = document.getElementById("App");

            // Nav buttons
            var loginButton = document.getElementById("login-button");
            var logoutButton = document.getElementById("logout-button");
            var signupButton = document.getElementById("signup-button");

            // Views
            var signupView = new App.Views.SignupForm(AppDiv);
            var loginView = new App.Views.LoginForm(AppDiv);
            var loadingView = new App.Views.Loading(AppDiv);
            var chatView = null;

            window.addEventListener('change-nav', function (event) {
                loginButton.disabled = event.detail;
                signupButton.disabled = event.detail;
                logoutButton.disabled = !event.detail;

                if (event.detail) {
                    loginButton.classList.add('hidden');
                    signupButton.classList.add('hidden');
                    logoutButton.classList.remove('hidden');
                } else {
                    loginButton.classList.remove('hidden');
                    signupButton.classList.remove('hidden');
                    logoutButton.classList.add('hidden');
                }
            }, false);

            // App events
            AppDiv.addEventListener('loading-', function (event) {
                event.preventDefault();
                loadingView.render();
            })

            AppDiv.addEventListener('chat', function (event) {
                loadingView.render();

                let evt = null

                // Try to open websocket
                chatView = new ChatView(AppDiv);

                chatView.initConnection().then(function () {
                    evt = new CustomEvent("change-nav", {
                        detail: true
                    });

                    // If successful, render chat and hide buttons
                    chatView.render();

                    window.dispatchEvent(evt);

                }).catch(function (err) {
                    console.log(err);

                    evt = new CustomEvent("change-nav", {
                        detail: false
                    });

                    // Render login and load buttons
                    loginView.render();


                    // onclick listeners
                    loginButton.onclick = function () {
                        loadingView.render();

                        loginView.render();
                    }

                    signupButton.onclick = function () {
                        loadingView.render();

                        signupView.render();
                    }

                    logoutButton.onclick = async function () {
                        loadingView.render();

                        chatView.logout();

                        loginView.render();

                        window.dispatchEvent(new CustomEvent("change-nav", {
                            detail: false
                        }));

                    }


                    window.dispatchEvent(evt);
                });

                var event = new Event('chat');
                AppDiv.dispatchEvent(event);



                break;

            });
    }
}


var App = {
    Utils: {
        /**
         * Aplica herencia a una clase
         * 
         * @param {Object} derived 
         * @param {Object} base 
         */
        inherit: function (derived, base) {
            derived.prototype = Object.create(base.prototype);
            derived.prototype.constructor = derived;
        }
    },

    Models: {
        Message: function () {
            this.senderName = senderName
            this.body = body
            this.sendAt = sendAt
            this.receivedAt = receivedAt
        }
    },

    Views: {
        View: function (anchor) {
            this.fragment = document.createDocumentFragment();

            if (anchor instanceof HTMLElement) {
                this.anchor = anchor;
            }
        },

        Chat: function (anchor) {
            this.ws = null

            this.formHTML();

            var self = this;

            // Message submit function
            this.form.onsubmit = function (event) {
                event.preventDefault();

                var message = {
                    body: self.messageInput.value
                };

                if (message.body !== "") {
                    try {
                        self.ws.send(JSON.stringify(message));

                        self.messageInput.value = "";

                    } catch (err) {
                        self.chat.innerHTML += `<h3>An error occured!</h3>`;
                        console.log(err);
                    }
                }
            }
        },

        Loading: function (anchor) {
            App.Views.View.call(this, anchor);
            this.loadingText = document.createElement('h3');
            this.loadingText.innerText = 'Loading...';
            this.loadingText.classList.add('fade');

            this.fragment.appendChild(this.loadingText);
        },

        LoginForm: function (anchor) {
            App.Views.View.call(this, anchor);

            this.loginForm = document.createElement('form');
            this.loginForm.classList.add('form');

            this.title = document.createElement('h3');
            this.title.innerText = 'Login'

            this.err = document.createElement('h3');
            this.err.classList.add('form-error');
            this.err.id = 'err-login';

            var usernameLabel = document.createElement('label');
            usernameLabel.innerText = 'Username';
            usernameLabel.for = 'username';
            usernameLabel.classList.add('form-label');

            this.usernameInput = document.createElement('input');
            this.usernameInput.name = 'username';
            this.usernameInput.placeholder = 'Username';
            this.usernameInput.classList.add('form-input');

            var passwordLabel = document.createElement('label');
            passwordLabel.innerText = 'Password';
            passwordLabel.for = 'password';
            passwordLabel.classList.add('form-label');

            this.passwordInput = document.createElement('input');
            this.passwordInput.type = 'password';
            this.passwordInput.name = 'password';
            this.passwordInput.placeholder = 'Password';
            this.passwordInput.classList.add('form-input');

            this.submitButton = document.createElement('button');
            this.submitButton.innerText = "Send"
            this.submitButton.type = 'submit';
            this.submitButton.classList.add('btn-main');

            this.loginForm.appendChild(this.title);
            this.loginForm.appendChild(this.err);
            this.loginForm.appendChild(usernameLabel);
            this.loginForm.appendChild(document.createElement('br'));
            this.loginForm.appendChild(this.usernameInput);
            this.loginForm.appendChild(document.createElement('br'));
            this.loginForm.appendChild(passwordLabel);
            this.loginForm.appendChild(document.createElement('br'));
            this.loginForm.appendChild(this.passwordInput);
            this.loginForm.appendChild(document.createElement('br'));
            this.loginForm.appendChild(this.submitButton);

            this.fragment.appendChild(this.loginForm)
        },

        SignupForm: function (anchor) {
            App.Views.View.call(this, anchor);

            this.signupForm = document.createElement('form');
            this.signupForm.classList.add('form');

            this.title = document.createElement('h3');
            this.title.innerText = 'Signup'

            this.err = document.createElement('h3');
            this.err.classList.add('form-error');
            this.err.id = 'err-signup';

            var usernameLabel = document.createElement('label');
            usernameLabel.innerText = 'Username';
            usernameLabel.for = 'username';
            usernameLabel.classList.add('form-label');

            var passwordLabel = document.createElement('label');
            passwordLabel.innerText = 'Password';
            passwordLabel.for = 'password';
            passwordLabel.classList.add('form-label');

            this.usernameInput = document.createElement('input');
            this.usernameInput.name = 'username';
            this.usernameInput.placeholder = 'Username';
            this.usernameInput.classList.add('form-input');

            this.passwordInput = document.createElement('input');
            this.passwordInput.type = 'password';
            this.passwordInput.name = 'password';
            this.passwordInput.placeholder = 'Password';
            this.passwordInput.classList.add('form-input');

            this.submitButton = document.createElement('button');
            this.submitButton.innerText = "Send"
            this.submitButton.type = 'submit';
            this.submitButton.classList.add('btn-main');

            this.signupForm.appendChild(this.title);
            this.signupForm.appendChild(this.err);
            this.signupForm.appendChild(usernameLabel);
            this.signupForm.appendChild(document.createElement('br'));
            this.signupForm.appendChild(this.usernameInput);
            this.signupForm.appendChild(document.createElement('br'));

            this.signupForm.appendChild(passwordLabel);
            this.signupForm.appendChild(document.createElement('br'));
            this.signupForm.appendChild(this.passwordInput);
            this.signupForm.appendChild(document.createElement('br'));
            this.signupForm.appendChild(this.submitButton);

            this.fragment.appendChild(this.signupForm)
        }
    }
};

App.Models.Message.asHTMLNode = function () {
    var messageElement = document.createElement("div");
    messageElement.classList.add("message");

    var sentAt = document.createElement("i");
    sentAt.innerText = `Sent at ${this.sendAt}`;
    sentAt.classList.add("meta");

    var sender = document.createElement("b");
    sender.innerText = this.senderName;


    var body = document.createElement("p");
    body.innerText = this.body;


    var recAt = document.createElement("i");
    recAt.innerText = `Received at ${this.receivedAt}`;
    recAt.classList.add("meta");

    messageElement.appendChild(sender);
    messageElement.appendChild(document.createElement('br'));
    messageElement.appendChild(sentAt);
    messageElement.appendChild(body);
    messageElement.appendChild(recAt);

    return messageElement;
}

David.Utils.inherit(App.Views.Chat, App.Views.View);
David.Utils.inherit(App.Views.SignupForm, App.Views.View);
David.Utils.inherit(App.Views.LoginForm, App.Views.View);

App.Views.View.render = function () {
    this.anchor.innerText = '';
    var fragmentClone = this.fragment.cloneNode(true);

    if (this.applyListeners)
        this.applyListeners(fragmentClone);

    return this.anchor.appendChild(fragmentClone);
};

App.Views.Chat.formHTML = function () {
    this.chatDiv = document.createElement("div");
    this.chatDiv.id = "chat-div";
    this.chatDiv.classList.add('flex-column');

    this.chat = document.createElement("div");
    this.chat.id = "chat";

    this.form = document.createElement("form");
    this.form.id = "message-form"
    this.messageInput = document.createElement("input");
    this.messageInput.id = "message";
    this.messageInput.placeholder = "Your message goes here";
    this.messageInput.autocomplete = "off";
    this.messageInput.maxlength = "1023";

    this.submitButton = document.createElement('button');
    this.submitButton.innerText = "Send";
    this.submitButton.type = 'submit';
    this.submitButton.classList.add('btn-main');

    this.closeConMessage = document.createElement('h3');

    this.form.appendChild(this.messageInput);
    this.form.appendChild(this.submitButton);
    this.chatDiv.appendChild(this.chat);
    this.chatDiv.appendChild(this.form);
    this.chatDiv.appendChild(this.closeConMessage);
    this.fragment.appendChild(this.chatDiv);
}

App.Views.Chat.initConnection = function () {
    var self = this;

    return new Promise(function (resolve, reject) {
        try {
            self.ws = new WebSocket("wss://infinite-ocean-38389.herokuapp.com/chat");


            self.ws.onopen = function () {
                this.send(JSON.stringify({
                    body: "Hello people!!"
                }))
                self.closeConMessage.innerText = ""
                resolve()
            }

            self.ws.onclose = function (event) {
                self.closeConMessage.innerText = "Connection closed, please refresh to reconnect"
            }

            self.ws.onmessage = function (event) {
                if (self.chat.children.length > 20) {
                    self.eraseMessages(10);
                }

                var message = JSON.parse(event.data);

                var messageElement = new App.Models.Message(message.sendername, message.body, message.sendat, message.receivedat);

                var node = messageElement.asHTMLNode();


                if (self.chat.scrollTop + self.chat.clientHeight >= self.chat.scrollHeight) {
                    self.chat.appendChild(node);
                    self.chat.scrollTop = node.offsetHeight + node.offsetTop;
                } else {
                    self.chat.appendChild(node);
                }
            }

            self.ws.onerror = function (err) {
                self.chat.innerHTML += `<h3 class="message">An error occured!</h3>`;
                reject(err)
            }

        } catch (err) {
            reject(err);
        }


    })
}

App.Views.LoginForm.applyListeners = function (cloneFragment) {
    var err = cloneFragment.getElementById('err-login');
    var anchor = this.anchor;

    cloneFragment.lastChild.onsubmit = function (event) {
        event.preventDefault();

        err.innerHTML = `<h3 class="fade">Loading...</h3>`;

        var inputs = this.querySelectorAll('input');

        var info = {}

        inputs.forEach(function (input) {
            info[input.name] = input.value;
        });

        var xhr = new XMLHttpRequest();

        xhr.onload = function () {
            switch (this.status) {
                case 200:
                    var chatEvent = new Event('chat');

                    err.innerText = "Success!!";
                    anchor.dispatchEvent(chatEvent);

                    break;

                default:
                    err.innerText = this.responseText;

                    break;
            }
        }

        xhr.open('POST', '/login', true);
        xhr.withCredentials = true;

        xhr.send(JSON.stringify(info));
    }
};

App.Views.SignupForm.applyListeners = function (cloneFragment) {
    var err = cloneFragment.getElementById('err-signup');
    var anchor = this.anchor;

    cloneFragment.lastChild.onsubmit = function (event) {
        event.preventDefault();

        err.innerHTML = `<h3 class="fade">Loading...</h3>`;

        var inputs = this.querySelectorAll('input');

        var info = {}

        inputs.forEach(function (input) {
            info[input.name] = input.value;
        });

        xhr.onload = function () {
            switch (this.status) {
                case 200:
                    err.innerText = "Success!!";

                    break;

                default:
                    err.innerText = this.responseText;

                    break;
            }
        }

        xhr.open('POST', '/signup', true);
        xhr.withCredentials = true;

        xhr.send(JSON.stringify(info));
    }
};