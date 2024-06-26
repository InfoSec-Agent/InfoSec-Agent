import 'jsdom-global/register.js';
import test from 'unit.js';
import {JSDOM} from 'jsdom';
import {jest} from '@jest/globals';
import {mockPageFunctions, storageMock} from './mock.js';

global.TESTING = true;

// Mock issue page
const dom = new JSDOM(`
<!DOCTYPE html>
<html>
<body>
    <div id="page-contents">
      <div id="progress-segment" class="data-segment">
        <div class="data-segment-header">
          <p id="lighthouse-progress-header" class="lang-lighthouse-progress"></p>
          <div id="lighthouse-progress-hoverbox">
            <img id="lighthouse-progress-tooltip">
            <p class="lighthouse-progress-tooltip-text lang-tooltip-text"></p>
          </div>
        </div>
        <div id="progress-bar-container" class="progress-container">
          <div class="progress-bar" id="progress-bar"></div>
        </div>
        <p id="progress-percentage-text" class="gamification-text"></p>
        <p id="progress-text" class="lang-progress-text gamification-text"></p>
        <p id="progress-almost-text" class="lang-progress-almost-text gamification-text"></p></p>
        <p id="progress-done-text" class="lang-progress-done-text gamification-text"</p>
      </div>
        <div id="share-modal" class="modal">
            <div class="modal-content">
              <div class="modal-header">
                <span id="close-share-modal" class="close">&times;</span>
                <p>Select where to share your progress, Save and download it, then share it with others!</p>
              </div>
              <div id="share-node" class="modal-body">
                <img class="api-key-image" src="https://placehold.co/600x315" alt="Step 1 Image">
              </div>
              <div id="share-buttons" class="modal-body">
                <a id="share-save-button" class="modal-button share-button">Save</a>
                <a class="share-button-break">|</a>
                <a id="select-facebook" class="select-button selected">Facebook</a>
                <a id="select-x" class="select-button">X</a>
                <a id="select-linkedin" class="select-button">LinkedIn</a>
                <a id="select-instagram" class="select-button">Instagram</a>
                <a class="share-button-break">|</a>
                <a id="share-button" class="modal-button share-button">Share</a>
              </div>
            </div>
        </div>
    </div>
</body>
</html>
`);
global.document = dom.window.document;
global.window = dom.window;

// mock createObjectURL
window.URL.createObjectURL = jest.fn().mockImplementation((input) => input);

// mock window.open
window.open = jest.fn();

// Mock sessionStorage
global.sessionStorage = storageMock;
global.localStorage = storageMock;

// Mock often used page functions
mockPageFunctions();

// Mock Chart constructor
jest.unstable_mockModule('html-to-image', () => ({
  toBlob: jest.fn().mockImplementation((node, config) => {
    return node.innerHTML + '_' + config.width.toString() + '_' + config.height.toString();
  }),
}));

jest.unstable_mockModule('browser-image-compression', () => ({
  imageCompression: jest.fn().mockImplementation((input, i) => input),
  default: jest.fn().mockImplementation((input) => input),
}));

// Mock openIssuesPage
jest.unstable_mockModule('../src/js/issues.js', () => ({
  getUserSettings: jest.fn().mockImplementationOnce(() => 2),
}));

// Mock Localize function
jest.unstable_mockModule('../wailsjs/go/main/App.js', () => ({
  GetImagePath: jest.fn().mockImplementation((input) => input),
  GetLighthouseState: jest.fn().mockImplementationOnce(() => 0)
    .mockImplementationOnce(() => 1)
    .mockImplementationOnce(() => 2)
    .mockImplementationOnce(() => 3)
    .mockImplementationOnce(() => 4)
    .mockImplementationOnce(() => 5)
    .mockImplementation(() => 0),
}));

describe('share functions', function() {
  beforeAll(() => {
    jest.useFakeTimers('modern');
    jest.setSystemTime(new Date(2000, 5, 1));
  });

  afterAll(() => {
    jest.useRealTimers();
  });

  it('setImage should set the right image as the background of the share image', async function() {
    // Arrange
    const share = await import('../src/js/share.js');

    // Act
    sessionStorage.setItem('ShareSocial', JSON.stringify(share.socialMediaSizes['facebook']));
    await share.setImage('facebook');

    // Assert
    test.value(document.getElementById('share-node').style.backgroundImage).isEqualTo('url(sharing/first-state.png)');

    // Act
    await share.setImage('facebook');

    // Assert
    test.value(document.getElementById('share-node').style.backgroundImage).isEqualTo('url(sharing/second-state.png)');

    // Act
    await share.setImage('facebook');

    // Assert
    test.value(document.getElementById('share-node').style.backgroundImage).isEqualTo('url(sharing/third-state.png)');

    // Act
    await share.setImage('facebook');

    // Assert
    test.value(document.getElementById('share-node').style.backgroundImage).isEqualTo('url(sharing/fourth-state.png)');

    // Act
    await share.setImage('facebook');

    // Assert
    test.value(document.getElementById('share-node').style.backgroundImage).isEqualTo('url(sharing/fifth-state.png)');

    // Act
    await share.setImage('facebook');

    // Assert
    test.value(document.getElementById('share-node').style.backgroundImage).isEqualTo('url(sharing/first-state.png)');
  });
  it('getImage should return a url of the passed node converted to an image', async function() {
    // Arrange
    const share = await import('../src/js/share.js');

    // Act
    const node = document.getElementById('share-node');
    const url = await share.getImage(node, 600, 315);

    // Assert
    test.value(url).isEqualTo(node.innerHTML + '_600_315');
  });
  it('saveProgress should get the image from the html node passed and download it', async function() {
    // Arrange
    const share = await import('../src/js/share.js');

    const linkElement = {
      download: '',
      href: '',
      click: jest.fn(),
    };
    jest.spyOn(document, 'createElement').mockImplementation(() => linkElement);

    // Act
    sessionStorage.setItem('ShareSocial', JSON.stringify(share.socialMediaSizes['facebook']));
    const node = document.getElementById('share-node');
    await share.saveProgress(node);

    // Assert
    test.value(linkElement.download).isEqualTo('Info-Sec-Agent_6-1-2000_facebook.png');
    test.value(linkElement.href).isEqualTo(node.innerHTML + '_600_315');
    expect(linkElement.click).toHaveBeenCalled();
  });
  it('shareProgress should call window.open to the selected social media page', async function() {
    // Arrange
    const share = await import('../src/js/share.js');
    jest.spyOn(window, 'open');

    // Act
    await share.shareProgress();

    // Assert
    expect(window.open).toHaveBeenCalledTimes(1);

    // Act
    sessionStorage.setItem('ShareSocial', JSON.stringify(share.socialMediaSizes['x']));
    share.shareProgress();

    // Assert
    expect(window.open).toHaveBeenCalledTimes(2);

    // Act
    sessionStorage.setItem('ShareSocial', JSON.stringify(share.socialMediaSizes['linkedin']));
    share.shareProgress();

    // Assert
    expect(window.open).toHaveBeenCalledTimes(3);

    // Act
    sessionStorage.setItem('ShareSocial', JSON.stringify(share.socialMediaSizes['instagram']));
    share.shareProgress();

    // Assert
    expect(window.open).toHaveBeenCalledTimes(4);

    // Act
    sessionStorage.setItem('ShareSocial', JSON.stringify(''));
    share.shareProgress();

    // Assert
    expect(window.open).toHaveBeenCalledTimes(4);
  });
  it('selectSocialMedia should select the right social media and set it in the session storage', async function() {
    // Arrange
    const share = await import('../src/js/share.js');
    const socialMedias = ['facebook', 'x', 'linkedin', 'instagram'];

    socialMedias.forEach((social) => {
      // Act
      share.selectSocialMedia(social);

      // Assert
      // The right social media box is selected
      socialMedias.forEach((social2) => {
        if (social == social2) {
          test.value(document.getElementById('select-' + social2).classList.contains('selected')).isTrue();
        } else test.value(document.getElementById('select-' + social2).classList.contains('selected')).isFalse();
      });

      const inStorage = JSON.parse(sessionStorage.getItem('ShareSocial'));
      test.value(inStorage.name).isEqualTo(share.socialMediaSizes[social].name);
      test.value(inStorage.height).isEqualTo(share.socialMediaSizes[social].height);
      test.value(inStorage.width).isEqualTo(share.socialMediaSizes[social].width);
    });
  });
});
