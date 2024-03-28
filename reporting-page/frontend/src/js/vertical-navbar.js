import logo from '../assets/images/logoTeamA-transformed.png';
import { GetLocalization } from './localize.js';

document.getElementById('logo').src = logo;

let navbarItems = ["settings", "home", "security-dashboard","privacy-dashboard", "issues", "integration", "about"]
let localizationIds = ["Navigation.Settings", "Navigation.Home", "Navigation.SecurityDashboard", "Navigation.PrivacyDashboard", "Navigation.Issues", "Navigation.Integration", "Navigation.About"]
for (let i = 0; i < navbarItems.length; i++) {
    GetLocalization(localizationIds[i], navbarItems[i])
}
