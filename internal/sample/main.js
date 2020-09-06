const xhr = new XMLHttpRequest();
const url = 'http://localhost:8080/api/v1/acted_in';
xhr.open("GET", url);
xhr.setRequestHeader("Content-Type", "application/json");
xhr.send();
var data

var container = document.getElementById('network');
var options = {
};

xhr.onload = (e) => {
    data = JSON.parse(xhr.responseText)
    var network = new vis.Network(container, data, options);
}
