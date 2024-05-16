import {Graph} from './graph.js';
import {PieChart} from './piechart.js';
import {getLocalization} from './localize.js';
import {LogError as logError} from '../../wailsjs/go/main/Tray.js';
import {closeNavigation, markSelectedNavigationItem} from './navigation-menu.js';
import {retrieveTheme} from './personalize.js';
import {scanTest} from './database.js';

/** Load the content of the Security Dashboard page */
export function openSecurityDashboardPage() {
  retrieveTheme();
  closeNavigation(document.body.offsetWidth);
  markSelectedNavigationItem('security-dashboard-button');
  sessionStorage.setItem('savedPage', '2');

  document.getElementById('page-contents').innerHTML = `
  <div class="dashboard">
    <div class="container-dashboard">
      <div class="dashboard-segment">
        <div class="data-segment-header">
          <p class="lang-security-stat"></p>
        </div>
        <div class="security-status">
          <p class="status-descriptor"></p>
        </div>
      </div>
      <div class="dashboard-segment">
        <div class="data-segment-header">
          <p class="lang-risk-level-counters"></p>
        </div>
        <div class="risk-counter high-risk">
          <div><p class="lang-high-risk-issues"></p></div>
          <div><p id="high-risk-counter">0</p></div>
        </div>
        <div class="risk-counter medium-risk">
          <div><p class="lang-medium-risk-issues"></p></div>
          <div><p id="medium-risk-counter">0</p></div>
        </div>
        <div class="risk-counter low-risk">
          <div><p class="lang-low-risk-issues"></p></div>
          <div><p id="low-risk-counter">0</p></div>
        </div>
        <div class="risk-counter info-risk">
          <div><p class="lang-info-risk-issues"></p></div>
          <div><p id="info-risk-counter">0</p></div>
        </div>
        <div class="risk-counter no-risk">
          <div><p class="lang-acceptable-issues"></p></div>
          <div><p id="no-risk-counter">0</p></div>
        </div>
      </div>      
    </div>
    <div class="container-dashboard">
      <div class="dashboard-segment">
        <div class="data-segment-header">
            <p class="lang-risk-level-distribution piechart-header"></p>
        </div>
        <div class="pie-chart-container">
          <canvas class="pie-chart" id="pie-chart-security"></canvas>
        </div>
      </div>
      <div class="dashboard-segment">
        <div class="data-segment-header">
          <p class="lang-risk-level-distribution"></p>
        </div>
        <div class="graph-segment-content">
          <div class="graph-buttons">
            <p class="lang-bar-graph-description">
            </p>
            <button id="dropbtn" class="dropbtn"><span class="lang-select-risks"></span></button>
            <div class="dropdown-selector" id="myDropdown">
              <p><input type="checkbox" checked="true" value="true" id="select-high-risk">
                <label for="select-high-risk" class="lang-high-risk-issues"></label><br>
              </p>
              <p><input type="checkbox" checked="true" value="true" id="select-medium-risk">
                <label for="select-medium-risk" class="lang-medium-risk-issues"></label>
              </p>
              <p><input type="checkbox" checked="true" value="true" id="select-low-risk">
                <label for="select-low-risk" class="lang-low-risk-issues"></label>
              </p>
              <p><input type="checkbox" checked="true" value="true" id="select-info-risk">
                <label for="select-info-risk" class="lang-info-risk-issues"></label>
              </p>
              <p><input type="checkbox" checked="true" value="true" id="select-no-risk">
                <label for="select-no-risk" class="lang-acceptable-issues"></label>
              </p>
            </div>
            <a class="interval-button">
              <p class="lang-change-interval"></p>
              <input type="number" value="1" id="graph-interval" min="1">
            </a>
          </div>
          <div class="interval-graph-container">
            <canvas id="interval-graph"></canvas>
          </div>
        </div>
      </div>
    </div>
    <div class="container-dashboard">
      <div class="dashboard-segment">
        <div class="data-segment-header">
          <p class="lang-choose-issue-description"></p>
        </div>
        <a class="issue-button lang-suggested-issue"><p></p></a>
        <a class="issue-button lang-quick-fix"><p></p></a>
        <a id="scan-now" class="issue-button lang-scan-now"></a>
      </div>
      <div class="dashboard-segment risk-areas">
        <div class="data-segment-header">
          <p class="lang-security-risk-areas"></p>
        </div>
        <div class="security-area">
          <a>
            <p>
              <span class="lang-applications"></span>
              <span class="material-symbols-outlined">apps_outage</span>
            </p>
          </a>
        </div>
        <div class="security-area">
          <a>
            <p><span class="lang-browser"></span><span class="material-symbols-outlined">travel_explore</span></p>
          </a>
        </div>
        <div class="security-area">
          <a>
            <p><span class="lang-devices"></span><span class="material-symbols-outlined">devices</span></p>
          </a>
        </div>
        <div class="security-area">
          <a>
            <p>
              <span class="lang-operating-system"></span>
              <span class="material-symbols-outlined">desktop_windows</span>
            </p>        
          </a>
        </div>
        <div class="security-area">
          <a>
            <p><span class="lang-passwords"></span><span class="material-symbols-outlined">key</span></p>
          </a>
        </div>
        <div class="security-area">
          <a>
            <p><span class="lang-other"></span><span class="material-symbols-outlined">view_cozy</span></p>
          </a>
        </div>
      </div>
    </div>
  </div>
  `;
  // Set counters on the page to the right values
  let rc = JSON.parse(sessionStorage.getItem('SecurityRiskCounters'));
  adjustWithRiskCounters(rc, document);
  setMaxInterval(rc, document);

  // Localize the static content of the dashboard
  const staticDashboardContent = [
    'lang-issues',
    'lang-high-risk-issues',
    'lang-medium-risk-issues',
    'lang-low-risk-issues',
    'lang-info-risk-issues',
    'lang-acceptable-issues',
    'lang-security-stat',
    'lang-risk-level-counters',
    'lang-risk-level-distribution',
    'lang-suggested-issue',
    'lang-quick-fix',
    'lang-scan-now',
    'lang-security-risk-areas',
    'lang-applications',
    'lang-browser',
    'lang-devices',
    'lang-operating-system',
    'lang-passwords',
    'lang-other',
    'lang-select-risks',
    'lang-change-interval',
    'lang-choose-issue-description',
    'lang-bar-graph-description',
  ];
  const localizationIds = [
    'Dashboard.Issues',
    'Dashboard.HighRisk',
    'Dashboard.MediumRisk',
    'Dashboard.LowRisk',
    'Dashboard.InfoRisk',
    'Dashboard.Acceptable',
    'Dashboard.SecurityStatus',
    'Dashboard.RiskLevelCounters',
    'Dashboard.RiskLevelDistribution',
    'Dashboard.SuggestedIssue',
    'Dashboard.QuickFix',
    'Dashboard.ScanNow',
    'Dashboard.SecurityRiskAreas',
    'Dashboard.Applications',
    'Dashboard.Browser',
    'Dashboard.Devices',
    'Dashboard.OperatingSystem',
    'Dashboard.Passwords',
    'Dashboard.Other',
    'Dashboard.SelectRisks',
    'Dashboard.ChangeInterval',
    'Dashboard.ChooseIssueDescription',
    'Dashboard.BarGraphDescription',
  ];
  for (let i = 0; i < staticDashboardContent.length; i++) {
    getLocalization(localizationIds[i], staticDashboardContent[i]);
  }

  // Create charts
  new PieChart('pie-chart-security', rc, 'Security');
  const g = new Graph('interval-graph', rc);
  addGraphFunctions(g);
  document.getElementById('scan-now').addEventListener('click', async () => {
    await scanTest();
    rc = JSON.parse(sessionStorage.getItem('SecurityRiskCounters'));
    adjustWithRiskCounters(rc, document);
    setMaxInterval(rc, document);
    g.rc = rc;
    await g.changeGraph();
  });
}

if (typeof document !== 'undefined') {
  try {
    document.getElementById('security-dashboard-button').addEventListener('click', () => openSecurityDashboardPage());
  } catch (error) {
    logError('Error in security-dashboard.js: ' + error);
  }
}

/** Changes the risk counters to show the correct values
 *
 * @param {RiskCounters} rc Risk counters from which the data is taken
 * @param {Document} doc Document in which the counters are located
 */
export function adjustWithRiskCounters(rc, doc) {
  // change counters according to collected data
  doc.getElementById('high-risk-counter').innerHTML = rc.lastHighRisk;
  doc.getElementById('medium-risk-counter').innerHTML = rc.lastMediumRisk;
  doc.getElementById('low-risk-counter').innerHTML = rc.lastLowRisk;
  doc.getElementById('info-risk-counter').innerHTML = rc.lastInfoRisk;
  doc.getElementById('no-risk-counter').innerHTML = rc.lastNoRisk;

  const securityStatus = doc.getElementsByClassName('status-descriptor')[0];
  if (rc.lastHighRisk > 1) {
    try {
      getLocalization('Dashboard.Critical', 'status-descriptor');
    } catch (error) {
      securityStatus.innerHTML = 'Critical';
    }
    securityStatus.style.backgroundColor = rc.highRiskColor;
    securityStatus.style.color = 'rgb(255, 255, 255)';
  } else if (rc.lastMediumRisk > 1) {
    try {
      getLocalization('Dashboard.MediumConcern', 'status-descriptor');
    } catch (error) {
      securityStatus.innerHTML = 'Medium concern';
    }
    securityStatus.style.backgroundColor = rc.mediumRiskColor;
    securityStatus.style.color = 'rgb(255, 255, 255)';
  } else if (rc.lastLowRisk > 1) {
    try {
      getLocalization('Dashboard.LowConcern', 'status-descriptor');
    } catch (error) {
      securityStatus.innerHTML = 'Low concern';
    }
    securityStatus.style.backgroundColor = rc.lowRiskColor;
    securityStatus.style.color = 'rgb(0, 0, 0)';
  } else if (rc.lastInfoRisk > 1) {
    try {
      getLocalization('Dashboard.InfoConcern', 'status-descriptor');
    } catch (error) {
      securityStatus.innerHTML = 'Informative';
    }
    securityStatus.style.backgroundColor = rc.infoColor;
    securityStatus.style.color = 'rgb(0, 0, 0)';
  } else {
    try {
      getLocalization('Dashboard.NoConcern', 'status-descriptor');
    } catch (error) {
      securityStatus.innerHTML = 'Acceptable';
    }
    securityStatus.style.backgroundColor = rc.noRiskColor;
    securityStatus.style.color = 'rgb(0, 0, 0)';
  }
}

/** Set the max number input of the 'graph-interval' element
 *
 * @param {RiskCounters} rc Risk counters from which the max count is taken
 * @param {Document} doc Document in which the counters are located
 */
export function setMaxInterval(rc, doc) {
  doc.getElementById('graph-interval').max = rc.count;
}

/** Adds event listeners to elements in graph-row section of the dashboard page
 *
 * @param {Graph} g Graph class containing the functions to be called
 */
export function addGraphFunctions(g) {
  document.getElementById('dropbtn').addEventListener('click', () => g.graphDropdown());
  document.getElementById('graph-interval').addEventListener('change', () => g.changeGraph());
  document.getElementById('select-high-risk').addEventListener('change', () => g.toggleRisks('high'));
  document.getElementById('select-medium-risk').addEventListener('change', () => g.toggleRisks('medium'));
  document.getElementById('select-low-risk').addEventListener('change', () => g.toggleRisks('low'));
  document.getElementById('select-info-risk').addEventListener('change', () => g.toggleRisks('info'));
  document.getElementById('select-no-risk').addEventListener('change', () => g.toggleRisks('no'));
}
