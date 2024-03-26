import cs from "../customize.json" assert { type: "json" };

/** Load the content of the Personalize page */
export function openPersonalizePage() {
  document.getElementById("page-contents").innerHTML = `
  <div class="personalize favicon-mode">
    <span class="personalize-description">Favicon</span>
    <div class="personalize-button-container">
      <label class="personalize-label" for="input-file-icon">Change favicon</label>
      <input class="personalize-input-invisible" type="file" id="input-file-icon" accept=".ico, .png">
    </div>
  </div>
  <hr class="solid">
  <div class="personalize">
    <span class="personalize-description">Navigation image</span>
    <div class="personalize-button-container">
      <label class="personalize-label" for="input-file-picture">Update image</label>
      <input class="personalize-input-invisible" type="file" id="input-file-picture" accept="image/jpeg, image/png, image/jpg">   
    </div>
  </div>
  <hr class="solid">
  <div class="personalize">
    <span class="personalize-description">Navigation title</span>
    <div class="personalize-button-container">
      <label class="personalize-label" for="newTitle">Update title</label>
      <input type="text" id="newTitle">
    </div>
  </div>
  <hr class="solid">
  <div class="personalize">
    <span class="personalize-description">Font</span>
  </div>
  <hr class="solid">
  <div class="personalize">
    <span class="personalize-description">Background color Left nav</span>
    <div class="personalize-button-container">
      <label class="personalize-label" for="input-color-background">Change background color</label>
      <input class="personalize-input-invisible" type="color" id="input-color-background">   
    </div>
  </div>
  <hr class="solid">
  <div class="personalize">
    <span class="personalize-description">text color</span>
  </div>
  `;
  const faviconInput = document.getElementById('input-file-icon');//add eventlistener for changing Favicon
  faviconInput.addEventListener('change', handleFaviconChange);
  
  const pictureInput = document.getElementById('input-file-picture'); //add eventlistener for changing navication picture
  pictureInput.addEventListener('change', handlePictureChange);
  
  const newTitleInput = document.getElementById('newTitle'); //add eventlistener for changing navigation title
  newTitleInput.addEventListener('input', handleTitleChange);

  const inputBackgroundNav = document.getElementById('input-color-background'); //add eventlistener for changing navigation title
  inputBackgroundNav.addEventListener('change', handleLeftBackgroundNav);
  }
  
/* Changes the favicon*/
export function handleFaviconChange(icon) {
  const file = icon.target.files[0]; // Get the selected file
  if (file) {
    const reader = new FileReader();
    reader.onload = function(e) {
      const picture = e.target.result;
      const favicon = document.querySelector('link[rel="icon"]');
      if(favicon){
        favicon.href = picture;
      }
      else{
        const newFavicon = document.createElement('link');
        newFavicon.rel = 'icon';
        newFavicon.href = picture;
        document.head.appendChild(newFavicon);
      }
      localStorage.setItem("favicon", picture);
    };
    reader.readAsDataURL(file); // Read the selected file as a Data URL
  }
}

/* Changes the navigation picture*/
export function handlePictureChange(picture) {
  const file = picture.target.files[0]; // Get the selected file
  const reader = new FileReader();
  reader.onload = function(e) {
    const logo = document.getElementById('logo');
    logo.src = e.target.result; // Set the source of the logo to the selected image
    localStorage.setItem("picture", e.target.result)
    };
  reader.readAsDataURL(file); // Read the selected file as a Data URL
}

/* Changes the title of the page to value of element with id:"newTitle" */
export function handleTitleChange() {
  const newTitle = document.getElementById('newTitle').value; // Get the value from the input field
  const titleElement = document.getElementById('title'); 
  titleElement.textContent = newTitle; // Set the text content to the new title
  localStorage.setItem("title", newTitle);
}

/*Change the left background of the navigation*/
export function handleLeftBackgroundNav(){
  const colorPicker = document.getElementById('input-color-background');
  const color = colorPicker.value;
  let temp = document.getElementById('left-nav');
  temp.style.backgroundColor = color;
}



//achtergrond navigation
//normale achtergrond
//kleur text
//text font
//kleur buttons