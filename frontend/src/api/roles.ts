import { apiClient } from './client';
import type { ChangeRoleRequest, Role } from '../types';

export const rolesApi = {
  // Изменить роль пользователя
  changeRole: async (
    subjectId: string,
    roleId: string,
    data: ChangeRoleRequest
  ): Promise<Role> => {
    const response = await apiClient.post<Role>(
      `/subjects/${subjectId}/roles/${roleId}/change`,
      data
    );
    return response.data;
  },

  // Удалить пользователя из предмета
  removeUser: async (subjectId: string, roleId: string): Promise<void> => {
    await apiClient.delete(`/subjects/${subjectId}/roles/${roleId}`);
  },
};

