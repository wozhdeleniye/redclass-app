import { z } from 'zod';

export const createProjectSchema = z.object({
  title: z
    .string()
    .min(1, 'Название обязательно')
    .max(200, 'Название не должно превышать 200 символов'),
  description: z
    .string()
    .min(1, 'Описание обязательно')
    .max(2000, 'Описание не должно превышать 2000 символов'),
});

export const joinProjectSchema = z.object({
  code: z
    .string()
    .min(1, 'Код обязателен'),
});

export type CreateProjectFormData = z.infer<typeof createProjectSchema>;
export type JoinProjectFormData = z.infer<typeof joinProjectSchema>;

