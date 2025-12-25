documentList = document.getElementById("documentList");

async function ListDocuments() {
    try {
        const res = await fetch("/api/list");
        if (!res.ok) {
            console.error("list response failed" + res.error);
        }
        const result = await res.json()
        console.log(result)

    }
    catch (error) {
        console.error("Error: " + error)
    }
}

ListDocuments()