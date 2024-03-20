import { openPersonalizePage } from "./personalize";
import { ChangeLanguage } from "../../wailsjs/go/main/Tray";

function updateLanguage() {
  ChangeLanguage()
    .then((result) => {
    })
    .catch((err) => {
        console.error(err);
    });
}

/** Load the content of the Settings page */
function openSettingsPage() {
  document.getElementById("page-contents").innerHTML = `
  <div class="setting dark-mode">
    <span class="setting-description">Dark mode</span>
    <label class="switch">
      <input type="checkbox" id="dark-mode-switch">
      <span class="slider round"></span>
    </label>
  </div>
  <hr class="solid">
  <div class="setting">
    <span class="setting-description">Other setting</span>
    <label class="switch">
      <input type="checkbox">
      <span class="slider round"></span>
    </label>
  </div>
  <hr class="solid">
  <div class="setting">
    <span class="setting-description">Other setting</span>
    <label class="switch">
      <input type="checkbox">
      <span class="slider round"></span>
    </label>
  </div>
  <hr class="solid">
  <div class="setting">
    <span class="setting-description">Personalize page</span>
    <div id="personalize-button-box">
      <div id="personalize-button">Click me</div>
    </div>
  </div> 
  `;

  document.getElementById("change-language").addEventListener("click", () => updateLanguage());
  document.getElementById("personalize-button").addEventListener("click", () => openPersonalizePage());
}

document.getElementById("settings-button").addEventListener("click", () => openSettingsPage());