mainEntry = document.getElementById("mainEntry");
titleEntry = document.getElementById("titleEntry");

const params = new URLSearchParams(window.location.search);
const id = params.get("id")

let lastSaveContent = "";
let lastSaveTitle = "";

async function loadcontent() {
    try {

        const res = await fetch("/api/load", {
                method: "Post",
                body: JSON.stringify({id}),
                headers: {"Content-Type": "application/json"}
            });
        if (!res.ok) {
            console.error("Response failed, status: " + res.status);
        }
        console.log(res)
        const result = await res.json();
        console.log(result);
        mainEntry.value = result.content;
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
            console.log("Saved")
        }
        catch (err) {
            console.log("Failed to save: " + err);
        }

        lastSaveContent = content;
        lastSaveTitle = title;
}

function autosave() {
    let currentText = mainEntry.value;
    let currentTitle = titleEntry.value;

    if (currentText !== lastSaveContent || currentTitle !== lastSaveTitle) {
        savecontent(currentText, currentTitle);
    }
}

function packagedocumentasjson(content, title) {
    return JSON.stringify({title: title ,content: content});
}




loadcontent();



setInterval(autosave, 5000);