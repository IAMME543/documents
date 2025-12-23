mainEntry = document.getElementById("mainEntry");

let lastSaveContent = "";

async function loadcontent() {
    try {
        const res = await fetch("/api/load");
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

loadcontent();

setInterval(autosave, 5000);