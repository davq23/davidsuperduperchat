import SignupView from "./js/view/Signup.js";
import LoginView from "./js/view/Login.js";
import ChatView from "./js/view/Chat.js"
import LoadingView from "./js/view/Loading.js";

document.onreadystatechange = async function() {
    if (document.readyState == 'complete') {
        const AppDiv = document.getElementById("App");  
        const loginButton = document.getElementById("login-button");
        const logoutButton = document.getElementById("logout-button");
        const signupButton = document.getElementById("signup-button");

        const signupView = new SignupView(AppDiv);
        const loginView = new LoginView(AppDiv);
        const loadingView = new LoadingView(AppDiv);

        let chatView = null;

        AppDiv.addEventListener('loading-', function(event) {
            event.preventDefault();
            loadingView.render();
        })

        AppDiv.addEventListener('chat', async function(event)  {
            loadingView.render();

            try {
                
                chatView = new ChatView(AppDiv);

                await chatView.initConnection()
    
                chatView.render();

                loginButton.disabled = true;
                signupButton.disabled = true;
                logoutButton.disabled = false;

                loginButton.classList.add('hidden');                
                signupButton.classList.add('hidden');
                logoutButton.classList.remove('hidden');       


            } catch(err) {
                console.log(err);
                loginView.render();

                loginButton.classList.remove('hidden');                
                signupButton.classList.remove('hidden');  
                logoutButton.classList.add('hidden');


                loginButton.disabled = false;
                signupButton.disabled = false;
                logoutButton.disabled = true;

               
            }

            loginButton.onclick = function() {
                loadingView.render();

                loginView.render();
            }

            signupButton.onclick = function() {
                loadingView.render();

                signupView.render();
            }

            logoutButton.onclick = async function() {
                loadingView.render();

                const response = await fetch('/logout', {
                    method: 'post', 
                    credentials: 'include'
                });
                
                if (response.status === 200) {
                    await chatView.ws.close();

                    const event = new Event('chat');
                    AppDiv.dispatchEvent(event);
                }  
            }
       
        });

        const event = new Event('chat');
        AppDiv.dispatchEvent(event);
    }
   
    
}

