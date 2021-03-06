import View from "./View.js";


export default class SignupView extends View {
    constructor(anchor) {
        super(anchor)
        this.signupForm = document.createElement('form');
        this.signupForm.classList.add('form');

        this.title = document.createElement('h3');
        this.title.innerText = 'Signup'

        this.err = document.createElement('h3');
        this.err.classList.add('form-error');
        this.err.id = 'err-signup';

        const usernameLabel = document.createElement('label'); 
        usernameLabel.innerText = 'Username';
        usernameLabel.for = 'username';
        usernameLabel.classList.add('form-label');

        const passwordLabel = document.createElement('label'); 
        passwordLabel.innerText = 'Password';
        passwordLabel.for = 'password';
        passwordLabel.classList.add('form-label');

        this.usernameInput = document.createElement('input');
        this.usernameInput.name = 'username';
        this.usernameInput.placeholder = 'Username';
        this.usernameInput.classList.add('form-input');

        this.passwordInput  = document.createElement('input');
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
        this.signupForm.appendChild(this.usernameInput);

        this.signupForm.appendChild(passwordLabel);
        this.signupForm.appendChild(this.passwordInput);
        this.signupForm.appendChild(this.submitButton);

        this.fragment.appendChild(this.signupForm)
    }

    applyListeners(cloneFragment) {
        const err = cloneFragment.getElementById('err-signup');
        const anchor = this.anchor;

        const submitButton = this.submitButton;

        cloneFragment.lastChild.onsubmit = async function(event) {
            event.preventDefault();

            err.innerHTML = `<h3 class="fade">Loading...</h3>`;

            submitButton.disabled = true;
    
            const formData = new FormData(this);
    
            const info = {}
    
            formData.forEach((value, key) => {
                info[key] = value;
            })
            
            const response = await fetch('/signup', {
                'method': 'post',
                'credentials': 'include',
                'body': JSON.stringify(info)
            });
    
            if (response.status !== 200) {
                err.innerText = await response.text();
            } else {
                err.innerText = "Success!!";
            }

            submitButton.disabled = false;
        }
    }

}