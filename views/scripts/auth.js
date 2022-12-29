const BASEURL = "http://127.0.0.1:5500/"
const API = "https://faae-41-43-245-201.eu.ngrok.io"

if(localStorage.getItem("cookie") !== ''){
    location.replace(BASEURL)
}
const signUpButton = document.getElementById('signUp');
const signInButton = document.getElementById('signIn');
const container = document.getElementById('container');

// TODO: if logged in, go to index page directly
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
    
    try{
        let rawResponse = await fetch(`${API}/login/`, requestOptions);
        let res = await rawResponse.json();
        console.log(res);
        if(res.response == true){
            res.redirect = `${API}/cookie/${email}/`
            // for login and logout purposes. Generate your own cookie because the API endpoint doesn't return a cookie anymore.
            localStorage.setItem("cookie", res.cookie); 
            return true;
        }
    }catch(error){
        console.error("[FATAL] ", error);
        return false;
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
    
    try{
        let rawResponse = await fetch(`${API}/signup/`, requestOptions);
        let res = await rawResponse.json();
        console.log(res);
        if(res.response == true){
            res.redirect = `${API}/cookie/${email}/`
            // for login and logout purposes. Generate your own cookie because the API endpoint doesn't return a cookie anymore.
            localStorage.setItem("cookie", res.cookie); 
            return true;
        }
    }catch(error){
        console.error("[FATAL] ", error);
        return false;
    }
    return false; 
}

const loginForm = document.getElementById("login-form");
loginForm.addEventListener("submit", async function(e){
    e.preventDefault()
    const email = document.getElementById("login-email").value;
    const password = document.getElementById("login-password").value;

    const status = await login(email, password);

    // TODO: add loading symbol in login button
    // TODO: show error tag
    if(status){
        location.replace(BASEURL)
    }
}) 


const signUpForm = document.getElementById("signup-form");
signUpForm.addEventListener("submit", async (e) =>{
   e.preventDefault();
   const email =  document.getElementById("signup-email").value;
   const password = document.getElementById("signup-password").value;

   const status = await singUp(email, password);
   
   if(status){
        location.replace(BASEURL + "templates/auth.html")
   }
});

