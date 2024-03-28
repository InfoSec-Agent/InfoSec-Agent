import data from "../database.json" assert { type: "json" };
import { openIssuePage } from "./issue.js";
import { GetLocalization } from './localize.js';

/** Load the content of the Issues page */
export function openIssuesPage() {
  const pageContents = document.getElementById("page-contents");
  pageContents.innerHTML = `
  <table class="issues-table" id="issues-table">
    <thead>
      <tr>
      <th class="issue-column"><span class="table-header name">Name</span><span class="material-symbols-outlined" id="sort-on-issue">swap_vert</span></th>
      <th class="type-column"><span class="table-header type">Type</span><span class="material-symbols-outlined" id="sort-on-type">swap_vert</span></th>
      <th class="risk-column"><span class="table-header risk">Risk level</span><span class="material-symbols-outlined" id="sort-on-risk">swap_vert</span></th>
      </tr>
    </thead>
    <tbody>
    </tbody>
  </table>
  `;  

  let tableHeaders = ["name", "type", "risk"]
  let localizationIds = ["Issues.Name", "Issues.Type", "Issues.Risk"]
  for (let i = 0; i < tableHeaders.length; i++) {
      GetLocalization(localizationIds[i], tableHeaders[i])
  }

  let issues = []; // retrieve issues from tray application
  issues = [ // dummy info
    { 
      "Id": "Windows defender", 
      "Result": ["Windows defender is disabled"],
      "ErrorMSG": null
    },
    { 
      "Id": "Camera and microphone access", 
      "Result": ["Something has access to camera", "Something has access to microphone"],
      "ErrorMSG": null
    }
  ];  

  const tbody = pageContents.querySelector('tbody');
  fillTable(tbody, issues);
}

/** Fill the table with issues
 * 
 * @param {HTMLTableSectionElement} tbody Table to be filled
 * @param {Issue} issues Issues to be filled in
 */ 
export function fillTable(tbody, issues) {
  issues.forEach(issue => {
    const currentIssue = data.find(element => element.Name === issue.Id);
    if (currentIssue) {
      const row = document.createElement('tr');
      row.innerHTML = `
        <td class="issue-link">${currentIssue.Name}</td>
        <td>${currentIssue.Type}</td>
        <td>${currentIssue.Risk}</td>
      `;
      tbody.appendChild(row);
    }
  });  

  // Add links to issue information pages
  const issueLinks = document.querySelectorAll(".issue-link");
  issueLinks.forEach((link, index) => {
    link.addEventListener("click", () => openIssuePage(issues[index].Id));
  });  

  // Add buttons to sort on columns
  document.getElementById("sort-on-issue").addEventListener("click", () => sortTable(tbody, 0));
  document.getElementById("sort-on-type").addEventListener("click", () => sortTable(tbody, 1));
  document.getElementById("sort-on-risk").addEventListener("click", () => sortTable(tbody, 2));
}

/** Sorts the table 
 * 
 * @param {HTMLTableSectionElement} tbody Table to be sorted
 * @param {string} column Column to sort the table on
 */ 
export function sortTable(tbody, column) {
  const table = tbody.closest("table");
  let direction = table.getAttribute("data-sort-direction");
  direction = direction === "ascending" ? "descending" : "ascending";
  const rows = Array.from(tbody.rows);
  rows.sort((a, b) => {
    if (column === 2) {
      // Custom sorting for the last column
      const order = { "high": 1, "medium": 2, "low": 3, "acceptable": 4 };
      const textA = a.cells[column].textContent.toLowerCase();
      const textB = b.cells[column].textContent.toLowerCase();
      if (direction === "ascending") {
        return order[textA] - order[textB];
      } else {
        return order[textB] - order[textA];
      }
    } else {
      // Alphabetical sorting for other columns
      const textA = a.cells[column].textContent.toLowerCase();
      const textB = b.cells[column].textContent.toLowerCase();
      if (direction === "ascending") {
        return textA.localeCompare(textB);
      } else {
        return textB.localeCompare(textA);
      }
    }
  });
  while (tbody.rows.length > 0) {
    tbody.deleteRow(0);
  }
  rows.forEach(row => {
    tbody.appendChild(row);
  });
  table.setAttribute("data-sort-direction", direction);
}
if (typeof document !== 'undefined') {
  document.getElementById("issues-button").addEventListener("click", () => openIssuesPage());
}