import 'jsdom-global/register.js';
import test from 'unit.js';
import {JSDOM} from 'jsdom';
import {jest} from '@jest/globals';
import {mockPageFunctions,
    mockGetLocalization,
    mockChangeLanguage,
    storageMock,
    clickEvent,
    mockOpenPageFunctions} from './mock.js';

global.TESTING = true;

// Mock issue page
const dom = new JSDOM(`
<!DOCTYPE html>
<html>
<body>
    <input type="file" class="personalize-input-invisible" id="faviconInput" accept=".png,.ico"> <!--Use id to dynamically change favicon -->
    <div class="header">
      <div class="header-hamburger container">
        <span id="header-hamburger" class="header-hamburger material-symbols-outlined">menu</span>
      </div>
      <div class="header-logo">
        <div id="logo-button" class="logo-name">
          <img id="logo" alt="logo" src="./src/assets/images/logoTeamA-transformed.png">  <!-- Use id to dynamically change logo -->
          <div class="header-name">
            <h1 id="title">Little Brother</h1><!-- Use id to dynamically change title -->
          </div>
        </div>
      </div>
      <div class="header-settings">
        <div class="nav-link settings-button">
          <span><span class="material-symbols-outlined">settings</span></span>
          <div class="dropdown-content">
            <a id="personalize-button">Personalize page</a>
            <a id="language-button">Change Language</a>
          </div>
        </div>
      </div>
    </div>
    <div id="page-contents"></div>
</body>
</html>
`);
global.document = dom.window.document;
global.window = dom.window;

// Mock often used page functions
mockPageFunctions();

// Mock all openPage functions
mockOpenPageFunctions();

// Mock Localize function
jest.unstable_mockModule('../wailsjs/go/main/App.js', () => ({
    Localize: jest.fn().mockImplementation((input) => mockGetLocalization(input)),
  }));

// Mock runtime functions
jest.unstable_mockModule('../wailsjs/runtime/runtime.js', () => ({
    WindowReload: jest.fn(),
    LogPrint: jest.fn(),
    WindowShow: jest.fn(),
    WindowMaximise: jest.fn(),
  }));

// Mock scanTest
jest.unstable_mockModule('../src/js/database.js', () => ({
    scanTest: jest.fn(),
  }));

// Mock changeLanguage
jest.unstable_mockModule('../wailsjs/go/main/Tray.js', () => ({
    ChangeLanguage: jest.fn().mockImplementationOnce(() => mockChangeLanguage(true))
        .mockImplementation(() => mockChangeLanguage(false)),
    LogError: jest.fn(),
  }));

// Mock session and localStorage
global.sessionStorage = storageMock;
global.localStorage = storageMock;

describe('Settings page', function() {
    it('updateLanguage calls changeLanguage and reloads the window', async function() {
        // Arrange
        const settings = await import('../src/js/settings.js');

        const tray = await import('../wailsjs/go/main/Tray.js');
        const changeLanguageMock = jest.spyOn(tray, 'ChangeLanguage');
        const logErrorMock = await jest.spyOn(tray, 'LogError');
        
        // Act
        // Clicking language-button calls updateLanguage()
        document.getElementById('language-button').dispatchEvent(clickEvent);

        // Assert
        expect(changeLanguageMock).toHaveBeenCalled();

        // Act
        await settings.updateLanguage();

        // Assert
        expect(logErrorMock).toHaveBeenCalled();
    });
    it('reloadPage reloads the correct page when called', async function() {
        // Arrange
        const settings = await import('../src/js/settings.js');

        const tray = await import('../wailsjs/go/main/Tray.js');
        const logErrorMock = jest.spyOn(tray, 'LogError');

        const paths = [
            '../src/js/home.js',
            '../src/js/security-dashboard.js',
            '../src/js/privacy-dashboard.js',
            '../src/js/issues.js',
            '../src/js/integration.js',
            '../src/js/about.js',
            '../src/js/personalize.js',
        ]

        const pageFunctions = [
            'openHomePage',
            'openSecurityDashboardPage',
            'openPrivacyDashboardPage',
            'openIssuesPage',
            'openIntegrationPage',
            'openAboutPage',
            'openPersonalizePage',
        ]

        paths.forEach(async (path, index) => {
            // Arrange
            const page = await import(path);
            const openPageMock = jest.spyOn(page, pageFunctions[index]);

            // Act
            sessionStorage.setItem('languageChanged',true);
            sessionStorage.setItem('savedPage',index+1);
            settings.reloadPage();

            // Assert
            expect(openPageMock).toHaveBeenCalled();
            test.value(sessionStorage.getItem('languageChanged')).isUndefined();
        })

        // Act
        sessionStorage.setItem('languageChanged',true);
        sessionStorage.setItem('savedPage',0);
        settings.reloadPage();

        // Assert
        expect(logErrorMock).toHaveBeenCalled();
    });
    it('the personalize page can be opened from the settings', async function() {
        // Arrange
        const settings = await import('../src/js/settings.js');
        const personalize = await import('../src/js/personalize.js');
        const openPageMock = jest.spyOn(personalize, 'openPersonalizePage');

        // Act
        document.getElementById('personalize-button').dispatchEvent(clickEvent);

        // Assert
        expect(openPageMock).toHaveBeenCalled();
    })
});
