// const BASEURL = "https://container-service-2.e41513gjaiic0.eu-central-1.cs.amazonlightsail.com/"
// var API = "https://container-service-1.e41513gjaiic0.eu-central-1.cs.amazonlightsail.com";
 const BASEURL = "http://localhost:5500/"
 const API = "http://localhost:8080"

if(localStorage.getItem("cookie") !== ''){
    location.replace(BASEURL)
}
const signUpButton = document.getElementById('signUp');
const signInButton = document.getElementById('signIn');
const container = document.getElementById('container');

signUpButton.addEventListener('click', () => {
	container.classList.add("right-panel-active");
});

signInButton.addEventListener('click', () => {
	container.classList.remove("right-panel-active");
});

/**
 * login with username and password
 * @returns {boolean} status flag for operation success
 */
async function login(email, password){
    let raw = `{\"email\":\"${email}\", \"password\":\"${password}\"}`;

    let  requestOptions = {
        method: 'POST',
        body: raw,
        credentials: 'include',
        redirect: 'follow'
    };
    
    let rawResponse = await fetch(`${API}/login/`, requestOptions);
    let res = await rawResponse.json();
    console.log(res)
    if(res.response == true){

        // for login and logout purposes, generate your own cookie
        // because the API endpoint doesn't return a cookie anymore.
        localStorage.setItem("cookie", "dummy"); 
        return true;
    }
    const errorLabel = document.getElementById("login-error");
    errorLabel.innerHTML = `* ${res.message}`;
    errorLabel.style.display = "flex";

    if(res.message.toLowerCase().includes("email")){
        const email =  document.getElementById("login-email");
        email.style.border = "1px red solid";
        email.classList.add("error");
    }else{
        const password = document.getElementById("login-password");
        password.style.border = "1px red solid";
        password.classList.add("error");
    }
    return false; 
    
}

async function singUp(email, password){
    let raw = `{\"email\":\"${email}\", \"password\":\"${password}\"}`;

    let  requestOptions = {
        method: 'POST',
        body: raw,
        credentials: 'include',
        redirect: 'follow'
    };
    
    
    let rawResponse = await fetch(`${API}/signup/`, requestOptions);
    let res = await rawResponse.json();
    if(res.response == true){
        return true;
    }

    const errorLabel = document.getElementById("signup-error");
    errorLabel.innerHTML = `* ${res.message}`;
    errorLabel.style.display = "block";

    if(res.message.toLowerCase().includes("email")){
        const email = document.getElementById("signup-email");
        email.style.border = "1px red solid";
        email.classList.add("error");
    }else{
        const password = document.getElementById("signup-password");
        password.style.border = "1px red solid";
        password.classList.add("error");
    }
    return false; 
}

const loginForm = document.getElementById("login-form");
loginForm.addEventListener("submit", async function(e){
    e.preventDefault()
    const email = document.getElementById("login-email").value;
    const password = document.getElementById("login-password").value;

    const status = await login(email, password);

    if(status){
        location.replace(`${API}/cookie/${email}/`);
        // location.replace(BASEURL)
    }
}) 


const signUpForm = document.getElementById("signup-form");
signUpForm.addEventListener("submit", async (e) =>{
   e.preventDefault();
   const email =  document.getElementById("signup-email").value;
   const password = document.getElementById("signup-password").value;

   const status = await singUp(email, password);
   
   if(status){
        location.replace(BASEURL + "templates/auth.html");
   }
});

