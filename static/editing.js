mainEntry = document.getElementById("mainEntry");
titleEntry = document.getElementById("titleEntry");

homebutton = document.getElementById("home")


const params = new URLSearchParams(window.location.search);
const id = params.get("id")

let lastSaveContent = "";
let lastSaveTitle = "";

async function loadcontent() {
    try {
        console.log(id)
        const res = await fetch("/api/load", {
                method: "Post",
                body: JSON.stringify({id: Number(id)}),
                headers: {"Content-Type": "application/json"}
            });
        if (!res.ok) {
            console.error("Response failed, status: " + res.status);
        }
        const result = await res.json();
        console.log(result);
        mainEntry.value = result.content;
        titleEntry.value = result.title
    }
    catch (error) {
        console.error("Error: " + error)
    }
}

async function savecontent(content, title) {
        try {
            await fetch("/api/update", {
                method: "Post",
                body: packagedocumentasjson(content, title),
                headers: {"Content-Type": "application/json"}
            });
            console.log(packagedocumentasjson(content, title))
            console.log("Saved")
        }
        catch (err) {
            console.log("Failed to save: " + err);
        }

        lastSaveContent = content;
        lastSaveTitle = title;
}
function takeMeHome() {
    window.location.replace("/")
}

function autosave() {
    let currentText = mainEntry.value;
    let currentTitle = titleEntry.value;

    if (currentText !== lastSaveContent || currentTitle !== lastSaveTitle) {
        savecontent(currentText, currentTitle);
    }
}

function packagedocumentasjson(content, title) {
    return JSON.stringify({id:Number(id), title: title ,content: content});
}




loadcontent();

homebutton.addEventListener('click', takeMeHome)

setInterval(autosave, 5000);