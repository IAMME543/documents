mainEntry = document.getElementById("mainEntry");
titleEntry = document.getElementById("titleEntry");

let lastSaveContent = "";
let lastSaveTitle = "";

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

async function savecontent(content, title) {
        try {
            await fetch("/api/save", {
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




loadcontent(currentDocID);



setInterval(autosave, 5000);