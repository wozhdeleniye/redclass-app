import { z } from 'zod';

export const createSubjectSchema = z.object({
  name: z
    .string()
    .min(1, 'Название обязательно')
    .max(100, 'Название не должно превышать 100 символов'),
  description: z
    .string()
    .min(1, 'Описание обязательно')
    .max(500, 'Описание не должно превышать 500 символов'),
  code: z
    .string()
    .min(4, 'Код должен содержать минимум 4 символа')
    .max(20, 'Код не должен превышать 20 символов')
    .regex(/^[A-Z0-9]+$/, 'Код должен содержать только заглавные буквы и цифры'),
});

export const updateSubjectSchema = z.object({
  name: z
    .string()
    .min(1, 'Название обязательно')
    .max(100, 'Название не должно превышать 100 символов')
    .optional(),
  description: z
    .string()
    .min(1, 'Описание обязательно')
    .max(500, 'Описание не должно превышать 500 символов')
    .optional(),
});

export const joinSubjectSchema = z.object({
  code: z
    .string()
    .min(1, 'Код обязателен'),
});

export type CreateSubjectFormData = z.infer<typeof createSubjectSchema>;
export type UpdateSubjectFormData = z.infer<typeof updateSubjectSchema>;
export type JoinSubjectFormData = z.infer<typeof joinSubjectSchema>;

