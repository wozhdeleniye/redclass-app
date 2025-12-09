import { z } from 'zod';

export const createProblemSchema = z.object({
  title: z
    .string()
    .min(1, 'Название обязательно')
    .max(200, 'Название не должно превышать 200 символов'),
  description: z
    .string()
    .min(1, 'Описание обязательно')
    .max(2000, 'Описание не должно превышать 2000 символов'),
  start_time: z
    .string()
    .min(1, 'Дата начала обязательна'),
  end_time: z
    .string()
    .min(1, 'Дата окончания обязательна'),
  assignee_ids: z
    .array(z.string())
    .default([]),
});

export const updateProblemSchema = z.object({
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
  start_time: z
    .string()
    .optional(),
  end_time: z
    .string()
    .optional(),
  assignee_ids: z
    .array(z.string())
    .optional(),
});

export type CreateProblemFormData = z.infer<typeof createProblemSchema>;
export type UpdateProblemFormData = z.infer<typeof updateProblemSchema>;

