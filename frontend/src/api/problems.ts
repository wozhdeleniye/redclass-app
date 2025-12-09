import { apiClient } from './client';
import type {
  Problem,
  ProblemWithResult,
  CreateProblemRequest,
  UpdateProblemRequest,
  Result,
  CreateResultRequest,
} from '../types';

export const problemsApi = {
  // Получить проблемы проекта
  getByProject: async (
    projectId: string,
    assignedOnly?: boolean
  ): Promise<Problem[]> => {
    const response = await apiClient.get<Problem[]>(
      `/projects/${projectId}/problems`,
      { params: assignedOnly ? { assigned_only: true } : {} }
    );
    return response.data;
  },

  // Получить проблему по ID (с результатом если solved=true)
  getById: async (id: string): Promise<ProblemWithResult> => {
    const response = await apiClient.get<ProblemWithResult>(`/problems/${id}`);
    return response.data;
  },

  // Создать подпроблему
  createSubproblem: async (
    problemId: string,
    data: CreateProblemRequest
  ): Promise<Problem> => {
    const response = await apiClient.post<Problem>(
      `/problems/${problemId}/subproblems`,
      data
    );
    return response.data;
  },

  // Обновить проблему
  update: async (
    id: string,
    data: UpdateProblemRequest
  ): Promise<Problem> => {
    const response = await apiClient.put<Problem>(`/problems/${id}`, data);
    return response.data;
  },

  // Получить подпроблемы
  getSubproblems: async (parentId: string): Promise<Problem[]> => {
    const response = await apiClient.get<Problem[]>(
      `/problems/${parentId}/subproblems`
    );
    return response.data;
  },

  // Получить результат проблемы
  getResult: async (problemId: string): Promise<Result> => {
    const response = await apiClient.get<Result>(
      `/problems/${problemId}/result`
    );
    return response.data;
  },

  // Создать результат
  createResult: async (
    problemId: string,
    data: CreateResultRequest
  ): Promise<Result> => {
    const response = await apiClient.post<Result>(
      `/problems/${problemId}/result`,
      data
    );
    return response.data;
  },
};

