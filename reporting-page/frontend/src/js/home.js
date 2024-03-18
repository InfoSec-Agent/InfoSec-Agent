import * as piechart from "./piechart";
import {ScanNow} from '../../wailsjs/go/main/Tray';
import {LogPrint, WindowSetAlwaysOnTop} from '../../wailsjs/runtime/runtime';

/** Load the content of the Home page */
function openHomePage() {
<<<<<<< HEAD
  document.getElementById("page-contents").innerHTML = `
  <div class="container-data">       
    <div class="data-column risk-counters">     
      <div class="data-column piechart">
        <canvas id="pieChart"></canvas>
      </div>
=======
    document.getElementById("page-contents").innerHTML = `
    <div class="container-data">       
        <div class="data-column risk-counters">     
            <div class="data-column piechart">
                <canvas id="pieChart"></canvas>
            </div>
        </div>
        <div class="data-column issue-buttons">
            <H2>You have some issues you can fix. 
                To start resolving a issue either navigate to the issues page, or pick a suggested issue below
            </H2>
            <a class="issue-button">Suggested Issue</a>
            <a class="issue-button">Quick Fix</a>
            <a class="issue-button" id="scan-button">Scan Now</a>
        </div>
>>>>>>> main
    </div>
    <div class="data-column issue-buttons">
      <H2>You have some issues you can fix. 
        To start resolving a issue either navigate to the issues page, or pick a suggested issue below
      </H2>
      <a class="issue-button">Suggested Issue</a>
      <a class="issue-button">Quick Fix</a>
    </div>
  </div>
  <h2 class="title-medals">Medals</h2>
  <div class="container">  
    <div class="medal-layout">
      <img src="src/assets/images/img_medal1.jpg" alt="Photo of medal">
      <p class="medal-name"> Medal 1</p>
      <p> 01-03-2024</p>
    </div>
    <div class="medal-layout">
      <img src="src/assets/images/img_medal1.jpg" alt="Photo of medal">
      <p class="medal-name"> Medal 2</p>
      <p> 01-03-2024</p>
    </div>
    <div class="medal-layout">
      <img src="src/assets/images/img_medal1.jpg" alt="Photo of medal">
      <p class="medal-name"> Medal 3</p>
      <p> 01-03-2024</p>
    </div><div class="medal-layout">
      <img src="src/assets/images/img_medal1.jpg" alt="Photo of medal">
      <p class="medal-name"> Medal 1</p>
      <p> 01-03-2024</p>
    </div>
  </div>
  `;  

<<<<<<< HEAD
  CreatePieChart();
=======
    CreatePieChart();

    document.getElementById("scan-button").addEventListener("click", () => scanNow());
>>>>>>> main
}

document.getElementById("logo-button").addEventListener("click", () => openHomePage());
document.getElementById("home-button").addEventListener("click", () => openHomePage());


function scanNow() {
    ScanNow()
    .then((result) => {
    })
    .catch((err) => {
        console.error(err);
    });
}

//#region PieChart

// Reusable snippit for other files
let pieChart;

/** Creates a pie chart for risks */
function CreatePieChart() {
  pieChart = new Chart("pieChart", {
    type: "doughnut",
    data: piechart.GetData(),
    options: piechart.GetOptions()
  });
}

//#endregion

document.onload = openHomePage();