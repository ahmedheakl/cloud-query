// const BASEURL = "https://container-service-2.e41513gjaiic0.eu-central-1.cs.amazonlightsail.com/"
// var API = "https://container-service-1.e41513gjaiic0.eu-central-1.cs.amazonlightsail.com";
const BASEURL = "http://localhost:5500/"
const API = "http://localhost:8080"

const navAuth = document.getElementById("nav-auth");
let element; 


if(localStorage.getItem("cookie") !== '' || localStorage.getItem("cookie") === null){
    navAuth.innerHTML = `
        <a href="${API+"/cookie/remove/"}" id="logout-button">
            <i class="glyphicon glyphicon-log-out" style="padding-right: 5px"></i>
            Logout
        </a>
    `
    const logoutButton = document.getElementById("logout-button")
    logoutButton.addEventListener("click", async function(e){
        localStorage.setItem("cookie", "");
        await fetch(`${API}/logout/`, {method: "GET", redirect: "follow"});
    }) 
    
    
}else{
    navAuth.innerHTML = `
    <a href="/templates/auth.html" id="login-button">
        <i class="glyphicon glyphicon-log-in" style="padding-right: 5px"></i>
        Login
    </a>
    `
    const loginButton = document.getElementById("login-button")
    loginButton.addEventListener("click", (e) => {
        location.replace(BASEURL+"templates/auth.html")
    }) 
}

const mainGrid = document.getElementById("main-grid")
const cartButton = document.getElementById("shopping-cart-button");


cartButton.addEventListener("click", (e) => {
    const cart = document.getElementById("cart");
    if(cart.style.display == "flex"){
        cart.style.display = "none";
    }else{
        cart.style.display = "flex";
    }
    
})

let data = []

/**
 * fetch query data from api
 * @param {String} query body of the query in the following form '"<body>"'
 * @returns {object} results of the query in JSON format
 */
async function getDataWithQuery(query, schema){
    
    let raw = `{"query":"${query}", "schema":"${schema}"}`;
    let requestOptions = {
        method: 'POST',
        credentials: 'include',
        body: raw,
        redirect: 'follow',
    };
    try{
        let dataJson = await fetch(`${API}/query/`, requestOptions)
        data = await dataJson.json();
    }catch(error){
        console.log(error);
    }

    return data;
}


/**
 * create/view elements for each item 
 * @param {object} data data object in JSON format for incoming items data
 * @returns {any}  none
 */
function createItemElements(data){
    mainGrid.innerHTML = "";
    for(let i = 0; i < data.length; i++){
        let dataItem = data[i];
        const item = `
            <div class="grid-item-container">
                <div class="grid-item-name">
                    <h3>
                        ${dataItem.name}
                    </h3>
                </div>
                <div class="grid-item-image-container">
                    <img src="assets/images/${dataItem.image}" />
                </div>
                <div class="grid-item-description-price">
                    <div class="grid-item-description">
                        <p>
                            ${dataItem.description}
                        </p>
                    </div>
                </div>
                <div class="grid-item-price">
                    <p> ${dataItem.price}$
                    </p>
                </div>
            </div>
            <div class="grid-item-buttons">
                <button class="grid-item-button button-add" type="button">
                    <i class="material-icons">exposure_plus_1</i>
                </button>
                <button class="grid-item-button button-delete" type="button"><i
                        class="material-icons">exposure_neg_1</i></button>
            </div>`
        const element = document.createElement("div")
        element.innerHTML = item
        element.className = "grid-item"
        element.setAttribute("data-id", dataItem.id);
        mainGrid.append(element)
    }
}

/**
 * Add item with itemid input to the purchases table
 * @param {number} itemid id of the cart item to be added
 * @returns {boolean} status flag for valid/invalid addition
 */
async function addItemToCart(itemid){
    let raw = `{"item": ${itemid}, "quantity": 1}`
    let requestOptions = {
        credentials: 'include',
        method: 'POST',
        body: raw,
        redirect: 'follow'
    };
    
    let rawResponse = await fetch(`${API}/add/`, requestOptions);
    let res = await rawResponse.json();
    if(res.response == false){
        alert(res.message)
        return false;
    }
    return true;
}

/**
 * Remove item with itemid input to the purchases table
 * @param {number} itemid id of the cart item to be removed
 * @returns {boolean} status flag for valid/invalid removal
 */
async function removeItemToCart(itemid){
    let raw = `{"item": ${itemid}, "quantity": 1}`
    let requestOptions = {
        credentials: 'include',
        method: 'POST',
        body: raw,
        redirect: 'follow'
    };
    
    let rawResponse = await fetch(`${API}/remove/`, requestOptions);
    let res = await rawResponse.json();
    if(res.response == false){
        alert(res.message)
        return false;
    }
    return true;
}


/**
 * update cart UI with items found in `cartItems`
 * @param {object} cartItems list of items 
 * @returns {any}
 */
async function updateCart(){
    let innerHTML = "";
    var data = await fetch(`${API}/items/`, {"method": "GET",'credentials':'include'});
    var items = await data.json();
    if(items.length === 0){
        innerHTML = "There is nothing your cart."
        cart.innerHTML = innerHTML;
    }else{
        let totalPrice = 0;
        innerHTML = `
        <table>
            <tr>
                <th></th>
                <th>Name</th>
                <th>Price</th>
                <th>QTY</th>
                <th>STotal</th>
            </tr>
        `
        for(var i = 0; i < items.length; i++){
            var item = items[i];
            let subtotalPrice = item.quantity * item.price;
            totalPrice += subtotalPrice;
            subtotalPrice = subtotalPrice.toFixed(2);
            innerHTML += `
                <tr>
                    <td>
                        <img src="assets/images/${item.image}">
                    </td>
                    <td>${item.name}</td>
                    <td>${item.price}$</td>
                    <td>${item.quantity}</td>
                    <td>${subtotalPrice}$</td>
                </tr>
            `
        }
        totalPrice = totalPrice.toFixed(2);
        innerHTML += `
            <tr>
                <td></td>
                <td></td>
                <td></td>
                <td>Total:</td>
                <td>${totalPrice}$</td>
            </tr>
        </table>
        <button class="checkout" id="checkout-button">Checkout</button>
        `

    const cart = document.getElementById("cart");
    cart.innerHTML = innerHTML;
    cart.style.display = "flex";

    // setup checkout event listener
    const checkoutButton = document.getElementById("checkout-button");
    checkoutButton.addEventListener("click", async function(e){
        let requestOptions = {
            credentials: 'include',
            method: 'GET',
            redirect: 'follow'
        };
        let rawResponse = await fetch(`${API}/checkout/`, requestOptions);
        let res = await rawResponse.json();
        console.log(res);
        if(res.response == true){
            alert(res.message);
            updateCart();
        }
    });
    }
}

function setUpCartEventListeners(){
    let elements = document.querySelectorAll(".grid-item")
    for(let element of elements){
        const elementId = element.dataset.id;

        // event listener for add buttons
        const addButton = element.getElementsByClassName("grid-item-button button-add")[0]
        addButton.addEventListener("click", async function(e){
            let status = await addItemToCart(elementId);
            if(status == false) return; 
            // const itemsQuery = await getDataWithQuery(`select * from items where id=${elementId}`, "items");
            // const item = itemsQuery[0];
            // if(item.id in cartData){
            //     cartData[item.id].quantity += 1;
            // }else{
            //     cartData[item.id] = {
            //         "id": item.id,
            //         "price": item.price,
            //         "name": item.name,
            //         "quantity": 1,
            //         "image": item.image,
            //     };
            // }
            
            updateCart();
        })

        // event listener for remove button
        const removeButton = element.getElementsByClassName("grid-item-button button-delete")[0]
        removeButton.addEventListener("click", async function(e){
            let status = await removeItemToCart(elementId);
            if(status == false) return; 
            // const itemId = parseInt(elementId);
            // if(itemId in cartData){
            //     if(cartData[elementId].quantity == 1){
            //         delete cartData[elementId];
            //     }else{
            //         cartData[elementId].quantity -= 1; 
            //     }
            // }
            updateCart();
        });

        
    }    
}


/**
 * fetch items from api, set event handlers, and view items
 * @returns {any} none
 */
async function main(){
    const brands = new Set();
    const data = await getDataWithQuery('select * from items', "items");
    createItemElements(data);
    brands.add("All");
    await updateCart();

    for(let dataItem of data){
        brands.add(dataItem.brand);
    }
    const filterMenu = document.getElementById("filter-menu");
    for(let brand of brands){
        const brandFilter = document.createElement("a");
        brandFilter.innerHTML = brand;
        brandFilter.addEventListener("click", function(e){
            let filteredData = [];
            for(let dataItem of data){
                if(dataItem.brand == brand || brand === "All"){
                    filteredData.push(dataItem);
                }
            }
            createItemElements(filteredData);
            setUpCartEventListeners();            
        });
        filterMenu.append(brandFilter);
    }

    

    // search handler
    const searchInput = document.getElementById("search-input");
    searchInput.addEventListener("input", function(e){
        const searchParam = searchInput.value;
        let filteredData = [];
        for(let dataItem of data){
            const dataName = dataItem.name.toLowerCase();
            if(dataName.match(`.*${searchParam}.*`) == null){
                continue;
            }
            filteredData.push(dataItem);
        }
        createItemElements(filteredData);
        setUpCartEventListeners();
    })

    setUpCartEventListeners();
}

main();
    
