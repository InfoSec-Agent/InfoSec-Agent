import dataDe from '../databases/database.de.json' assert { type: 'json' };
import dataEnGB from '../databases/database.en-GB.json' assert { type: 'json' };
import dataEnUS from '../databases/database.en-US.json' assert { type: 'json' };
import dataEs from '../databases/database.es.json' assert { type: 'json' };
import dataFr from '../databases/database.fr.json' assert { type: 'json' };
import dataNl from '../databases/database.nl.json' assert { type: 'json' };
import dataPt from '../databases/database.pt.json' assert { type: 'json' };

import {openIssuesPage, getUserSettings} from './issues.js';
import {getLocalization} from './localize.js';
import {retrieveTheme} from './personalize.js';
import {closeNavigation, markSelectedNavigationItem} from './navigation-menu.js';
import {GetImagePath as getImagePath} from '../../wailsjs/go/main/App.js';
import {LogError as logError, LogDebug as logDebug} from '../../wailsjs/go/main/Tray.js';
import {scanTest} from './database.js';
import {openAllChecksPage} from './all-checks.js';
import {openHomePage} from './home.js';

let stepCounter = 0;
const issuesWithResultsShow =
    ['11', '21', '60', '70', '80', '90', '100', '110', '160', '172', '173',
      '201', '230', '250', '260', '271', '300', '311', '320', '351', '361'];

/** Update contents of solution guide
 *
 * @param {HTMLParagraphElement} solutionText Element in which textual solution step is shown
 * @param {HTMLImageElement} solutionScreenshot Element in which screenshot of solution step is shown
 * @param {*} issue Issue to update the solution step for
 * @param {int} stepCounter Counter specifying the current step
 */
export async function updateSolutionStep(solutionText, solutionScreenshot, issue, stepCounter) {
  solutionText.innerHTML = `${stepCounter + 1}. ${getVersionSolution(issue, stepCounter)}`;
  const screenshot = await getVersionScreenshot(issue, stepCounter);
  logDebug('screenshot source updateSolutionStep: ' + screenshot);
  solutionScreenshot.src = screenshot;
  // Hide/show buttons based on the current step
  const previousButton = document.getElementById('previous-button');
  const nextButton = document.getElementById('next-button');
  const scanButton = document.getElementById('scan-button');

  if (previousButton && nextButton) {
    if (stepCounter === 0) {
      previousButton.style.display = 'none';
    } else {
      previousButton.style.display = 'block';
    }
    if (stepCounter === issue.Solution.length - 1) {
      nextButton.style.display = 'none';
      scanButton.style.display = ' block';
    } else {
      nextButton.style.display = 'block';
      scanButton.style.display = 'none';
    }
  }
}

/** Go to next step of solution guide
 *
 * @param {HTMLParagraphElement} solutionText Element in which textual solution step is shown
 * @param {HTMLImageElement} solutionScreenshot Element in which screenshot of solution step is shown
 * @param {*} issue Issue to update the solution step for
 */
export function nextSolutionStep(solutionText, solutionScreenshot, issue) {
  if (stepCounter < issue.Solution.length - 1) {
    stepCounter++;
    updateSolutionStep(solutionText, solutionScreenshot, issue, stepCounter);
  }
}

/** Go to previous step of solution guide
 *
 * @param {HTMLParagraphElement} solutionText Element in which textual solution step is shown
 * @param {HTMLImageElement} solutionScreenshot Element in which screenshot of solution step is shown
 * @param {*} issue Issue to update the solution step for
 */
export function previousSolutionStep(solutionText, solutionScreenshot, issue) {
  if (stepCounter > 0) {
    stepCounter--;
    updateSolutionStep(solutionText, solutionScreenshot, issue, stepCounter);
  }
}

/** Load the content of the issue page
 *
 * @param {string} issueId Id of the issue to open
 * @param {string} resultId result of the issue to open
 * @param {string} back if not undefined, back navigation to allChecksPage enabled
 */
export async function openIssuePage(issueId, resultId, back = undefined) {
  retrieveTheme();
  closeNavigation(document.body.offsetWidth);
  markSelectedNavigationItem('issue-button');
  stepCounter = 0;
  sessionStorage.setItem('savedPage', JSON.stringify([issueId, resultId]));
  const language = await getUserSettings();

  let currentIssue;
  switch (language) {
  case 0:
    currentIssue = dataDe[issueId];
    break;
  case 1:
    currentIssue = dataEnGB[issueId];
    break;
  case 2:
    currentIssue = dataEnUS[issueId];
    break;
  case 3:
    currentIssue = dataEs[issueId];
    break;
  case 4:
    currentIssue = dataFr[issueId];
    break;
  case 5:
    currentIssue = dataNl[issueId];
    break;
  case 6:
    currentIssue = dataPt[issueId];
    break;
  default:
    currentIssue = dataEnGB[issueId];
    break;
  }
  const issueData = currentIssue[resultId];
  let riskLevel = '';
  if (issueData.Severity == 0) {
    riskLevel = '<span class="risk-indicator lang-acceptable-risk" data-severity="0"></span>';
  } else if (issueData.Severity == 1) riskLevel = '<span class="risk-indicator lang-low" data-severity="1"></span>';
  else if (issueData.Severity == 2) riskLevel = '<span class="risk-indicator lang-medium" data-severity="2"></span>';
  else if (issueData.Severity == 3) riskLevel = '<span class="risk-indicator lang-high" data-severity="3"></span>';
  else riskLevel = '<span class="risk-indicator lang-info" data-severity="4"></span>';
  if (resultId < 0) {
    const pageContents = document.getElementById('page-contents');
    pageContents.innerHTML = `
    <div class="issue-data">
      <div class="issue-page-header" data-severity="${issueData.Severity}">
        <h2 class="issue-name">${issueData.Name}</h2>
        ${riskLevel}
      </div>
      <div class="issue-information">
        <p id="description">${currentIssue.Information}</p>
        <h2 id="solution">${issueData.Failed}</h2>
        <div class="issue-solution">
          <p id="solution-text">${getVersionSolution(issueData, stepCounter)}</p>
        </div>
        <div class="solution-buttons">
          <div class="button-box">
            <div class="lang-scan-again button" id="scan-button"></div>
          </div>
        </div>
        <div class="button" id="back-button"></div>
      </div>
    </div>
    `;
  } else if (issueData.Severity == 0 && issueData.Screenshots.length == 0) {
    // Check issue severity and if the issue has no screenshots, if so, display that there is no issue (acceptable)
    const pageContents = document.getElementById('page-contents');
    pageContents.innerHTML = `
    <div class="issue-data">
      <div class="issue-page-header" data-severity="${issueData.Severity}">
        <h2 class="issue-name">${issueData.Name}</h2>
        ${riskLevel}
      </div>
      <div class="issue-information">
        <p id="description">${currentIssue.Information}</p>
        <h2 id="solution" class="lang-acceptable"></h2>
        <div class="issue-solution">
          <p id="solution-text">${getVersionSolution(issueData, stepCounter)}</p>
          <div class="solution-buttons">
            <div class="button-box">
              <div class="lang-scan-again button" id="scan-button"></div>
            </div>
          </div>
        </div>
        <div class="button" id="back-button"></div>
      </div>
    </div>
    `;
  } else { // Issue has screenshots, display the solution guide
    const pageContents = document.getElementById('page-contents');
    if (issuesWithResultsShow.includes(issueId.toString() + resultId.toString())) {
      pageContents.innerHTML = parseShowResult(issueId, resultId, currentIssue);
    } else {
      pageContents.innerHTML = `
      <div class="issue-data">
        <div class="issue-page-header" data-severity="${issueData.Severity}">
          <h2 class="issue-name">${issueData.Name}</h2>
          ${riskLevel}
        </div>
        <div class="issue-information">
          <p id="description">${currentIssue.Information}</p>
          <h2 id="solution" class="lang-solution"></h2>
          <div class="issue-solution">
            <p id="solution-text">${stepCounter +1}. ${getVersionSolution(issueData, stepCounter)}</p>
            <input type="checkbox" id="zoom-check">
            <label for="zoom-check">
              <img class="zoom-img" style='display:block; width:65%;height:auto' id="step-screenshot"></img>
            </label>
            <div class="solution-buttons">
              <div class="button-box">
                <div id="previous-button" class="lang-previous-button button"></div>
                <div id="next-button" class="lang-next-button button"></div>
                <div class="lang-scan-again button" id="scan-button"></div>
              </div>
            </div>
          </div>
          <div class="button" id="back-button"></div>
        </div>
      </div>
      `;
    }

    // Add functions to page for navigation
    const solutionText = document.getElementById('solution-text');
    const solutionScreenshot = document.getElementById('step-screenshot');
    document.getElementById('next-button').addEventListener('click', () =>
      nextSolutionStep(solutionText, solutionScreenshot, issueData));
    document.getElementById('previous-button').addEventListener('click', () =>
      previousSolutionStep(solutionText, solutionScreenshot, issueData));

    // Initial check to hide/show buttons
    try {
      await updateSolutionStep(solutionText, solutionScreenshot, issueData, stepCounter);
    } catch (error) {
      logError('Error in updateSolutionStep: ' + error);
    }
  }

  if (back == undefined) {
    document.getElementById('back-button').addEventListener('click', () => openIssuesPage());
    document.getElementById('back-button').classList.add('lang-back-button-issues');
  } else if (back == 'home') {
    document.getElementById('back-button').addEventListener('click', () => openHomePage());
    document.getElementById('back-button').classList.add('lang-back-button-home');
  } else {
    document.getElementById('back-button').addEventListener('click', () => openAllChecksPage(back));
    document.getElementById('back-button').classList.add('lang-back-button-checks');
  }

  const texts = ['lang-findings', 'lang-solution', 'lang-previous-button',
    'lang-next-button', 'lang-back-button-issues', 'lang-back-button-home', 'lang-back-button-checks', 'lang-port',
    'lang-password', 'lang-acceptable', 'lang-cookies', 'lang-permissions', 'lang-scan-again',
    'lang-info', 'lang-medium', 'lang-high', 'lang-low', 'lang-acceptable-risk',
  ];
  const localizationIds = ['Issues.Findings', 'Issues.Solution', 'Issues.Previous',
    'Issues.Next', 'Issues.BackIssues', 'Issues.BackHome', 'Issues.BackChecks', 'Issues.Port', 'Issues.Password',
    'Issues.Acceptable', 'Issues.Cookies', 'Issues.Permissions', 'Issues.ScanAgain',
    'Dashboard.InfoRisk', 'Dashboard.MediumRisk', 'Dashboard.HighRisk', 'Dashboard.LowRisk', 'Dashboard.Acceptable',
  ];
  for (let i = 0; i < texts.length; i++) {
    getLocalization(localizationIds[i], texts[i]);
  }

  document.getElementById('scan-button').addEventListener('click', async () => {
    await scanTest(true);
    const issueData = JSON.parse(sessionStorage.getItem('ScanResult'));
    const resultId = findResultId(issueData, issueId);
    sessionStorage.setItem('savedPage', JSON.stringify([issueId, resultId]));
    openIssuePage(issueId, resultId);
  });

  const image = document.getElementsByClassName('zoom-img');
  if (image.length > 0) {
    image[0].addEventListener('click', function() {
      if (image[0].classList.contains('zoomed')) {
        image[0].classList.remove('zoomed');
      } else {
        image[0].classList.add('zoomed');
      }
    });
  }
}

/**
 * Finds the result_id corresponding to the given issueId in the provided data array.
 *
 * @param {Array} data - An array of objects where each object contains an issue_id and a result_id.
 * @param {number|string} issueId - The issue ID to search for in the data array. This can be a number or a string.
 * @return {number|null} The result_id corresponding to the provided issueId, or null if not found.
 */
function findResultId(data, issueId) {
  for (const item of data) {
    if (parseInt(item.issue_id) === parseInt(issueId)) {
      return item.result_id;
    }
  }
  return null;
}

/** Check if the issue is a show result issue
 *
 * @param {string} issue checks if the issue is a show result issue
 * @return {boolean} if the issue is a show result issue
 */
export function checkShowResult(issue) {
  return issue.Name.includes('Applications with');
}

/** Parse the show result of an issue
 *
 * @param {string} issueId of the issue
 * @param {string} resultId of the issue we are looking at
 * @param {string} currentIssue current issue we are looking at
 * @return {string} result of the show result
 */
export function parseShowResult(issueId, resultId, currentIssue) {
  let issues = [];
  issues = JSON.parse(sessionStorage.getItem('ScanResult'));
  let resultLine = '';

  switch (issueId.toString() + resultId.toString()) {
  case '11':
    generateBulletList(issues, 1);
    break;
  case '21':
    generateBulletList(issues, 2);
    break;
  case '60':
    resultLine = permissionShowResults(issues);
    break;
  case '70':
    resultLine = permissionShowResults(issues);
    break;
  case '80':
    resultLine = permissionShowResults(issues);
    break;
  case '90':
    resultLine = permissionShowResults(issues);
    break;
  case '100':
    resultLine = permissionShowResults(issues);
    break;
  case '110':
    resultLine += `<p class="lang-port"></p>`;
    const portTable = processPortsTable(issues.find((issue) => issue.issue_id === 11).result);
    resultLine += `<table class="issues-table ports-table">`;
    resultLine += `<thead><tr><th>Process</th><th>Port(s)</th></tr></thead>`;
    portTable.forEach((entry) => {
      resultLine += `<tr><td style="width: 30%">${entry.portProcess}</td>
        <td style="width: 30%">${entry.ports.join('<br>')}</td></tr>`;
    });
    resultLine += '</table>';
    break;
  case '160':
    issues.find((issue) => issue.issue_id === 16).result.forEach((issue) => {
      resultLine += `<p class="lang-password"></p>`;
      resultLine += `<p class="information">${issue}</p>`;
    });
    break;
  case '172':
    generateBulletList(issues, 17);
    break;
  case '173':
    generateBulletList(issues, 17);
    break;
  case '201':
    generateBulletList(issues, 20);
    break;
  case '230':
    generateBulletList(issues, 23);
    break;
  case '250':
    generateBulletList(issues, 25);
    break;
  case '260':
    generateBulletList(issues, 26);
    break;
  case '271':
    resultLine += '<p class="lang-cookies"</p>';
    resultLine += cookiesTable(issues.find((issue) => issue.issue_id === 27).result);
    break;
  case '300':
    generateBulletList(issues, 30);
    break;
  case '311':
    generateBulletList(issues, 31);
    break;
  case '320':
    const cisTable = cisregistryTable(issues.find((issue) => issue.issue_id === 32).result);
    resultLine += `<table class = "issues-table">`;
    cisTable.forEach((entry) => {
      resultLine += `<tr><td style="width: 30%; word-break: break-all">${entry.registryKey}</td>
        <td>${entry.values.join('<br>')}</td></tr>`;
    });
    resultLine += '</table>';
    break;
  case '351':
    resultLine += '<p class="lang-cookies"</p>';
    resultLine += cookiesTable(issues.find((issue) => issue.issue_id === 35).result);
    break;
  case '361':
    resultLine += '<p class="lang-cookies"</p>';
    resultLine += cookiesTable(issues.find((issue) => issue.issue_id === 36).result);
    break;
  default:
    break;
  }

  /**
   * Generate a bullet list for each entry of a result of certain issues
   * @param {string} issues to generate a bullet list for
   * @param {int} issueId of the issue
   * @return {string} html tags for a bullet list
   */
  function generateBulletList(issues, issueId) {
    resultLine += `<ul>`;
    issues.find((issue) => issue.issue_id == issueId).result.forEach((issue) => {
      resultLine += `<li>${issue}</li>`;
    });
    resultLine += `</ul>`;
    return resultLine;
  }

  /**
   *
   * @param {string} issues with the permission results
   * @return {string} resultLine with the permission results
   */
  function permissionShowResults(issues) {
    let applications = '<ul>';
    issues.forEach((issue) => {
      if (issue.issue_id.toString() === issueId.toString()) {
        const issueResult = issue.result;
        issueResult.forEach((application) => {
          applications += `<li>${application}</li>`;
        });
      }
    });
    applications += '</ul>'; // Close the list
    resultLine += `<p class="lang-permissions"></p>`;
    resultLine += `${applications}`;
    return resultLine;
  }

  /**
   * Create a table for the CIS registry issues
   * @param {string} issues list of incorrect registry keys
   * @return {*[]} table with registry keys and values
   */
  function cisregistryTable(issues) {
    const table = [];
    let currentKey = null;
    let currentValues = [];

    issues.forEach((issue) => {
      if (issue.includes('SYSTEM') || issue.includes('SOFTWARE')) {
        if (currentKey) {
          table.push({registryKey: currentKey, values: currentValues});
        }
        currentKey = issue;
        currentValues = [];
      } else if (currentKey) {
        currentValues.push(issue);
      }
    });

    if (currentKey) {
      table.push({registryKey: currentKey, values: currentValues});
    }

    return table;
  }

  /**
   * Create a table for the process ports
   * @param {string} issues list of processes and ports
   * @return {*[]} table with process names and ports
   */
  function processPortsTable(issues) {
    const table = [];
    issues.forEach((issue) => {
      const parts = issue.split(/[ ,]+/); // Split on space and comma
      const processIndex = parts.indexOf('process:');
      const portIndex = parts.indexOf('port:');

      if (processIndex !== -1 && portIndex !== -1) {
        const processName = parts[processIndex + 1];
        const ports = new Set(parts.slice(portIndex + 1));
        table.push({portProcess: processName, ports: Array.from(ports)});
      }
    });

    return table;
  }

  /**
   * Create a table to display found (possible) tracking cookies
   * @param {string} issues list of cookies and their host
   * @return {string} HTML table with cookies and their host
   */
  function cookiesTable(issues) {
    const cookiesByHost = {};
    for (let i = 0; i < issues.length; i += 2) {
      const host = issues[i+1];

      if (!cookiesByHost[host]) {
        cookiesByHost[host] = true;
      }
    }

    // Generate HTML for table
    let tableHTML = '<table class="issues-table">';
    for (const host in cookiesByHost) {
      if (cookiesByHost.hasOwnProperty(host)) {
        tableHTML += `<tr><td style="width: 30%; word-break: break-all">${host}</td></tr>`;
      }
    }
    tableHTML += '</table>';

    return tableHTML;
  }

  const issueData = currentIssue[resultId];
  let riskLevel = '';
  if (issueData.Severity == 0) {
    riskLevel = '<span class="risk-indicator lang-acceptable-risk" data-severity="0"></span>';
  } else if (issueData.Severity == 1) riskLevel = '<span class="risk-indicator lang-low" data-severity="1"></span>';
  else if (issueData.Severity == 2) riskLevel = '<span class="risk-indicator lang-medium" data-severity="2"></span>';
  else if (issueData.Severity == 3) riskLevel = '<span class="risk-indicator lang-high" data-severity="3"></span>';
  else riskLevel = '<span class="risk-indicator lang-info" data-severity="4"></span>';
  const result = `
  <div class="issue-data">
    <div class="issue-page-header" data-severity="${issueData.Severity}">
      <h2 class="issue-name">${issueData.Name}</h2>
      ${riskLevel}
    </div>
    <div class="issue-information">
      <p id="description">${currentIssue.Information}</p>
      <h2 id="information" class="lang-findings"></h2>
      <p id="findings">${resultLine}</p>
      <h2 id="solution" class="lang-solution"></h2>
      <div class="issue-solution">
        <p id="solution-text">${stepCounter +1}. ${getVersionSolution(issueData, stepCounter)}</p>
        <input type="checkbox" id="zoom-check">
        <label for="zoom-check">
          <img class="zoom-img" style='display:block; width:65%;height:auto' id="step-screenshot"></img>
        </label>
        <div class="solution-buttons">
          <div class="button-box">
            <div id="previous-button" class="button lang-previous-button"></div>
            <div id="next-button" class="button lang-next-button"></div>
            <div class="lang-scan-again button" id="scan-button"></div>
          </div>
        </div>
      </div>
      <div class="lang-back-button button" id="back-button"></div>
    </div>
  </div>
  `;

  return result;
}

/**
 * Get the screenshot for an issue with the correct windows version detected.
 * If the version is not found, returns windows 11 screenshot.
 * If the screenshot is not found, returns no path.
 * @param {string} issue issue of which to get the screenshot
 * @param {int} index index of the desired screenshot in the list of screenshots
 * @return {string} path to the screenshot
 */
export async function getVersionScreenshot(issue, index) {
  const imagesDirectory = await getImagePath('');
  const version = 'windows-' + sessionStorage.getItem('WindowsVersion') + '/';

  let screenshot;
  if (version === 'windows-10/' && issue.ScreenshotsWindows10 !== undefined) {
    screenshot = issue.ScreenshotsWindows10[index];
  } else {
    screenshot = issue.Screenshots[index];
  }

  // Return empty source if no screenshot is found
  if (screenshot == undefined) {
    return '';
  }
  // Construct full image path
  return imagesDirectory + version + screenshot;
}

/**
 * Get the solution for an issue with the correct windows version detected.
 * @param {string} issue issue of which to get the solution
 * @param {int} index index of the desired solution in the list of solutions
 * @return {string} solution
 */
export function getVersionSolution(issue, index) {
  let solution = issue.Solution[index];
  if (solution == undefined) solution = '';
  switch (sessionStorage.getItem('WindowsVersion')) {
  case ('10'):
    const solutions = issue.SolutionWindows10;
    if (solutions !== undefined) solution = solutions[index];
    return solution;
  case ('11'):
    return solution;
  default:
    return solution;
  }
}

/**
 * Function to scroll to an element
 * @param {HTMLElement} element node to scroll to
 */
export function scrollToElement(element) {
  /* istanbul ignore next */
  element.scrollIntoView({behavior: 'smooth', block: 'center', inline: 'nearest'});
}
