mainEntry = document.getElementById("mainEntry");
titleEntry = document.getElementById("titleEntry");

let lastSaveContent = "";
let lastSaveTitle = "";
let currentDocID = "";

async function loadcontent(id) {
    try {
        const res = await fetch("/api/load", {
                method: "Post",
                body: json.stringify({id}),
                headers: {"Content-Type": "application/json"}
            });
        if (!res.ok) {
            console.error("Response failed, status: " + res.status);
        }
        const result = await res.json();
        console.log(result);
        mainEntry.value = result.content;
    }
    catch (error) {
        console.error("Error: " + error)
    }
}

async function savecontent(content, title, id) {
        try {
            await fetch("/api/save", {
                method: "Post",
                body: packagedocumentasjson(content, title, id),
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
        if (currentDocID == "") {
            currentDocID = generatedocumentid();
        }
        savecontent(currentText, currentTitle, currentDocID);
    }
}

function packagedocumentasjson(content, title, id) {
    return JSON.stringify({id: id, title: title ,content: content});
}

function generatedocumentid(length = 8) {
    const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
    let id = '';
    for (let i = 0; i < length; i++) {
        id += chars.charAt(Math.floor(Math.random() * chars.length));
    }
    return id;
}

if (currentDocID != "") {
    loadcontent(currentDocID);
}


setInterval(autosave, 5000);