import { z } from 'zod';

export const createTaskSchema = z.object({
  title: z
    .string()
    .min(1, 'Название обязательно')
    .max(200, 'Название не должно превышать 200 символов'),
  description: z
    .string()
    .min(1, 'Описание обязательно')
    .max(2000, 'Описание не должно превышать 2000 символов'),
  due_date: z
    .string()
    .min(1, 'Срок выполнения обязателен'),
});

export const updateTaskSchema = z.object({
  title: z
    .string()
    .min(1, 'Название обязательно')
    .max(200, 'Название не должно превышать 200 символов')
    .optional(),
  description: z
    .string()
    .min(1, 'Описание обязательно')
    .max(2000, 'Описание не должно превышать 2000 символов')
    .optional(),
  due_date: z
    .string()
    .optional(),
});

export type CreateTaskFormData = z.infer<typeof createTaskSchema>;
export type UpdateTaskFormData = z.infer<typeof updateTaskSchema>;

