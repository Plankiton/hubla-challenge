let table = document.getElementById("sales"),
  countSpan = document.getElementById("count"),
  sumSpan = document.getElementById("sum");

const limit = 10;
let offset = 0, pageCount = 0, currPage = 0;
let salesRes = null;
th = table.rows[0].outerHTML;

function makeTr(sale) {
  sale.date = new Date(sale.date);
  return `
<tr>
  <td>${sale.type}</td>
  <td>${sale.date.toLocaleString("pt-BR")}</td>
  <td>${sale.product}</td>
  <td>RS$ ${sale.value.toFixed(2).replace(".", ",")}</td>
  <td>${sale.seller}</td>
</tr>
      `;
}

function updateSaleMeta() {
  let xhr = new XMLHttpRequest();
  xhr.open('GET', "/api/sales/meta");
  xhr.onreadystatechange = function(e) {
    if (xhr.readyState == 4) {
      // Everything is good!
      r = JSON.parse(xhr.response);

      countSpan.innerHTML = r.meta.count;
      sumSpan.innerHTML = `R$ ${r.meta.total.toFixed(2)}`;

      console.log(r);
    }
  };
  xhr.send();
}

function updateSaleList(onResponse) {
  let xhr = new XMLHttpRequest();
  xhr.open('GET', "/api/sales?limit=" + limit + "&offset=" + offset, true);
  xhr.onreadystatechange = function(e) {
    if (xhr.readyState == 4) {
      // Everything is good!
      salesRes = JSON.parse(xhr.response);
      console.log("get sales res: ", salesRes);

      pageCount = salesRes.meta.total / limit;
      onResponse(salesRes);
    }
  };

  xhr.send();
  updateSaleMeta();
}

function pageButtons(pCount, cur) {
  var prevDis = (cur == 1) ? "disabled" : "",
    nextDis = (cur == pCount) ? "disabled" : "",
    buttons = "<input type='button' value='&lt;&lt; Prev' onclick='updateSaleList(s => sort(s, " + (cur - 1) + "))' " + prevDis + ">";
  for (i = 1; i <= pCount; i++)
    buttons += "<input type='button' id='id" + i + "'value='" + i + "' onclick='updateSaleList(s => sort(s, " + i + "))'>";
  buttons += "<input type='button' value='Next &gt;&gt;' onclick='updateSaleList(s => sort(s, " + (cur + 1) + "))' " + nextDis + ">";
  return buttons;
}

function sort(salesRes, page) {
  offset = (page * limit) - 1

  let rows = th;
  for (let i = 0; i < salesRes.meta.limit; i++) {
    rows += makeTr(salesRes.data[i]);
  }

  table.innerHTML = rows;

  document.getElementById("buttons").innerHTML = pageButtons(pageCount, page);
  document.getElementById("id" + page).setAttribute("class", "active");

  currPage = page;
}

function fileSelectHandler(e) {
  // Fetch FileList object
  let files = e.target.files || e.dataTransfer.files;

  // Cancel event and hover styling
  fileDragHover(e);

  // Process all File objects
  for (let i = 0, f; f = files[i]; i++) {
    parseFile(f);
    uploadFile(f);
  }
}
// Output
//
function output(msg) {
  // Response
  let m = document.getElementById('messages');
  m.innerHTML = msg;
}

function parseFile(file) {

  console.log(file.name);
  output(
    '<strong>' + encodeURI(file.name) + '</strong>'
  );

  let fileType = file.type;
  console.log("file type:", fileType);

  document.getElementById('start').classList.add("hidden");
  document.getElementById('response').classList.remove("hidden");
}

function setProgressMaxValue(e) {
  let pBar = document.getElementById('file-progress');

  if (e.lengthComputable) {
    pBar.max = e.total;
  }
}

function updateFileProgress(e) {
  let pBar = document.getElementById('file-progress');

  if (e.lengthComputable) {
    pBar.value = e.loaded;
  }
}

function uploadFile(file) {
  let xhr = new XMLHttpRequest(),
    fileInput = document.getElementById('class-roster-file'),
    pBar = document.getElementById('file-progress'),
    fileSizeLimit = 1024; // In MB

  xhr.onreadystatechange = function(e) {
    if (xhr.readyState == 4) {
      output("File already sent! Please refresh the page to send new transactions");
      pBar.style.display = "none";

      updateSaleList(s => sort(s, currPage))
    }
  };

  if (xhr.upload) {
    // Check if file is less than x MB
    if (file.size <= fileSizeLimit * 1024 * 1024) {
      pBar.style.display = 'inline';
      xhr.upload.addEventListener('loadstart', setProgressMaxValue, false);
      xhr.upload.addEventListener('progress', updateFileProgress, false);
      xhr.open('POST', "/api/sales?filename=sales", true);

      const formData = new FormData();
      formData.append("sales", file);
      xhr.send(formData);
    } else {
      output('Please upload a smaller file (< ' + fileSizeLimit + ' MB).');
    }
  }
}

function fileDragHover(e) {
  let fileDrag = document.getElementById('file-drag');

  e.stopPropagation();
  e.preventDefault();

  fileDrag.className = (e.type === 'dragover' ? 'hover' : 'modal-body file-upload');
}


function Init() {
  let fileSelect = document.getElementById('file-upload'),
    fileDrag = document.getElementById('file-drag');
  fileSelect.addEventListener('change', fileSelectHandler, false);

  let xhr = new XMLHttpRequest();
  if (xhr.upload) {
    // File Drop
    fileDrag.addEventListener('dragover', fileDragHover, false);
    fileDrag.addEventListener('dragleave', fileDragHover, false);
    fileDrag.addEventListener('drop', fileSelectHandler, false);
  }

  console.log("Upload Initialized");
}

if (window.File && window.FileList && window.FileReader) {
  Init();
} else {
  document.getElementById('file-drag').style.display = 'none';
}

updateSaleList((sales) => {
  if (pageCount > 1)
    table.insertAdjacentHTML("afterend", "<div id='buttons'></div");

  sort(sales, 1);
});
