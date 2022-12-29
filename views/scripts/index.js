const BASEURL = "http://127.0.0.1:5500/views/"

const navAuth = document.getElementById("nav-auth");
let element; 
if(localStorage.getItem("cookie") !== ''){
    navAuth.innerHTML = `
        <a href="/views/templates/auth.html" id="logout-button">
            <i class="glyphicon glyphicon-log-out" style="padding-right: 5px"></i>
            Logout
        </a>
    `
    const logoutButton = document.getElementById("logout-button")
    logoutButton.addEventListener("click", (e) => {
        localStorage.setItem("cookie", "");
        location.replace(BASEURL+"templates/auth.html");
    }) 
    
    
}else{
    navAuth.innerHTML = `
    <a href="/views/templates/auth.html" id="login-button">
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
let cartData = {};


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
        body: raw,
        redirect: 'follow'
    };
    try{
        let dataJson = await fetch("http://localhost:8080/query/", requestOptions)
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
 * @returns {any}
 */
async function addItemToCart(itemid){
    let raw = `{"items": "${itemid}", "quantity": 1}`
    let requestOptions = {
        method: 'POST',
        body: raw,
        redirect: 'follow'
    };
    
    try{
        let rawResponse = fetch("http://localhost:8080/add/", requestOptions);
        // let res = await rawResponse.text();
    }catch(error){
        console.log(error);
    } 
}

/**
 * update cart UI with items found in `cartItems`
 * @param {object} cartItems list of items 
 * @returns {any}
 */
function updateCart(cartItems){
    let innerHTML = "";
    if(Object.keys(cartItems).length == 0){
        innerHTML = "There is nothing your cart."
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
        for(let item of Object.values(cartItems)){
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
        <button class="checkout">Checkout</button>
        `
    }
    const cart = document.getElementById("cart");
    cart.innerHTML = innerHTML;
    cart.style.display = "flex";
}

function setUpCartEventListeners(){
    let elements = document.querySelectorAll(".grid-item")
    for(let element of elements){
        const elementId = element.dataset.id;

        // event listener for add buttons
        const addButton = element.getElementsByClassName("grid-item-button button-add")[0]
        addButton.addEventListener("click", async function(e){
            const itemsQuery = await getDataWithQuery(`select * from items where id=${elementId}`, "items");
            const item = itemsQuery[0];
            if(item.id in cartData){
                cartData[item.id].quantity += 1;
            }else{
                cartData[item.id] = {
                    "id": item.id,
                    "price": item.price,
                    "name": item.name,
                    "quantity": 1,
                    "image": item.image,
                };
            }
            
            updateCart(cartData);
        })

        // event listener for remove button
        const removeButton = element.getElementsByClassName("grid-item-button button-delete")[0]
        removeButton.addEventListener("click", async function(e){
            const itemId = parseInt(elementId);

            if(itemId in cartData){
                if(cartData[elementId].quantity == 1){
                    delete cartData[elementId];
                }else{
                    cartData[elementId].quantity -= 1; 
                }
            }
            updateCart(cartData);
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
    