import { Localize } from '../../wailsjs/go/main/App';

export function GetLocalization(messageId, elementClass) {
    Localize(messageId).then((result) => {
        let elements = document.getElementsByClassName(elementClass);
        for (let i = 0; i < elements.length; i++) {
            elements[i].innerHTML = result;
        }
    });
}