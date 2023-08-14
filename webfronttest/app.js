const baseURL = 'https://4ade-180-183-137-112.ap.ngrok.io/course'
fetch(baseURL)

.then(response => response.json())
.then(response => {
    appendData(response)
    // console.log(response)

})
.catch(function (err){
    console.log(`error: ` + err)
})

function appendData(data){
    var mainContainer = document.getElementById("myData");
    for (let i = 0; i < data.length; i++) {
        let div = document.createElement("div");
        div.innerHTML = `CourseID: ` + data[i].ID + ` `+ data[i].Name + ` ` + data[i].Price + ` ` + data[i].Instructor + ` `;
        mainContainer.appendChild(div);        
    }
    // document.querySelector("pre").innerHTML = JSON.stringify(data, null, 2);
}