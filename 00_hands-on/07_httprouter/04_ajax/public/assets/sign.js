const forms = document.querySelectorAll("form");
let lastPath;

function checkName(name) {
    return (
        (name === "") ? "First name field is empty" :
        (name < 3) ? "First name is too short" : null
    )
}

function checkEmail(email) {
    return (
        (email === "") ? "Email field is empty" :
        (email.length < 5 ) ? "Email is too short" :
        (/^\w+([\.-]?\w+)*@\w+([\.-]?\w+)*(\.\w{2,3})+$/.test(email)) ?  null : "Email is invalid"
    )
}

function checkPassword(password) {
    // Password should contain at least one digit
    // Password should contain at least one lowercase
    // Password should contain at least one uppercast
    // Password should contain at least 8 from the characters above
    passRex = /^(?=.*\d)(?=.*[a-z])(?=.*[A-Z])[0-9a-zA-Z]{8,}$/
    return (
        (password === "") ? "Password field is empty" :
        (password.length < 6) ? "Password is too shorty" :
        (passRex.test(password)) ? null : "Password is too weak"
    )
}

function showMessages(path, message) {
    const err = document.createElement("p");
    err.className = "warning-error";
    
    // Check if element exists in the DOM
    find = document.querySelector(".warning-error")
    if ((path === lastPath) && (find)) {
        console.log("Path equals!!")
        find.innerText = message;
        return
    }

    const text = document.createTextNode(message);
    err.appendChild(text);
    path.insertAdjacentElement('afterbegin', err);
    lastPath = path 
}

forms.forEach((form) => {
    const location = form.action.split("/", 3).join("/")+"/";
    const action = form.action.split("/", 4)[3]+"/";
    const url = location + action;
    form.addEventListener("submit", (e) => {
        if (action === "login/") {
            e.preventDefault();

            // Set form values
            email = form["email"].value;
            password = form["password"].value;

            // Validate form
            if (checkEmail(email) !== null) { showMessages(form, checkEmail(email)); return;}
            if (checkPassword(password) !== null) { showMessages(form, checkPassword(password)); return;}

            // Make the request 
            const xhr = new XMLHttpRequest();
            xhr.open("POST", url)
            xhr.send(JSON.stringify({email, password}))

            // Get the response
            xhr.addEventListener("readystatechange", () => {
                if (xhr.readyState === 4) {
                    if (xhr.status === 200) {
                        window.location.href = "/profile"
                    } else {
                        console.log("Something went wrong on: ", xhr.responseText);
                        showMessages(form, xhr.responseText)
                    }
                }
            })

        } else if (action === "create-account/") {
            e.preventDefault()

            // Set form values
            first = form["first"].value
            email = form["email"].value;
            password = form["password"].value;
            password2 = form["password2"].value

            // Validate form
            if (checkName(first) !== null) { showMessages(form, checkName(first)); return}
            if (checkEmail(email) !== null) { showMessages(form, checkEmail(email)); return}
            if (checkPassword(password) !== null) { showMessages(form, checkPassword(password)); return}
            else if(password !== password2 ) { showMessages(form, "Passwords must be equals."); return}

            // Set the request
            xhr = new XMLHttpRequest()
            xhr.open("POST", url)
            xhr.send(JSON.stringify({first, email, password}))

            // Get the response
            xhr.addEventListener("readystatechange", () => {
                if (xhr.readyState === 4) {
                    if (xhr.status === 200) { 
                        signin.style.display = "block";
                        signup.style.display = "none";
                    } else {
                        console.log("Something went wrong: ", xhr.responseText)
                        showMessages(form, xhr.responseText)
                    }
                }
            })
        } else {
            e.preventDefault()
            return
        }
    })
})

const signin = document.getElementById("signin")
const signup = document.getElementById("signup")
const signinBtn = document.getElementById("signinBtn");
const sigupBtn = document.getElementById("signupBtn");

signinBtn.addEventListener("click", () => {
    signin.style.display = "block";
    signup.style.display = "none";
})
signupBtn.addEventListener("click", () => {
    signin.style.display = "none";
    signup.style.display = "block";
})