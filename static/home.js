documentList = document.getElementById("documentList");
createButton = document.getElementById("create")

async function GetDocumentsIndex() {
    try {
        const res = await fetch("/api/list");
        if (!res.ok) {
            console.error("list response failed" + res.error);
        }
        const result = await res.json()
        console.log(result)
        return result

    }
    catch (error) {
        console.error("Error: " + error)
    }
}

function InsertToList(docindex) {
    console.log(docindex)

    docindex.forEach((doc) => {
        let listitem = document.createElement("li")
        let anchor = document.createElement("a")

        anchor.target = "_blank";
        anchor.rel = "noopener noreferrer";
        anchor.href = "/editing/?id=" + doc.id
        if (doc.title != "") {
            anchor.textContent = doc.title
        } else {
            anchor.textContent = "Untitled Document " + doc.id
        }


        listitem.appendChild(anchor)
        documentList.appendChild(listitem)

    });
}

async function createDocument() {
    try {
        const res = await fetch("/api/create");
        if (!res.ok) {
            console.error("create fetch failed" + res.error)
        }
        const result = await res.text();
        console.log(result)
        window.location.replace("/editing?id=" + result)

    }
        catch (error) {
        console.error("Error: " + error)
    }
}

GetDocumentsIndex().then(docindex => {
    InsertToList(docindex);
});

createButton.addEventListener('click', createDocument)
