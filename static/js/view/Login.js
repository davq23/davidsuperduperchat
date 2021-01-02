import View from "./View.js";

export default class LoginView extends View {
    constructor(anchor) {
        super(anchor)
        this.loginForm = document.createElement('form');
        this.loginForm.classList.add('form');

        this.title = document.createElement('h3');
        this.title.innerText = 'Login'

        this.err = document.createElement('h3');
        this.err.classList.add('form-error');
        this.err.id = 'err-login';

        const usernameLabel = document.createElement('label'); 
        usernameLabel.innerText = 'Username';
        usernameLabel.for = 'username';
        usernameLabel.classList.add('form-label');

        this.usernameInput = document.createElement('input');
        this.usernameInput.name = 'username';
        this.usernameInput.placeholder = 'Username';
        this.usernameInput.classList.add('form-input');
        
        const passwordLabel = document.createElement('label'); 
        passwordLabel.innerText = 'Password';
        passwordLabel.for = 'password';
        passwordLabel.classList.add('form-label');

        this.passwordInput  = document.createElement('input');
        this.passwordInput.type = 'password';
        this.passwordInput.name = 'password';
        this.passwordInput.placeholder = 'Password';
        this.passwordInput.classList.add('form-input');



        this.submitButton = document.createElement('button');
        this.submitButton.innerText = "Send"
        this.submitButton.type = 'submit';

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
    }

    applyListeners(cloneFragment) {
        const err = cloneFragment.getElementById('err-login');
        const anchor = this.anchor;

        cloneFragment.lastChild.onsubmit = async function(event) {
            event.preventDefault();

            const formData = new FormData(this);

            err.innerHTML = `<h3 class="fade">Loading...</h3>`;

            const info = {}
    
            formData.forEach((value, key) => {
                info[key] = value;
            });

                     
            const response = await fetch('/login', {
                'method': 'post',
                'credentials': 'include',
                'body': JSON.stringify(info)
            });

            const chatEvent = new Event('chat');
    
            if (response.status !== 200) {
                err.innerText = await response.text();
            } else {
                err.innerText = "Success!!";
                anchor.dispatchEvent(chatEvent);
            }
        }
    }

}