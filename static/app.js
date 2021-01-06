import SignupView from "./js/view/Signup.js";
import LoginView from "./js/view/Login.js";
import ChatView from "./js/view/Chat.js"
import LoadingView from "./js/view/Loading.js";

document.onreadystatechange = async function() {

    // App starts when page is completely loaded
    if (document.readyState == 'complete') {
        const AppDiv = document.getElementById("App");  

        // Nav buttons
        const loginButton = document.getElementById("login-button");
        const logoutButton = document.getElementById("logout-button");
        const signupButton = document.getElementById("signup-button");

        // Views
        const signupView = new SignupView(AppDiv);
        const loginView = new LoginView(AppDiv);
        const loadingView = new LoadingView(AppDiv);
        let chatView = null;

        AppDiv.addEventListener('change-nav', function(event) {
            event.preventDefault();
            if (event.chat) {
                loginButton.disabled = false;
                signupButton.disabled = false;
                logoutButton.disabled = true;

                loginButton.classList.remove('hidden');                
                signupButton.classList.remove('hidden');  
                logoutButton.classList.add('hidden');
                
            } else {
                loginButton.disabled = true;
                signupButton.disabled = true;
                logoutButton.disabled = false;

                loginButton.classList.add('hidden');                
                signupButton.classList.add('hidden');
                logoutButton.classList.remove('hidden');      
            }      
        });  

        // App events
        AppDiv.addEventListener('loading-', function(event) {
            event.preventDefault();
            loadingView.render();
        })

        AppDiv.addEventListener('chat', async function(event)  {
            loadingView.render();

            let evt = null

            try {
                // Try to open websocket
                chatView = new ChatView(AppDiv);

                evt = new CustomEvent("NavChange", {chat: true});

                await chatView.initConnection()
    
                 // If successful, render chat and hide buttons
                chatView.render();

               

            } catch(err) {
                console.log(err);

                evt = new CustomEvent("NavChange", {chat: false});

                // Render login and load buttons
                loginView.render();
            }

             // onclick listeners
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

                await chatView.logout();

                loginView.render();

                AppDiv.dispatchEvent(new CustomEvent("NavChange", {chat: false}));
            }

            AppDiv.dispatchEvent(evt);
        });

        const event = new Event('chat');
        AppDiv.dispatchEvent(event);
    }
   
    
}

