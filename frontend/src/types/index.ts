// User types
export interface User {
  id: string;
  email: string;
  nickname: string;
  created_at: string;
}

// Auth types
export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
  nickname: string;
}

export interface AuthResponse {
  access_token: string;
  refresh_token: string;
  user: User;
}

export interface RefreshTokenRequest {
  refresh_token: string;
}

// Subject types
export interface Subject {
  id: string;
  name: string;
  description: string;
  code: string;
  created_at: string;
  updated_at: string;
}

export interface CreateSubjectRequest {
  name: string;
  description: string;
  code: string;
}

export interface UpdateSubjectRequest {
  name?: string;
  description?: string;
}

export interface JoinSubjectRequest {
  code: string;
}

// Role types
export type RoleType = 'student' | 'teacher' | 'admin';

export interface Role {
  id: string;
  user_id: string;
  subject_id: string;
  role_type: RoleType;
  created_at: string;
  user?: User;
}

export interface ChangeRoleRequest {
  role_type: RoleType;
}

// Task types
export interface Task {
  id: string;
  subject_id: string;
  created_by_id: string;
  title: string;
  description: string;
  due_date: string;
  created_at: string;
  updated_at: string;
  subject?: Subject;
  created_by?: User;
}

export interface CreateTaskRequest {
  title: string;
  description: string;
  due_date: string;
}

export interface UpdateTaskRequest {
  title?: string;
  description?: string;
  due_date?: string;
}

// Project types
export interface Project {
  id: string;
  task_id: string;
  creator_id: string;
  title: string;
  description: string;
  code: string;
  created_at: string;
  updated_at: string;
  task?: Task;
  creator?: User;
  members?: User[];
}

export interface CreateProjectRequest {
  title: string;
  description: string;
}

export interface JoinProjectRequest {
  code: string;
}

// Problem Assignee types
export interface ProblemAssignee {
  id: string;
  problem_id: string;
  user_id: string;
  created_at: string;
  updated_at: string;
  user: User;
  problem?: Problem | null;
}

// Problem types
export interface Problem {
  id: string;
  project_id: string;
  creator_id: string;
  parent_id: string | null;
  number: number;
  title: string;
  description: string;
  start_time: string;
  end_time: string;
  solved: boolean;
  created_at: string;
  updated_at: string;
  creator?: User;
  assignees?: ProblemAssignee[];
  subproblems?: Problem[];
}

export interface ChildrenStatistics {
  completed: number;
  incomplete: number;
  total: number;
}

export interface ProblemWithResult {
  problem: Problem;
  result: Result | null;
  children_statistics?: ChildrenStatistics;
}

export interface ProblemStatistics {
  completed: number;
  incomplete: number;
  total: number;
  percentage: number;
}

export interface CreateProblemRequest {
  title: string;
  description: string;
  start_time: string;
  end_time: string;
  assignee_ids: string[];
}

export interface UpdateProblemRequest {
  title?: string;
  description?: string;
  start_time?: string;
  end_time?: string;
  assignee_ids?: string[];
}

// Result types
export interface Result {
  id: string;
  problem_id: string;
  creator_id: string;
  done: boolean;
  comment: string;
  created_at: string;
}

export interface CreateResultRequest {
  done: boolean;
  comment: string;
}

// Pagination
export interface PaginationParams {
  limit?: number;
  offset?: number;
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  limit: number;
  offset: number;
}

// API Error
export interface ApiError {
  message: string;
  status: number;
}

