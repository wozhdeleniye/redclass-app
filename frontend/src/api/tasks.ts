import { apiClient } from './client';
import type {
  Task,
  CreateTaskRequest,
  UpdateTaskRequest,
  PaginationParams,
  PaginatedResponse,
} from '../types';

export const tasksApi = {
  // Получить задания предмета
  getBySubject: async (
    subjectId: string,
    params?: PaginationParams
  ): Promise<PaginatedResponse<Task>> => {
    const response = await apiClient.get<PaginatedResponse<Task>>(
      `/subjects/${subjectId}/tasks`,
      { params }
    );
    return response.data;
  },

  // Получить задание по ID
  getById: async (id: string): Promise<Task> => {
    const response = await apiClient.get<Task>(`/tasks/${id}`);
    return response.data;
  },

  // Создать задание
  create: async (subjectId: string, data: CreateTaskRequest): Promise<Task> => {
    const response = await apiClient.post<Task>(
      `/subjects/${subjectId}/tasks`,
      data
    );
    return response.data;
  },

  // Обновить задание
  update: async (id: string, data: UpdateTaskRequest): Promise<Task> => {
    const response = await apiClient.put<Task>(`/tasks/${id}`, data);
    return response.data;
  },
};

