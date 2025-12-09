import { apiClient } from './client';
import type {
  Project,
  CreateProjectRequest,
  JoinProjectRequest,
  User,
  ProblemStatistics,
} from '../types';

export const projectsApi = {
  // Получить проекты задания
  getByTask: async (taskId: string): Promise<Project[]> => {
    const response = await apiClient.get<Project[]>(
      `/tasks/${taskId}/projects`
    );
    return response.data;
  },

  // Создать проект
  create: async (
    taskId: string,
    data: CreateProjectRequest
  ): Promise<Project> => {
    const response = await apiClient.post<Project>(
      `/tasks/${taskId}/projects`,
      data
    );
    return response.data;
  },

  // Присоединиться к проекту по коду
  join: async (data: JoinProjectRequest): Promise<Project> => {
    const response = await apiClient.post<Project>('/projects/join', data);
    return response.data;
  },

  // Получить мои проекты
  getMy: async (): Promise<Project[]> => {
    const response = await apiClient.get<Project[]>('/projects/my');
    return response.data;
  },

  // Получить пользователей проекта
  getUsers: async (projectId: string): Promise<User[]> => {
    const response = await apiClient.get<User[]>(
      `/projects/${projectId}/users`
    );
    return response.data;
  },

  // Получить статистику проекта
  getStatistics: async (projectId: string): Promise<ProblemStatistics> => {
    const response = await apiClient.get<ProblemStatistics>(
      `/projects/${projectId}/statistics`
    );
    return response.data;
  },
};

