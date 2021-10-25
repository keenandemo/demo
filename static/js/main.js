
var APP = {
    xhr: null
};


APP.Init = function() {
    document.addEventListener("click", APP.onClick)
}

APP.onClick = function(t) {
    let targ = t.target;

    switch (targ.getAttribute("data-cmd")) {
        case "clear-filter":
            console.log("filter cleared!");
            break;
        case "filter-ingest":
            console.log("filter mode: ingest");
            break;
        case "filter-movies":
            console.log("filter mode: movies");
            break;
        case "filter-shows":
            console.log("filter mode: shows");
            break;
    }
}

APP.xhr = function(method, url, data) {
    let req = new XMLHttpRequest;
    req.open(method, url, true);
    req.send(data);
}


document.addEventListener("DOMContentLoaded", APP.Init)