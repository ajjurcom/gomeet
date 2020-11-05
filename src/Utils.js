import {Message} from 'view-design';

export const showMessage = (type, content) => {
    return Message[type] && Message[type]({
        content,
        duration: 4,
        closable: true
    });
};