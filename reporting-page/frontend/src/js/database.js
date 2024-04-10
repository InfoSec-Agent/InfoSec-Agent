import {LogError as logError, ScanNow as scanNowGo} from '../../wailsjs/go/main/Tray.js';
import {GetDataBaseData as getDataBaseData} from '../../wailsjs/go/main/DataBase.js';
import {openHomePage} from './home.js';
import * as runTime from '../../wailsjs/runtime/runtime.js';
import * as rc from './risk-counters.js';

/** Call ScanNow in backend and store result in sessionStorage */
try {
  scanNowGo()
    .then((result) => {
      // place result in session storage
      sessionStorage.setItem('ScanResult', JSON.stringify(result));
      // place severities in session storage
      setSeverities(result);

      runTime.WindowShow();
      runTime.WindowMaximise();
      runTime.LogPrint(sessionStorage.getItem('ScanResult'));
    })
    .catch((err) => {
      logError('Error in scanNowGo: ' + err);
    });
} catch (err) {
  logError('Error in scanNowGo: ' + err);
}

// counts the occurences of each level: 0 = acceptable, 1 = low, 2 = medium, 3 = high
const countOccurences = (severities, level) => severities.filter((item) => item.severity === level).length;

/** Sets the severities collected from the checks and database in session storage
 *
 * @param {Check[]} input Checks to get severities from
 * @param {int[]} ids List of result ids to get corresponding severities
 */
function setSeverities(input) {
  getDataBaseData(input)
    .then((result) => {
      sessionStorage.setItem('DataBaseData', JSON.stringify(result));
      const high = countOccurences(result, 3);
      const medium = countOccurences(result, 2);
      const low = countOccurences(result, 1);
      const acceptable = countOccurences(result, 0);
      sessionStorage.setItem('RiskCounters', JSON.stringify(new rc.RiskCounters(high, medium, low, acceptable)));
      if(sessionStorage.getItem('RiskCounters') === null || sessionStorage.getItem('RiskCounters') === undefined){
        sessionStorage.setItem('RiskCounters', JSON.stringify(new rc.RiskCounters(high, medium, low, acceptable)));
      }
      else{
        riskCounter = sessionStorage.getItem('RiskCounters');
        riskCounter.updateRiskcounter(high, medium, low, acceptable);
        sessionStorage.setItem('RiskCounters', JSON.stringify(riskCounter));
      }
      openHomePage();
    })
    .catch((err) => {
      console.error(err);
    });
}

