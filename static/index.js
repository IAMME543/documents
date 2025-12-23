mainEntry = document.getElementById("mainEntry");

let lastSaveContent = "";

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