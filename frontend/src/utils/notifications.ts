import { notifications } from '@mantine/notifications';

export const showSuccessNotification = (message: string) => {
  notifications.show({
    title: 'Успешно',
    message,
    color: 'green',
  });
};

export const showErrorNotification = (message: string) => {
  notifications.show({
    title: 'Ошибка',
    message,
    color: 'red',
  });
};

export const showInfoNotification = (message: string) => {
  notifications.show({
    title: 'Информация',
    message,
    color: 'blue',
  });
};

