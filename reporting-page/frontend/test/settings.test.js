import 'jsdom-global/register.js';
import test from 'unit.js';
import {JSDOM} from 'jsdom';
import {jest} from '@jest/globals';
import {mockPageFunctions,
  mockGetLocalization,
  mockChangeLanguage,
  mockChangeScanInterval,
  storageMock,
  clickEvent,
  mockOpenPageFunctions} from './mock.js';

global.TESTING = true;

// Mock issue page
const dom = new JSDOM(`
<!DOCTYPE html>
<html>
<body>
    <div class="header">
      <div class="header-hamburger container">
        <span id="header-hamburger" class="header-hamburger material-symbols-outlined">menu</span>
      </div>
      <div class="header-logo">
        <div id="logo-button" class="logo-name">
          <img id="logo" alt="logo" src="../src/assets/images/InfoSec-Agent-logo.png">  
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
            <a id="windows-version-button" class="">Windows Version</a>
            <a id="scan-interval-button" class="lang-scan-interval"></a>
          </div>
        </div>
      </div>
    </div>
    <div id="window-version-modal" class="modal">
      <div class="modal-content">
        <div class="modal-header">
          <span id="close-window-version-modal" class="close">&times;</span>
          <p>Select windows version</p>
        </div>
        <div id="windows-10"><a id="windows-10-button">Windows 10</a></div>
        <div id="windows-11"><a id="windows-11-button">Windows 11</a></div>
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
  ChangeScanInterval: jest.fn().mockImplementationOnce(() => mockChangeScanInterval(true))
    .mockImplementation(() => mockChangeScanInterval(false)),
  LogError: jest.fn(),
}));

// Mock issues
jest.unstable_mockModule('../src/js/issues.js', () => ({
  openIssuesPage: jest.fn(),
  getUserSettings: jest.fn().mockImplementation(() => 1),
  sortTable: jest.fn(),
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
      '../src/js/programs.js',
      '../src/js/all-checks.js',
      '../src/js/integration.js',
      '../src/js/about.js',
      '../src/js/personalize.js',
    ];

    const pageFunctions = [
      'openHomePage',
      'openSecurityDashboardPage',
      'openPrivacyDashboardPage',
      'openIssuesPage',
      'openProgramsPage',
      'openAllChecksPage',
      'openIntegrationPage',
      'openAboutPage',
      'openPersonalizePage',
    ];

    paths.forEach(async (path, index) => {
      // Arrange
      const page = await import(path);
      const openPageMock = jest.spyOn(page, pageFunctions[index]);

      // Act
      sessionStorage.setItem('languageChanged', true);
      sessionStorage.setItem('savedPage', index+1);
      settings.reloadPage();

      // Assert
      expect(openPageMock).toHaveBeenCalled();
      test.value(sessionStorage.getItem('languageChanged')).isUndefined();
    });

    // Act
    sessionStorage.setItem('languageChanged', true);
    sessionStorage.setItem('savedPage', 0);
    settings.reloadPage();

    // Assert
    expect(logErrorMock).toHaveBeenCalled();
  });
  it('reloadPage doesnt reload when no page is called', async function() {
    // Arrange
    const settings = await import('../src/js/settings.js');

    const tray = await import('../wailsjs/go/main/Tray.js');
    const logErrorMock = jest.spyOn(tray, 'LogError');

    const paths = [
      '../src/js/home.js',
      '../src/js/security-dashboard.js',
      '../src/js/privacy-dashboard.js',
      '../src/js/issues.js',
      '../src/js/programs.js',
      '../src/js/all-checks.js',
      '../src/js/integration.js',
      '../src/js/about.js',
      '../src/js/personalize.js',
    ];

    const pageFunctions = [
      'openHomePage',
      'openSecurityDashboardPage',
      'openPrivacyDashboardPage',
      'openIssuesPage',
      'openProgramsPage',
      'openAllChecksPage',
      'openIntegrationPage',
      'openAboutPage',
      'openPersonalizePage',
    ];

    paths.forEach(async (path, index) => {
      // Arrange
      const page = await import(path);
      const openPageMock = jest.spyOn(page, pageFunctions[index]);

      // Act
      sessionStorage.setItem('languageChanged', true);
      sessionStorage.setItem('savedPage', index+1);
      settings.reloadPage();

      // Assert
      expect(openPageMock).toHaveBeenCalled();
      test.value(sessionStorage.getItem('languageChanged')).isUndefined();
    });

    // Act
    sessionStorage.setItem('savedPage', 0);
    settings.reloadPage();

    // Assert
    expect(logErrorMock).toHaveBeenCalled();
  });
  it('the personalize page can be opened from the settings', async function() {
    // Arrange
    await import('../src/js/settings.js');
    const personalize = await import('../src/js/personalize.js');
    const openPageMock = jest.spyOn(personalize, 'openPersonalizePage');

    // Act
    document.getElementById('personalize-button').dispatchEvent(clickEvent);

    // Assert
    expect(openPageMock).toHaveBeenCalled();
  });
  it('clicking the windows version button should call showWindowsVersion', async function() {
    // Arrange
    await import('../src/js/settings.js');
    sessionStorage.setItem('WindowsVersion', '10');

    // Act
    document.getElementById('windows-version-button').dispatchEvent(clickEvent);
    const modal = document.getElementById('window-version-modal');
    const selected10 = document.getElementById('windows-10-button').classList.contains('selected');
    const selected11 = document.getElementById('windows-11-button').classList.contains('selected');

    // Arrange
    test.value(modal.style.display == 'block').isTrue();
    test.value(selected10).isTrue();
    test.value(selected11).isFalse();
  });
  it('clicking on one of the windows versions selects it', async function() {
    // Arrange
    await import('../src/js/settings.js');
    sessionStorage.setItem('WindowsVersion', '10');

    // Act
    document.getElementById('windows-10').dispatchEvent(clickEvent);
    let selected10 = document.getElementById('windows-10-button').classList.contains('selected');
    let selected11 = document.getElementById('windows-11-button').classList.contains('selected');
    let version = sessionStorage.getItem('WindowsVersion');
    let changed = sessionStorage.getItem('WindowsVersionChanged');

    // Arrange
    test.value(selected10).isTrue();
    test.value(selected11).isFalse();
    test.value(version).isEqualTo('10');
    test.value(changed).isEqualTo('true');

    // clear session storage
    sessionStorage.removeItem('WindowsVersionChanged');

    // Act
    document.getElementById('windows-11').dispatchEvent(clickEvent);
    selected10 = document.getElementById('windows-10-button').classList.contains('selected');
    selected11 = document.getElementById('windows-11-button').classList.contains('selected');
    version = sessionStorage.getItem('WindowsVersion');
    changed = sessionStorage.getItem('WindowsVersionChanged');

    // Arrange
    test.value(selected10).isFalse();
    test.value(selected11).isTrue();
    test.value(version).isEqualTo('11');
    test.value(changed).isEqualTo('true');

    // clear session storage
    sessionStorage.removeItem('WindowsVersionChanged');
  });
  it('clicking the close button should close the modal', async function() {
    // Arrange
    await import('../src/js/settings.js');
    sessionStorage.setItem('ModalOpen', 'true');
    const modal = document.getElementById('window-version-modal');
    modal.style.display = 'block';

    // Act
    document.getElementById('close-window-version-modal').dispatchEvent(clickEvent);

    // Assert
    test.value(modal.style.display).isEqualTo('none');
    test.value(sessionStorage.getItem('ModalOpen')).isUndefined();
  });
  it('clicking outside the modal should close it', async function() {
    // Arrange
    await import('../src/js/settings.js');
    sessionStorage.setItem('ModalOpen', 'true');
    const modal = document.getElementById('window-version-modal');
    modal.style.display = 'block';

    // Act
    const event = new window.MouseEvent('click', {
      view: window,
      bubbles: true,
      cancelable: true,
    });
    modal.dispatchEvent(event);

    // Assert
    test.value(modal.style.display).isEqualTo('none');

    // Act
    document.dispatchEvent(event);

    // Assert
    test.value(modal.style.display).isEqualTo('none');
  });
  it('should not open the modal if it is already open', async function() {
    // Arrange
    await import('../src/js/settings.js');
    sessionStorage.setItem('ModalOpen', 'true');
    const modal = document.getElementById('window-version-modal');
    modal.style.display = 'none';

    // Act
    document.getElementById('windows-version-button').dispatchEvent(clickEvent);

    // Assert
    test.value(modal.style.display).isEqualTo('none');

    // Act
    document.getElementById('windows-10').dispatchEvent(clickEvent);
    // Assert
    test.value(modal.style.display).isEqualTo('none');
  });
});
