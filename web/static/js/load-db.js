function loadDatabase() {
    var e = document.getElementById("neighbors");
    var xhr = new XMLHttpRequest();
    xhr.onreadystatechange = function() {
        if (xhr.readyState == 4 && xhr.status == 200) {
            e.outerHTML = xhr.responseText;
        }
    }
    xhr.open("GET", "/nmap/more", true);
    try {
        xhr.send();
    } catch (err) {
        /* handle error */
        console.error("xhr.send()");
    }
}
