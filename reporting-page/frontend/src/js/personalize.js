import {closeNavigation, markSelectedNavigationItem} from './navigation-menu.js';
import {getLocalization} from './localize.js';
import {GetImagePath as getImagePath} from '../../wailsjs/go/main/App.js';

/** Load the content of the Personalize page */
export function openPersonalizePage() {
  closeNavigation();
  markSelectedNavigationItem('personalize-button');
  sessionStorage.setItem('savedPage', '9');

  document.getElementById('page-contents').innerHTML = `  
  <div class="personalize-container">
    <h2 class="lang-personalize-title"></h2>
    <div class="personalize-item">
      <span class="personalize-description lang-nav-image"></span>
      <div class="personalize-button-container">
        <button class="personalize-button logo-button lang-change-image" type="button"></button>    
        <input class="personalize-input-invisible" type="file" 
        id="input-file-picture" accept="image/jpeg, image/png, image/jpg">
      </div>
    </div>
    <hr class="solid">
    <div class="personalize-item">
      <span class="personalize-description lang-nav-title"></span>
      <div class="personalize-button-container">
        <button class="personalize-button title-button lang-change-title" type="button"></button>
        <div id="custom-modal" class="modal">
          <div class="modal-content">
            <input type="text" id="new-title-input" placeholder="Enter new title">
            <button id="saveTitleButton">Save</button>
          </div>
        </div>
      </div>
    </div>
    <hr class="solid">
    <div class="personalize-item">
      <form action="" class="color-picker">
        <fieldset>
          <legend class="lang-pick-theme"></legend>
          <label for="normal" class="lang-light"></label>
          <input type="radio" name="theme" id="normal" checked>
          <label for="dark" class="lang-dark"></label>
          <input type="radio" name="theme" id="dark">
        </fieldset>
      </form>
    </div>
    <div class="personalize-button-container">
      <button class="personalize-button reset-button lang-reset-button" type="button">Reset to Default</button>    
    </div>
  </div>
  `;

  // Localize the static content of the personalize page
  const staticPersonalizePageContent = [
    'lang-navigation-image',
    'lang-change-image',
    'lang-navigation-title',
    'lang-change-title',
    'lang-reset-button',
    'lang-personalize-title',
    'lang-pick-theme',
    'lang-nav-image',
    'lang-nav-title',
    'lang-light',
    'lang-dark',
  ];
  const localizationIds = [
    'Personalize.navImage',
    'Personalize.ChangeImage',
    'Personalize.Title',
    'Personalize.ChangeTitle',
    'Personalize.Reset',
    'Personalize.PersonalizeTitle',
    'Personalize.PickTheme',
    'Personalize.NavImage',
    'Personalize.NavTitle',
    'Personalize.Light',
    'Personalize.Dark',
  ];
  for (let i = 0; i < staticPersonalizePageContent.length; i++) {
    getLocalization(localizationIds[i], staticPersonalizePageContent[i]);
  }

  // add event-listener for changing navigation picture
  const changeLogoButton = document.getElementsByClassName('logo-button')[0];
  const inputLogo = document.getElementById('input-file-picture');

  changeLogoButton.addEventListener('click', function() {
    inputLogo.click();
  });

  inputLogo.addEventListener('change', handlePictureChange);

  // add event-listener for changing navigation title
  const changeTitleButton = document.getElementsByClassName('title-button')[0];
  const customModal = document.getElementById('custom-modal');
  const newTitleInput = document.getElementById('new-title-input');
  const saveTitleButton = document.getElementById('saveTitleButton');

  changeTitleButton.addEventListener('click', function() {
    customModal.style.display = 'block';
    newTitleInput.focus();
  });

  saveTitleButton.addEventListener('click', function() {
    const newTitle = newTitleInput.value.trim();
    if (newTitle !== '') {
      handleTitleChange(newTitle);
      customModal.style.display = 'none';
    }
  });

  /* save themes*/
  const themes = document.querySelectorAll('[name="theme"]');
  themes.forEach((themeOption) => {
    themeOption.addEventListener('click', () => {
      localStorage.setItem('theme', themeOption.id);
      markSelectedNavigationItem('personalize-button');
    });
  });

  const activeTheme = localStorage.getItem('theme');
  themes.forEach((themeOption) => {
    if (themeOption.id === activeTheme) {
      themeOption.checked = true;
    }
  });
  document.documentElement.className= activeTheme;

  // add event-listener for resetting settings
  const changeResetButton = document.getElementsByClassName('reset-button')[0];
  changeResetButton.addEventListener('click', async function() {
    await resetSettings();
  });
}

/**
 * Handles the change event when selecting a new picture file.
 * Updates the source of the specified image element with the selected image.
 * Saves the selected image URL in the localStorage.
 * @param {Event} picture - The event object representing the change of the picture input.
 */
export function handlePictureChange(picture) {
  const file = picture.target.files[0]; // Get the selected file
  const reader = new FileReader();
  reader.onload = function(e) {
    const logo = document.getElementById('logo');
    logo.src = e.target.result; // Set the source of the logo to the selected image
    localStorage.setItem('picture', e.target.result.toString());
  };
  reader.readAsDataURL(file); // Read the selected file as a Data URL
}

/**
 * Handles the change event when updating the title.
 * Updates the text content of the specified title element with the new title value.
 * Saves the new title value in the localStorage.
 * @param {string} newTitle - The new title to set.
 */
export function handleTitleChange(newTitle) {
  const titleElement = document.getElementById('title');
  titleElement.textContent = newTitle; // Set the text content to the new title
  localStorage.setItem('title', newTitle);
}

/**
 * Retrieves the active theme from localStorage and applies it to the document's root element.
 * The active theme class name is retrieved from the 'theme' key in localStorage.
 */
export function retrieveTheme() {
  window.scrollTo(0, 0);
  document.documentElement.className = localStorage.getItem('theme');
}
/**
 * Resets the settings by clearing localStorage and restoring default values.
 */
export async function resetSettings() {
  localStorage.clear();
  const logoPhoto = await getImagePath('InfoSec-Agent-logo.png');
  const logo = document.getElementById('logo');
  logo.src = logoPhoto;

  const title = document.getElementById('title');
  title.textContent = 'InfoSec Agent';

  // Reset theme to light mode
  document.documentElement.className = 'normal';
  localStorage.setItem('theme', 'normal');

  // Update the radio button to reflect the light theme
  const themeRadioButton = document.getElementById('normal');
  themeRadioButton.checked = true;
  markSelectedNavigationItem('personalize-button');
}

