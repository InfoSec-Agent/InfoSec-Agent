import data from "../database.json" assert { type: "json" };
import { openIssuesPage } from "./issues.js";
import { GetLocalization } from './localize.js';

let stepCounter = 0;

/** Update contents of solution guide 
 * 
 * @param {[string]} solution List of textual solution steps
 * @param {[image]} screenshots List of images of solution steps
 * @param {int} stepCounter Counter specifying the current step

 */ 
export function updateSolutionStep(solution, screenshots, stepCounter) {
  const solutionStep = document.getElementById("solution-text");
  solutionStep.innerHTML = solution[stepCounter];
  document.getElementById("step-screenshot").src = screenshots[stepCounter];
}

/** Go to next step of solution guide 
 * 
 * @param {[string]} solution List of textual solution steps
 * @param {[image]} screenshots List of images of solution steps
 */ 
export function nextSolutionStep(solution, screenshots) {
  if (stepCounter < solution.length - 1) {
    stepCounter++;
    updateSolutionStep(solution, screenshots, stepCounter);
  }
}

/** Go to previous step of solution guide 
 * 
 * @param {[string]} solution List of textual solution steps
 * @param {[image]} screenshots List of images of solution steps
 */ 
export function previousSolutionStep(solution, screenshots) {
  if (stepCounter > 0) {
    stepCounter--;
    updateSolutionStep(solution, screenshots, stepCounter);
  }
}

/** Load the content of the issue page 
 * 
 * @param {string} issueId Id of the issue to open
 */ 
export function openIssuePage(issueId) {
  stepCounter = 0;
  const currentIssue = data.find((element) => element.Name === issueId);
  const pageContents = document.getElementById("page-contents");
  pageContents.innerHTML = `
    <h1 class="issue-name">${currentIssue.Name}</h1>
    <div class="issue-information">
      <h2 id="information">Information</h2>
      <p>${currentIssue.Information}</p>
      <h2 id="solution">Solution</h2>
      <div class="issue-solution">
        <p id="solution-text">${currentIssue.Solution[stepCounter]}</p>
        <img style='display:block; width:500px;height:auto' id="step-screenshot"></img>
        <div class="solution-buttons">
          <div class="button-box">
            <div id="previous-button" class="step-button">&laquo; Previous step</div>
            <div id="next-button" class="step-button">Next step &raquo;</div>
          </div>
        </div>
      </div>
    </div>
    <div id="back-button">Back to issues overview</div>
  `;  

    let texts = ["information", "solution", "previous-button", "next-button", "back-button"]
    let localizationIds = ["Issues.Information", "Issues.Solution", "Issues.Previous", "Issues.Next", "Issues.Back"]
    for (let i = 0; i < texts.length; i++) {
        GetLocalization(localizationIds[i], texts[i])
    }

    try {
        document.getElementById("step-screenshot").src = currentIssue.Screenshots[stepCounter];
    } catch (error) { }

  // Add functions to page for navigation
  document.getElementById("next-button").addEventListener("click", () => nextSolutionStep(currentIssue.Solution, currentIssue.Screenshots));
  document.getElementById("previous-button").addEventListener("click", () => previousSolutionStep(currentIssue.Solution, currentIssue.Screenshots));
  document.getElementById("back-button").addEventListener("click", () => openIssuesPage());
}
