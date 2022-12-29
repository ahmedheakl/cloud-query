BASEURL = "http://127.0.0.1:5500/views/"

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
        redirect: 'follow'
    };
    
    try{
        let rawResponse = await fetch("http://localhost:8080/signin/", requestOptions);
        let res = await rawResponse.json();
        console.log(res);
        if(res.response == true){
            localStorage.setItem("cookie", res.cookie)
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
        redirect: 'follow'
    };
    
    try{
        let rawResponse = await fetch("http://localhost:8080/signup/", requestOptions);
        let res = await rawResponse.json();
        console.log(res);
        if(res.response == true){
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

