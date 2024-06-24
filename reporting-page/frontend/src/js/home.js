import {PieChart} from './piechart.js';
import {getLocalization} from './localize.js';
import {closeNavigation, markSelectedNavigationItem} from './navigation-menu.js';
import {retrieveTheme} from './personalize.js';
import {scanTest} from './database.js';
import {LogError as logError, LogDebug as logDebug} from '../../wailsjs/go/main/Tray.js';
import {GetImagePath as getImagePath, GetLighthouseState as getLighthouseState} from '../../wailsjs/go/main/App.js';
import {openIssuePage} from './issue.js';
import {saveProgress, shareProgress, selectSocialMedia, setImage, socialMediaSizes} from './share.js';
import data from '../databases/database.en-GB.json' assert { type: 'json' };
import {showModal} from './settings.js';

let lighthousePath;
/** Load the content of the Home page */
export async function openHomePage() {
  // Load the video background path
  const lighthouseState = await getLighthouseState();
  switch (lighthouseState) {
  case 0:
    lighthousePath = 'state0.mkv';
    break;
  case 1:
    lighthousePath = 'state1.mkv';
    break;
  case 2:
    lighthousePath = 'state2.mkv';
    break;
  case 3:
    lighthousePath = 'state3.mkv';
    break;
  case 4:
    lighthousePath = 'state4.mkv';
    break;
  default:
    lighthousePath = 'state0.mkv';
  }

  const lighthouseFullPath = await getImagePath('gamification/' + lighthousePath);
  logDebug('lighthouseState: ' + lighthouseFullPath);

  retrieveTheme();
  closeNavigation(document.body.offsetWidth);
  markSelectedNavigationItem('home-button');
  sessionStorage.setItem('savedPage', 1);

  document.getElementById('page-contents').innerHTML = `
  <video autoplay muted loop class="video-background">
        <source id="lighthouse-background" type="video/mp4">
        Your browser does not support HTML5 video.
  </video>
  <div class="home-page">
    <div class="container-home"> 
      <div class="data-segment">
        <div class="data-segment-header">
          <p class="lang-piechart-header"></p>
        </div>
        <div class="pie-chart-container">
          <canvas class="pie-chart" id="pie-chart-home"></canvas>
        </div>
      </div>
      <div class="data-segment">
        <div class="data-segment-header">
          <p class="lang-choose-issue-description"></p>
        </div>
        <a id="suggested-issue" class="issue-button lang-suggested-issue"></a>
        <a id="scan-now" class="issue-button lang-scan-now"></a>
        <a id="share-progress" class="issue-button lang-share-button"></a>
      </div>
      <div id="progress-segment" class="data-segment">
        <div class="data-segment-header">
          <p class="lang-lighthouse-progress"></p>
          <div id="lighthouse-progress-hoverbox">
            <img id="lighthouse-progress-tooltip">
            <p class="lighthouse-progress-tooltip-text lang-tooltip-text"></p>
          </div>
        </div>
        <div class="progress-container">
          <div class="progress-bar" id="progress-bar"></div>
        </div>
        <p id="progress-text"></p>
      </div>
    </div>
  </div>
  <div id="share-modal" class="modal">
    <div class="modal-content">
      <div class="modal-header">
        <span id="close-share-modal" class="close">&times;</span>
        <p class="lang-share-text"></p>
      </div>
      <div id="share-node" class="modal-body share-image">
      </div>
      <div id="share-buttons" class="modal-body">
        <a id="share-save-button" class="modal-button share-button lang-save-text"></a>
        <a class="share-button-break">|</a>
        <a id="select-facebook" class="select-button selected">Facebook</a>
        <a id="select-x" class="select-button">X</a>
        <a id="select-linkedin" class="select-button">LinkedIn</a>
        <a id="select-instagram" class="select-button">Instagram</a>
        <a class="share-button-break">|</a>
        <a id="share-button" class="modal-button share-button lang-share"></a>
      </div>
    </div>
  </div>
  `;

  document.getElementById('lighthouse-background').src = lighthouseFullPath;

  const tooltip = await getImagePath('tooltip.png');
  document.getElementById('lighthouse-progress-tooltip').src = tooltip;

  const rc = JSON.parse(sessionStorage.getItem('RiskCounters'));
  new PieChart('pie-chart-home', rc, 'Total');

  // Localize the static content of the home page
  const staticHomePageContent = [
    'lang-piechart-header',
    'lang-suggested-issue',
    'lang-scan-now',
    'lang-title-medals',
    'lang-choose-issue-description',
    'lang-share-button',
    'lang-share-text',
    'lang-save-text',
    'lang-share',
    'lang-lighthouse-progress',
    'lang-tooltip-text',
  ];
  const localizationIds = [
    'Dashboard.RiskLevelDistribution',
    'Dashboard.SuggestedIssue',
    'Dashboard.ScanNow',
    'Dashboard.Medals',
    'Dashboard.ChooseIssueDescription',
    'Dashboard.ShareButton',
    'Dashboard.ShareText',
    'Dashboard.SaveText',
    'Dashboard.Share',
    'Dashboard.LighthouseProgress',
    'Dashboard.TooltipText',
  ];
  for (let i = 0; i < staticHomePageContent.length; i++) {
    getLocalization(localizationIds[i], staticHomePageContent[i]);
  }

  document.getElementById('scan-now').addEventListener('click', () => scanTest(true));
  document.getElementById('suggested-issue').addEventListener('click', () => suggestedIssue(''));
  document.getElementById('share-progress').addEventListener('click', () => showModal('share-modal'));
  document.getElementById('share-save-button').addEventListener('click',
    () => saveProgress(document.getElementById('share-node')));
  document.getElementById('share-button').addEventListener('click', () => shareProgress());

  document.getElementById('select-facebook').addEventListener('click', () => selectSocialMedia('facebook'));
  document.getElementById('select-x').addEventListener('click', () => selectSocialMedia('x'));
  document.getElementById('select-linkedin').addEventListener('click', () => selectSocialMedia('linkedin'));
  document.getElementById('select-instagram').addEventListener('click', () => selectSocialMedia('instagram'));

  // Progress bar
  document.addEventListener('DOMContentLoaded', () => {
    const progressBar = document.getElementById('progress-bar');
    const progressText = document.getElementById('progress-text');

    // Assuming the points are stored in local storage under the key 'userPoints'
    const userPoints = parseInt(localStorage.getItem('userPoints')) || 0;
    const pointsToNextState = 100; // The points required to reach the next state

    // Calculate the progress percentage
    const progressPercentage = Math.min((userPoints / pointsToNextState) * 100, 100);

    // Update the progress bar width and text
    progressBar.style.width = progressPercentage + '%';
    progressText.textContent = `${userPoints} / ${pointsToNextState} (${progressPercentage.toFixed(2)}%)`;
  });

  // on startup set the social media to share to facebook
  sessionStorage.setItem('ShareSocial', JSON.stringify(socialMediaSizes['facebook']));
  setImage(document.getElementById('share-node'), document.getElementById('progress-segment'));
}

/** Opens the issue page of the issue with the highest risk level
 *
 * @param {string} type Type of issue to open the issue page of (e.g. 'Security', 'Privacy', and '' for all types)
*/
export function suggestedIssue(type) {
  // Get the issues from the database
  const issues = JSON.parse(sessionStorage.getItem('ScanResult'));

  // Skip informative issues
  let issueIndex = 0;
  let maxSeverityIssue = issues[issueIndex].issue_id;
  let maxSeverityResult = issues[issueIndex].result_id;
  while (getSeverity(maxSeverityIssue, maxSeverityResult) === 4 ||
        getSeverity(maxSeverityIssue, maxSeverityResult) === undefined) {
    issueIndex++;
    maxSeverityIssue = issues[issueIndex].issue_id;
    maxSeverityResult = issues[issueIndex].result_id;
  }

  // Find the issue with the highest severity
  for (let i = 0; i < issues.length; i++) {
    const severity = getSeverity(issues[i].issue_id, issues[i].result_id);
    if (getSeverity(maxSeverityIssue, maxSeverityResult) < severity &&
      severity !== 4) {
      if (type == '' || data[issues[i].issue_id].Type === type) {
        maxSeverityIssue = issues[i].issue_id;
        maxSeverityResult = issues[i].result_id;
      }
    }
  }

  // Open the issue page of the issue with the highest severity
  openIssuePage(maxSeverityIssue, maxSeverityResult, 'home');
  document.getElementById('scan-now').addEventListener('click', () => scanTest(true));
}

/**
 * Gets the severity of an issue
 * @param {string} issueId id of the issue
 * @param {string} resultId result id of the issue
 *
 * @return {string} severity
 */
export function getSeverity(issueId, resultId) {
  const issue = data[issueId];
  if (issue == undefined) return undefined;
  const issueData = issue[resultId];
  if (issueData == undefined) return undefined;
  return issueData.Severity;
}

/* istanbul ignore next */
if (typeof document !== 'undefined') {
  try {
    document.getElementById('logo-button').addEventListener('click', () => openHomePage());
    document.getElementById('home-button').addEventListener('click', () => openHomePage());
  } catch (error) {
    logError('Error in security-dashboard.js: ' + error);
  }
}


window.onload = function() {
  const savedImage = localStorage.getItem('picture');
  const savedText = localStorage.getItem('title');
  if (savedImage) {
    const logo = document.getElementById('logo');
    logo.src = savedImage;
  }
  if (savedText) {
    const title = document.getElementById('title');
    title.textContent = savedText;
  }
};


