import 'jsdom-global/register.js';
import test from 'unit.js';
import {JSDOM} from 'jsdom';
import {jest} from '@jest/globals';
import data from '../src/databases/database.en-GB.json' assert { type: 'json' };
import dataDe from '../src/databases/database.de.json' assert { type: 'json' };
import dataEnUS from '../src/databases/database.en-US.json' assert { type: 'json' };
import dataEs from '../src/databases/database.es.json' assert { type: 'json' };
import dataFr from '../src/databases/database.fr.json' assert { type: 'json' };
import dataNl from '../src/databases/database.nl.json' assert { type: 'json' };
import dataPt from '../src/databases/database.pt.json' assert { type: 'json' };
import {mockPageFunctions, clickEvent, storageMock, scanResultMock} from './mock.js';

global.TESTING = true;

const dom = new JSDOM(`
  <div id="page-contents"></div>
  <div class="page-contents"></div>
`);
global.document = dom.window.document;
global.window = dom.window;

/** empty the table to have it empty for next tests
 *
 * @param {HTMLTableElement} table table to delete all rows from
 */
function emptyTable(table) {
  while (table.rows.length > 0) {
    table.deleteRow(0);
  }
}

// Mock sessionStorage
global.sessionStorage = storageMock;

/** Mock of getLocalizationString function
 *
 * @param {string} messageID - The ID of the message to be localized.
 * @return {string} The localized string.
 */
function mockGetLocalizationString(messageID) {
  const myPromise = new Promise(function(myResolve, myReject) {
    switch (messageID) {
    case 'Issues.IssueTable':
      myResolve('Issue table');
    case 'Issues.AcceptableFindings':
      myResolve('Acceptable findings');
    case 'Issues.Name':
      myResolve('Name');
    case 'Issues.Type':
      myResolve('Type');
    case 'Issues.Risk':
      myResolve('Risk');
    case 'Dashboard.HighRisk':
      myResolve('HighRisk');
    case 'Dashboard.MediumRisk':
      myResolve('MediumRisk');
    case 'Dashboard.LowRisk':
      myResolve('LowRisk');
    case 'Dashboard.Acceptable':
      myResolve('Acceptable');
    case 'Dashboard.InfoRisk':
      myResolve('InfoRisk');
    case 'Dashboard.SelectRisks':
      myResolve('SelectRisks');
    case 'Issues.Acceptable':
      myResolve('Acceptable');
    case 'Issues.Low':
      myResolve('Low');
    case 'Issues.Medium':
      myResolve('Medium');
    case 'Issues.High':
      myResolve('High');
    case 'Issues.Info':
      myResolve('Info');
    case 'Issues.Failed':
      myResolve('Failed');
    default:
      myReject(new Error('Wrong message ID'));
    }
  });
  return myPromise;
}

// Mock often used page functions
mockPageFunctions();

// Mock Localize function
jest.unstable_mockModule('../wailsjs/go/main/App.js', () => ({
  Localize: jest.fn().mockImplementation((input) => mockGetLocalizationString(input)),
  LoadUserSettings: jest.fn(),
}));

// Mock openIssuesPage
jest.unstable_mockModule('../src/js/issue.js', () => ({
  openIssuePage: jest.fn(),
}));

// Test cases
describe('Issues table', function() {
  it('openIssuesPage should add the issues to the page-contents', async function() {
    // Arrange
    const issue = await import('../src/js/issues.js');
    // Arrange input issues
    const issues = scanResultMock;

    sessionStorage.setItem('ScanResult', JSON.stringify(issues));
    sessionStorage.setItem('IssuesFilter', JSON.stringify(
      {'high': 1, 'medium': 1, 'low': 1, 'acceptable': 1, 'info': 1},
    ));

    // Act
    await issue.openIssuesPage();

    const name = document.getElementsByClassName('lang-name')[0].innerHTML;
    const type = document.getElementsByClassName('lang-type')[0].innerHTML;
    const risk = document.getElementsByClassName('lang-risk')[0].innerHTML;

    const sorting = JSON.parse(sessionStorage.getItem('IssuesSorting'));

    // Assert
    test.value(name).isEqualTo('Name');
    test.value(type).isEqualTo('Type');
    test.value(risk).isEqualTo('Risk');
    test.value(sorting.column).isEqualTo(2);
    test.value(sorting.direction).isEqualTo('descending');

    // Make issues table empty
    const issueTable = document.getElementById('issues-table').querySelector('tbody');
    emptyTable(issueTable);
  });
  it('toRiskLevel should return the right risk level', async function() {
    // Arrange
    const issue = await import('../src/js/issues.js');

    // act
    const risks = [
      '<td><span class="table-risk-level lang-acceptable"></span></td>',
      '<td><span class="table-risk-level lang-low"></span></td>',
      '<td><span class="table-risk-level lang-medium"></span></td>',
      '<td><span class="table-risk-level lang-high"></span></td>',
      '<td><span class="table-risk-level lang-info"></span></td>',
    ];

    // Assert
    risks.forEach((value, index) => {
      test.value(issue.toRiskLevel(index)).isEqualTo(value);
    });
  });
  it('fillTable should fill the issues table with information from the provided JSON array', async function() {
    // Arrange input issues
    const result = scanResultMock;
    sessionStorage.setItem('ScanResult', JSON.stringify(result));

    const defaultSorting = {'column': 2, 'direction': 'ascending'};
    sessionStorage.setItem('IssuesSorting', JSON.stringify(defaultSorting));

    const issue = await import('../src/js/issues.js');
    const issues = await issue.getIssues();

    // Act
    const issueTable = document.getElementById('issues-table').querySelector('tbody');
    issue.fillTable(issueTable, issues);

    // Assert
    const row = issueTable.rows[0];
    test.value(row.cells[0].textContent).isEqualTo(issues[0].name);
    test.value(row.cells[1].textContent).isEqualTo(issues[0].type);
    test.value(row.cells[2].childNodes[0].classList.contains('lang-acceptable')).isTrue();

    // Make issues table empty
    emptyTable(issueTable);
  });
  it('sortTable should sort the issues table', async function() {
    // Arrange table rows
    const table = dom.window.document.getElementById('issues-table');
    const tbody = table.querySelector('tbody');
    tbody.innerHTML = `
      <tr data-severity="3">
        <td>Windows defender</td>
        <td>Security</td>
        <td class="lang-high">High</td>
      </tr>
      <tr data-severity="1">
        <td>Camera and microphone access</td>
        <td>Privacy</td>
        <td class="lang-low">Low</td>
      </tr>
      <tr data-severity="2">
        <td>Firewall settings</td>
        <td>Security</td>
        <td class="lang-medium">Medium</td>
      </tr>
    `;

    await import('../src/js/issues.js');

    // Act - sort by issue name (first column)
    document.getElementById('sort-on-issue').dispatchEvent(clickEvent);

    // Assert
    let sortedRows = Array.from(tbody.rows);
    const sortedNames = sortedRows.map((row) => row.cells[0].textContent);
    test.array(sortedNames).is(['Windows defender', 'Firewall settings', 'Camera and microphone access']);

    // Act - sort by issue name (first column) again to sort descending
    document.getElementById('sort-on-issue').dispatchEvent(clickEvent);

    // Assert
    let sortedRowsDescending = Array.from(tbody.rows);
    const sortedNamesDescending = sortedRowsDescending.map((row) => row.cells[0].textContent);
    test.array(sortedNamesDescending).is(['Camera and microphone access', 'Firewall settings', 'Windows defender']);

    // Act - sort by type (second column)
    document.getElementById('sort-on-type').dispatchEvent(clickEvent);

    // Assert
    sortedRows = Array.from(tbody.rows);
    const sortedTypes = sortedRows.map((row) => row.cells[1].textContent);
    test.array(sortedTypes).is(['Security', 'Security', 'Privacy']);

    // Act - sort by type (second column) again to sort descending
    document.getElementById('sort-on-type').dispatchEvent(clickEvent);

    // Assert
    sortedRowsDescending = Array.from(tbody.rows);
    const sortedTypesDescending = sortedRowsDescending.map((row) => row.cells[1].textContent);
    test.array(sortedTypesDescending).is(['Privacy', 'Security', 'Security']);

    // Act - sort by severity (third column)
    document.getElementById('sort-on-risk').dispatchEvent(clickEvent);

    // Assert
    sortedRows = Array.from(tbody.rows);
    const sortedRisks = sortedRows.map((row) => row.cells[2].textContent);
    test.array(sortedRisks).is(['High', 'Medium', 'Low']);

    // Act - sort by severity (third column) again to sort ascending
    document.getElementById('sort-on-risk').dispatchEvent(clickEvent);

    // Assert
    sortedRowsDescending = Array.from(tbody.rows);
    const sortedRisksDescending = sortedRowsDescending.map((row) => row.cells[2].textContent);
    test.array(sortedRisksDescending).is(['Low', 'Medium', 'High']);
  });
  it('changeTable should update the table with selected risks', async function() {
    // Arrange
    const issue = await import('../src/js/issues.js');

    // Arrange input issues
    const issues = scanResultMock;
    // Arrange expected table data
    const expectedData = [];
    issues.forEach((issue) => {
      expectedData.push(data[issue.jsonkey]);
    });
    sessionStorage.setItem('ScanResult', JSON.stringify(issues));

    const ids = [
      'select-low-risk-table',
      'select-medium-risk-table',
      'select-high-risk-table',
      'select-info-risk-table',
    ];

    for (let i = -1; i < issues.length - 1; i++) {
      // Act
      let issueTable = document.getElementById('issues-table').querySelector('tbody');
      issue.fillTable(issueTable, issues, true);

      if (i >= 0) document.getElementById(ids[i]).checked = false;
      issue.changeTable();
      issueTable = document.getElementById('issues-table').querySelector('tbody');

      // Assert
      expectedData.forEach((expectedIssue, index) => {
        if (index > i) {
          // const row = issueTable.rows[index - 1 - i];
          // test.value(row.cells[0].textContent).isEqualTo(expectedIssue.Name);
          // test.value(row.cells[1].textContent).isEqualTo(expectedIssue.Type);
          // test.value(row.cells[2].textContent).isEqualTo(issue.toRiskLevel(issues[index].severity));
        }
      });
    }
  });
  it('clicking on an issue should open the issue page', async function() {
    // Arrange
    const issue = await import('../src/js/issue.js');
    const issueLinks = document.querySelectorAll('.issue-link');
    const openIssuePageMock = jest.spyOn(issue, 'openIssuePage');

    // Assert
    issueLinks.forEach((link) => {
      link.parentElement.dispatchEvent(clickEvent);
      expect(openIssuePageMock).toHaveBeenCalled();
    });
  });
  it('clicking the select risks toggles show', async function() {
    // Arrange
    await import('../src/js/issues.js');
    const button = document.getElementById('dropbtn-table');
    const myDropdownTable = document.getElementById('myDropdown-table');

    // Act
    button.dispatchEvent(clickEvent);

    // Arrange
    expect(myDropdownTable.classList.contains('show')).toBe(true);

    // Act
    button.dispatchEvent(clickEvent);

    // Arrange
    expect(myDropdownTable.classList.contains('show')).toBe(false);
  });
  it('should show when a check has failed', async function() {
    // make sure filters are on
    const filters = {high: 1, medium: 1, low: 1, acceptable: 1, info: 1};
    sessionStorage.setItem('IssuesFilter', JSON.stringify(filters));

    // Arrange input issues
    const result = [{issue_id: 1, result_id: -1, result: []}];
    sessionStorage.setItem('ScanResult', JSON.stringify(result));

    const issue = await import('../src/js/issues.js');
    const issues = await issue.getIssues();

    // Act
    const issueTable = document.getElementById('issues-table').querySelector('tbody');
    issue.fillTable(issueTable, issues);

    // Assert
    const row = issueTable.rows[0];
    test.value(row.cells[0].textContent).isEqualTo(issues[0].name);
    test.value(row.cells[1].textContent).isEqualTo(issues[0].type);
    test.value(row.cells[0].classList.contains('issue-check-failed')).isTrue();
  });
  it('should use the correct data object based on user language settings', async () => {
    // make sure filters are on
    const filters = {high: 1, medium: 1, low: 1, acceptable: 1, info: 1};
    sessionStorage.setItem('IssuesFilter', JSON.stringify(filters));

    // Define the language settings and the corresponding expected data
    const languageSettings = [
      {language: 0, expectedData: dataDe},
      {language: 1, expectedData: data},
      {language: 2, expectedData: dataEnUS},
      {language: 3, expectedData: dataEs},
      {language: 4, expectedData: dataFr},
      {language: 5, expectedData: dataNl},
      {language: 6, expectedData: dataPt},
      {language: 999, expectedData: data}, // Default case
    ];
    const loadUserSettingsMock = jest.spyOn(await import('../wailsjs/go/main/App.js'), 'LoadUserSettings');

    for (const {language, expectedData} of languageSettings) {
      loadUserSettingsMock.mockResolvedValueOnce({Language: language});
      // Prepare the issues array
      const data = [
        {issue_id: 5, result_id: 1, result: []}, // assuming 51 exists in all datasets
      ];

      sessionStorage.setItem('ScanResult', JSON.stringify(data));
      const {getIssues} = await import('../src/js/issues.js');
      const issues = await getIssues();

      // Act
      const {fillTable} = await import('../src/js/issues.js');
      const issueTable = document.createElement('tbody');
      fillTable(issueTable, issues);

      // Assert
      const row = issueTable.rows[0];
      const currentIssue = expectedData[issues[0].issue_id];
      test.value(row.cells[0].textContent).isEqualTo(currentIssue[issues[0].result_id].Name);
      test.value(row.cells[1].textContent).isEqualTo(currentIssue.Type);
    }
  });
});
