import { showSuccessNotification, showErrorNotification, showInfoNotification } from '../utils/notifications';

export const useNotifications = () => {
  return {
    showSuccess: showSuccessNotification,
    showError: showErrorNotification,
    showInfo: showInfoNotification,
  };
};

