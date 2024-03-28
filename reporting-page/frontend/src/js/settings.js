import { openPersonalizePage } from "./personalize";
import { ChangeLanguage } from "../../wailsjs/go/main/Tray";
import { GetLocalization } from './localize.js';

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
  <div class="setting personalize">
    <span class="setting-description personalize-title">Personalization</span>
    <button class="setting-button personalize-button" type="button">Personalize</button>
    <!--<div id="personalize-button-box">
      <div class="personalize-button">Personalize</div>
    </div>-->
  </div> 
  <hr class="solid">
  <div class="setting language">
    <span class="setting-description language-title">Language</span>
    <button class="setting-button language-button" type="button">Change Language</button>
  </div> 
  `;

  // Localize the static content of the settings page
  let staticSettingsContent = ["personalize-title", "personalize-button", "language-title", "language-button"]
  let localizationIds = ["Settings.PersonalizeTitle", "Settings.PersonalizeButton", "Settings.ChangeLanguageTitle", "Settings.ChangeLanguageButton"]
  for (let i = 0; i < staticSettingsContent.length; i++) {
      GetLocalization(localizationIds[i], staticSettingsContent[i])
  }

  document.getElementsByClassName("language-button")[0].addEventListener("click", () => updateLanguage());
  document.getElementsByClassName("personalize-button")[0].addEventListener("click", () => openPersonalizePage());
}

document.getElementById("settings-button").addEventListener("click", () => openSettingsPage());