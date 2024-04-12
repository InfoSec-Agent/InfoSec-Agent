import 'jsdom-global/register.js';
import test from 'unit.js';
import {fillTable} from '../src/js/issues.js';
import {sortTable} from '../src/js/issues.js';
import {JSDOM} from 'jsdom';

global.TESTING = true;

// Mock page
const dom = new JSDOM(`
<!DOCTYPE html>
<html>
<body>
  <table class="issues-table" id="issues-table">
    <thead>
      <tr>
        <th class="issue-column">
          <span class="table-header">Name</span>
          <span class="material-symbols-outlined" id="sort-on-issue">swap_vert</span>
        </th>
        <th class="type-column">
          <span class="table-header">Type</span>
          <span class="material-symbols-outlined" id="sort-on-type">swap_vert</span>
        </th>
        <th class="risk-column">
          <span class="table-header">Risk level</span>
          <span class="material-symbols-outlined" id="sort-on-risk">swap_vert</span>
        </th>
      </tr>
    </thead>
    <tbody>
    </tbody>
  </table>
</body>
</html>
`);

// Test cases
describe('Issues table', function() {
  it('fillTable should fill the issues table with information from the provided JSON array', function() {
    // Arrange input issues
    let issues = [];
    issues = [
      {
        'Id': 'Windows defender',
        'Result': ['Windows defender is disabled'],
        'ErrorMSG': null,
      },
      {
        'Id': 'Camera and microphone access',
        'Result': ['Something has access to camera', 'Something has access to microphone'],
        'ErrorMSG': null,
      },
    ];
    // Arrange expected table data
    const expectedData = [
      {
        'Name': 'Windows defender',
        'Type': 'Security',
        'Risk': 'High',
      },
      {
        'Name': 'Camera and microphone access',
        'Type': 'Privacy',
        'Risk': 'Low',
      },
    ];
    // Act
    const tbody = global.document.querySelector('tbody');
    fillTable(tbody, issues);
    // Assert
    expectedData.forEach((expectedIssue, index) => {
      const row = tbody.rows[index];
      test.value(row.cells[0].textContent).isEqualTo(expectedData[index].Name);
      test.value(row.cells[1].textContent).isEqualTo(expectedData[index].Type);
      test.value(row.cells[2].textContent).isEqualTo(expectedData[index].Risk);
    });
  });

  it('sortTable should sort the issues table', function() {
    // Arrange table rows
    const table = dom.window.document.getElementById('issues-table');
    const tbody = table.querySelector('tbody');
    tbody.innerHTML = `
      <tr>
        <td>Windows defender</td>
        <td>Security</td>
        <td>High</td>
      </tr>
      <tr>
        <td>Camera and microphone access</td>
        <td>Privacy</td>
        <td>Low</td>
      </tr>
      <tr>
        <td>Firewall settings</td>
        <td>Security</td>
        <td>Medium</td>
      </tr>
    `;

    // Act
    sortTable(tbody, 0);

    // Assert
    let sortedRows = Array.from(tbody.rows);
    const sortedNames = sortedRows.map((row) => row.cells[0].textContent);
    test.array(sortedNames).is(['Camera and microphone access', 'Firewall settings', 'Windows defender']);

    // Act
    sortTable(tbody, 0);

    // Assert
    let sortedRowsDescending = Array.from(tbody.rows);
    const sortedNamesDescending = sortedRowsDescending.map((row) => row.cells[0].textContent);
    test.array(sortedNamesDescending).is(['Windows defender', 'Firewall settings', 'Camera and microphone access']);

    // Act
    sortTable(tbody, 2);

    // Assert
    sortedRows = Array.from(tbody.rows);
    const sortedRisks = sortedRows.map((row) => row.cells[2].textContent);
    test.array(sortedRisks).is(['High', 'Medium', 'Low']);

    // Act
    sortTable(tbody, 2);

    // Assert
    sortedRowsDescending = Array.from(tbody.rows);
    const sortedRisksDescending = sortedRowsDescending.map((row) => row.cells[2].textContent);
    test.array(sortedRisksDescending).is(['Low', 'Medium', 'High']);
  });
});
