const search = document.querySelector("#search");

search.addEventListener("input", (e) => {
    e.preventDefault();
    const data = e.data;

    // CREATE AJAX REQUEST
    const xhr = new XMLHttpRequest();

    // SEND AJAX REQUEST
    xhr.open("post", "/api/check");
    xhr.send(data);

    // RECEIVE AJAX REQUEST
    xhr.addEventListener("readystatechange", () => {
        if (xhr.readyState === 4 && xhr.status === 200){
            const response = xhr.responseText;
            console.log("Received: ", response);
        } else if(xhr.status === 500) {
            console.log("Something went wrong.");
        }
    })
})