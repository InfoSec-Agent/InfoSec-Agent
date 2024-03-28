import { GetLocalization } from './localize.js';

/** Load the content of the Privacy Dashboard page */
function openPrivacyDashboardPage() {
  document.getElementById("page-contents").innerHTML = `
  <div class="dashboard-data"></div>
  `;
}

document.getElementById("privacy-dashboard-button").addEventListener("click", () => openPrivacyDashboardPage());