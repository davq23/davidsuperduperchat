import View from "./View.js";


export default class SignupView extends View {
    constructor(anchor) {
        super(anchor)
        this.signupForm = document.createElement('form');

        this.title = document.createElement('h3');
        this.title.innerText = 'Signup'

        this.err = document.createElement('h3');
        this.err.id = 'err-signup';

        this.usernameInput = document.createElement('input');
        this.usernameInput.name = 'username';

        this.passwordInput  = document.createElement('input');
        this.passwordInput.type = 'password';
        this.passwordInput.name = 'password';


        this.submitButton = document.createElement('button');
        this.submitButton.innerText = "Send"
        this.submitButton.type = 'submit';

        this.signupForm.appendChild(this.title);
        this.signupForm.appendChild(this.err);
        this.signupForm.appendChild(this.usernameInput);
        this.signupForm.appendChild(this.passwordInput);
        this.signupForm.appendChild(this.submitButton);

        this.fragment.appendChild(this.signupForm)
    }

    applyListeners(cloneFragment) {
        const err = cloneFragment.getElementById('err-signup');
        const anchor = this.anchor;

        cloneFragment.lastChild.onsubmit = async function(event) {
            event.preventDefault();

            err.innerHTML = `<h3 class="fade">Loading...</h3>`;
    
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
        }
    }

}