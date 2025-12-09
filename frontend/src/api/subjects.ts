import { apiClient } from './client';
import type {
  Subject,
  CreateSubjectRequest,
  UpdateSubjectRequest,
  JoinSubjectRequest,
  PaginationParams,
  PaginatedResponse,
  Role,
} from '../types';

export const subjectsApi = {
  // Получить все предметы
  getAll: async (params?: PaginationParams): Promise<PaginatedResponse<Subject>> => {
    const response = await apiClient.get<PaginatedResponse<Subject>>('/subjects', { params });
    return response.data;
  },

  // Получить предмет по ID
  getById: async (id: string): Promise<Subject> => {
    const response = await apiClient.get<Subject>(`/subjects/${id}`);
    return response.data;
  },

  // Создать предмет
  create: async (data: CreateSubjectRequest): Promise<Subject> => {
    const response = await apiClient.post<Subject>('/subjects', data);
    return response.data;
  },

  // Обновить предмет
  update: async (id: string, data: UpdateSubjectRequest): Promise<Subject> => {
    const response = await apiClient.put<Subject>(`/subjects/${id}`, data);
    return response.data;
  },

  // Удалить предмет
  delete: async (id: string): Promise<void> => {
    await apiClient.delete(`/subjects/${id}`);
  },

  // Присоединиться к предмету по коду
  join: async (data: JoinSubjectRequest): Promise<Subject> => {
    const response = await apiClient.post<Subject>('/subjects/join', data);
    return response.data;
  },

  // Получить мои предметы
  getMy: async (): Promise<Subject[]> => {
    const response = await apiClient.get<Subject[]>('/subjects/get/my');
    return response.data;
  },

  // Получить участников предмета
  getMembers: async (id: string): Promise<Role[]> => {
    const response = await apiClient.get<Role[]>(`/subjects/${id}/members`);
    return response.data;
  },
};

