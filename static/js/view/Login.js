import View from "./View.js";

export default class LoginView extends View {
    constructor(anchor) {
        super(anchor)
        this.loginForm = document.createElement('form');

        this.title = document.createElement('h3');
        this.title.innerText = 'Login'

        this.err = document.createElement('h3');
        this.err.id = 'err-login';

        this.usernameInput = document.createElement('input');
        this.usernameInput.name = 'username';

        this.passwordInput  = document.createElement('input');
        this.passwordInput.type = 'password';
        this.passwordInput.name = 'password';


        this.submitButton = document.createElement('button');
        this.submitButton.innerText = "Send"
        this.submitButton.type = 'submit';

        this.loginForm.appendChild(this.title);
        this.loginForm.appendChild(this.err);
        this.loginForm.appendChild(this.usernameInput);
        this.loginForm.appendChild(this.passwordInput);
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