import test from 'unit.js';
import { JSDOM } from "jsdom";
import { MarkSelectedNavigationItem, CloseNavigation, ToggleNavigationHamburger, ToggleNavigationResize } from '../src/js/navigation-menu';

// Mock page
const dom = new JSDOM(`
`);
global.document = dom.window.document
global.window = dom.window

// Test cases
describe('MarkSelectedNavigationItem', function() {
  it('should mark the item in the navigation menu that is selected', function() {
    // Arrange

    // Act

    // Assert
  });
});

describe('CloseNavigation', function() {
  it('should close the navigation menu if a menu item is selected and the screen size is smaller than 800 px', function() {
    // Arrange

    // Act

    // Assert
  });
});

describe('ToggleNavigationHamburger', function() {
  it('should open or close the navigation menu when the hamburger menu is clicked', function() {
    // Arrange

    // Act

    // Assert
  });
});

describe('ToggleNavigationResize', function() {
  it('should open or close the navigation menu when the screen is resized', function() {
    // Arrange

    // Act

    // Assert
  });
});