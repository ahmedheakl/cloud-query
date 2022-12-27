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
        redirect: 'follow'
    };
    
    try{
        let rawResponse = await fetch("http://localhost:8080/signin/", requestOptions);
        let res = await rawResponse.json();
        
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
    if(status){
        location.replace("http://127.0.0.1:5500/views/")
    }
})

