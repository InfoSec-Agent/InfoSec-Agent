import {PieChart} from "./piechart";
import {ScanNow} from '../../wailsjs/go/main/Tray';
import {RiskCounters} from "./risk-counters";
import { GetLocalization } from './localize.js';

/** Load the content of the Home page */
function openHomePage() {
  document.getElementById("page-contents").innerHTML = `
  <div class="container-data">       
    <div class="data-column risk-counters">     
      <div class="data-column piechart">
        <canvas id="pieChart"></canvas>
      </div>
    </div>
    <div class="data-column issue-buttons">
      <H2 class="choose-issue-description">You have some issues you can fix. 
        To start resolving a issue either navigate to the issues page, or pick a suggested issue below
      </H2>
      <a class="issue-button suggested-issue">Suggested Issue</a>
      <a class="issue-button quick-fix">Quick Fix</a>
      <a class="issue-button scan-now">Scan Now</a>
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
    let rc = new RiskCounters();
    new PieChart("pieChart",rc);

    // Localize the static content of the home page
    let staticHomePageConent = [
      "suggested-issue", 
      "quick-fix", 
      "scan-now", 
      "title-medals", 
      "security-status",
      "high-risk-issues", 
      "medium-risk-issues",
      "low-risk-issues",
      "safe-issues",
      "choose-issue-description"
      ]
      let localizationIds = [
        "Dashboard.SuggestedIssue", 
        "Dashboard.QuickFix", 
        "Dashboard.ScanNow", 
        "Dashboard.Medals", 
        "Dashboard.SecurityStatus",
        "Dashboard.HighRisk", 
        "Dashboard.MediumRisk",
        "Dashboard.LowRisk",
        "Dashboard.Safe",
        "Dashboard.ChooseIssueDescription"
      ]
      for (let i = 0; i < staticHomePageConent.length; i++) {
          GetLocalization(localizationIds[i], staticHomePageConent[i])
    }

    document.getElementsByClassName("scan-now")[0].addEventListener("click", () => scanNow());
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

document.onload = openHomePage();