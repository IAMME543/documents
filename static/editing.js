mainEntry = document.getElementById("mainEntry");
titleEntry = document.getElementById("titleEntry");

savedState = document.getElementById("saveState");


homebutton = document.getElementById("home");
sharebutton = document.getElementById("share");
archivebutton = document.getElementById("archive");

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
        document.title = result.title + " | A Typing Site"
        autoResize(mainEntry)
    }
    catch (error) {
        console.error("Error: " + error)
    }
}

async function savecontent(content, title) {
        try {
            savedState.textContent = "Saving..."
            await fetch("/api/update", {
                method: "Post",
                body: packagedocumentasjson(content, title),
                headers: {"Content-Type": "application/json"}
            });
            console.log(packagedocumentasjson(content, title))
            console.log("Saved")
            savedState.textContent = "Saved ✓"
        }
        catch (err) {
            console.log("Failed to save: " + err);
            savedState.textContent = "Last Save Failed X" 
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


function autoResize(el) {
    el.style.height = "auto"
    el.style.height = el.scrollHeight + "px"
}

async function copyToClipboard(text, type) {
    try {
        await navigator.clipboard.writeText(text)
        alert(type + " URL Copied To Clipboard")
    }
catch {
    alert("Failed To Copy URL")
}
}
function share() {
    url = document.URL
    copyToClipboard(url, "Share")
}
function archive() {
    url = "https://atypingsite.masondoesthings.com/archive/" + id
    copyToClipboard(url, "Archive")
    
}


loadcontent();

homebutton.addEventListener('click', takeMeHome)
mainEntry.addEventListener('input', () => autoResize(mainEntry))

sharebutton.addEventListener('click', share)
archivebutton.addEventListener('click', archive)



setInterval(autosave, 5000);

