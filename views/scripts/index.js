const mainGrid = document.getElementById("main-grid")

/**
 * fetch query data from api
 * @param {String} query body of the query in the following form '"<body>"'
 * @returns {object} results of the query in JSON format
 */
async function getDataWithQuery(query){
    
    let raw = '{"query":'+query+', "schema":"items"}';
    let requestOptions = {
        method: 'POST',
        body: raw,
        redirect: 'follow'
      };
    let data;
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
        mainGrid.append(element)
    }
}


/**
 * fetch items from api and view them
 * @returns {any} none
 */
async function getItems(){
    const data = await getDataWithQuery('"select * from items"');
    createItemElements(data);
}

getItems();




// fetch("http://localhost:8080/query/", requestOptions)
// .then(response => response.text())
// .then(result => console.log(result))
// .catch(error => console.log('error', error));

