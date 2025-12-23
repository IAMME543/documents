mainEntry = document.getElementById("mainEntry");

let lastSaveContent = "";

async function loadcontent() {
    try {
        const res = await fetch("/api/save");
        if (!res.ok) {
            throw new Error("Response failed, status: " + res.status);
        }
        const result = await res.json;
        mainEntry.value = result;
    }
    catch (error) {
        throw new Error("Error: " + error)
    }
}

async function savecontent(content) {
        try {
            await fetch("/api/save", {
                method: "Post",
                body: JSON.stringify({content}),
                headers: {"Content-Type": "application/json"}
            });
            lastSaveContent = content;
            console.log("Saved")
        }
        catch (err) {
            console.log("Failed to save: " + err);
        }
}

function autosave() {
    let currentText = mainEntry.value;
    if (currentText != lastSaveContent) {
        savecontent(currentText);
    }
}

setInterval(autosave, 5000);